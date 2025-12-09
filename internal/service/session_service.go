package service

import (
	"time"

	"pomodoro-backend/internal/domain"
)

// Errores estándar usados por el servicio

type SessionService struct {
	sessionRepo domain.SessionRepository
	taskRepo    domain.TaskRepository
}

func NewSessionService(sr domain.SessionRepository, tr domain.TaskRepository) *SessionService {
	return &SessionService{
		sessionRepo: sr,
		taskRepo:    tr,
	}
}

//
// ─────────────────────────────────────────────────────────────
//   CREAR E INICIAR SESIÓN POMODORO
// ─────────────────────────────────────────────────────────────
//

func (s *SessionService) CreateAndStartSession(
	userID string,
	projectID *string,
	taskID *string,
	focusMin int,
	breakMin int,
) (*domain.Session, error) {

	now := time.Now()

	session := &domain.Session{
		UserID:        userID,
		ProjectID:     projectID,
		TaskID:        taskID,
		FocusMinutes:  focusMin,
		BreakMinutes:  breakMin,
		State:         domain.SessionStateRunning,
		StartedAt:     now,
		CreatedAt:     now,
		UpdatedAt:     now,
		Interruptions: 0,
	}

	// Guardar la sesión
	if err := s.sessionRepo.CreateSession(session); err != nil {
		return nil, err
	}

	// Si está ligada a una tarea → ponerla en progreso
	if taskID != nil {
		task, err := s.taskRepo.FindByID(*taskID)
		if err == nil {
			task.Status = domain.TaskStatusInProgress
			task.UpdatedAt = now
			_ = s.taskRepo.Update(task)
		}
	}

	return session, nil
}

//
// ─────────────────────────────────────────────────────────────
//   PAUSAR SESIÓN
// ─────────────────────────────────────────────────────────────
//

func (s *SessionService) PauseSession(id string) (*domain.Session, error) {
	session, err := s.sessionRepo.FindByID(id)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if session.State != domain.SessionStateRunning {
		return nil, ErrInvalidStateTransition
	}

	now := time.Now()
	session.State = domain.SessionStatePaused
	session.PausedAt = &now
	session.UpdatedAt = now

	// Contar interrupciones
	session.Interruptions++

	return session, s.sessionRepo.UpdateSession(session)
}

//
// ─────────────────────────────────────────────────────────────
//   REANUDAR SESIÓN
// ─────────────────────────────────────────────────────────────
//

func (s *SessionService) ResumeSession(id string) (*domain.Session, error) {
	session, err := s.sessionRepo.FindByID(id)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if session.State != domain.SessionStatePaused {
		return nil, ErrInvalidStateTransition
	}

	now := time.Now()
	session.State = domain.SessionStateRunning
	session.PausedAt = nil
	session.UpdatedAt = now

	return session, s.sessionRepo.UpdateSession(session)
}

// ─────────────────────────────────────────────────────────────
//
//	FINALIZAR SESIÓN (FIN DEL FOCUS)
//
// ─────────────────────────────────────────────────────────────
func (s *SessionService) FinishSession(id string) (*domain.Session, PomodoroCycleState, error) {
	session, err := s.sessionRepo.FindByID(id)
	if err != nil {
		return nil, PomodoroCycleState{}, ErrSessionNotFound
	}

	if session.State != domain.SessionStateRunning &&
		session.State != domain.SessionStatePaused {
		return nil, PomodoroCycleState{}, ErrInvalidStateTransition
	}

	now := time.Now()
	session.State = domain.SessionStateFinished
	session.FinishedAt = &now
	session.UpdatedAt = now

	// Estado neutro por defecto (cuando no hay tarea asociada)
	cycleState := PomodoroCycleState{
		TotalPomodoros: 0,
		IndexInCycle:   0,
		IsCycleEnd:     false,
		NextBreakMin:   DefaultPomodoroConfig.ShortBreakMinutes,
		CyclesDone:     0,
	}

	// Si está ligada a una tarea → sumamos métrica del focus y calculamos ciclo
	if session.TaskID != nil {
		focusMinutes := session.FocusMinutes

		// 1Leer la tarea actual para saber cuántos pomodoros llevaba
		task, errTask := s.taskRepo.FindByID(*session.TaskID)
		if errTask == nil {
			totalPomodorosAfter := task.PomodorosCompleted + 1

			// Usa tu helper de pomodoro_cycle.go para saber:
			// - qué pomodoro del ciclo es
			// - si toca descanso corto o largo
			// - cuántos ciclos completos llevas
			cycleState = ComputePomodoroState(totalPomodorosAfter, DefaultPomodoroConfig)
		}

		// Actualizar métricas en la tarea (minutos y conteo de pomodoros)
		_ = s.taskRepo.AddRealMinutes(*session.TaskID, focusMinutes)
		_ = s.taskRepo.IncrementPomodoroCount(*session.TaskID)
	}

	// Guardar la sesión actualizada
	if err := s.sessionRepo.UpdateSession(session); err != nil {
		return nil, PomodoroCycleState{}, err
	}

	return session, cycleState, nil
}

//
// ─────────────────────────────────────────────────────────────
//   INICIAR BREAK
// ─────────────────────────────────────────────────────────────
//

func (s *SessionService) StartBreak(id string) (*domain.Session, error) {
	session, err := s.sessionRepo.FindByID(id)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if session.State != domain.SessionStateFinished {
		return nil, ErrInvalidStateTransition
	}

	now := time.Now()
	session.State = domain.SessionStateBreakRunning
	session.BreakStartedAt = &now
	session.UpdatedAt = now

	return session, s.sessionRepo.UpdateSession(session)
}

//
// ─────────────────────────────────────────────────────────────
//   PAUSAR BREAK
// ─────────────────────────────────────────────────────────────
//

func (s *SessionService) PauseBreak(id string) (*domain.Session, error) {
	session, err := s.sessionRepo.FindByID(id)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if session.State != domain.SessionStateBreakRunning {
		return nil, ErrInvalidStateTransition
	}

	now := time.Now()
	session.State = domain.SessionStateBreakPaused
	session.PausedAt = &now
	session.UpdatedAt = now

	return session, s.sessionRepo.UpdateSession(session)
}

//
// ─────────────────────────────────────────────────────────────
//   REANUDAR BREAK
// ─────────────────────────────────────────────────────────────
//

func (s *SessionService) ResumeBreak(id string) (*domain.Session, error) {
	session, err := s.sessionRepo.FindByID(id)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if session.State != domain.SessionStateBreakPaused {
		return nil, ErrInvalidStateTransition
	}

	now := time.Now()
	session.State = domain.SessionStateBreakRunning
	session.PausedAt = nil
	session.UpdatedAt = now

	return session, s.sessionRepo.UpdateSession(session)
}

//
// ─────────────────────────────────────────────────────────────
//   FINALIZAR BREAK
// ─────────────────────────────────────────────────────────────
//

func (s *SessionService) FinishBreak(id string) (*domain.Session, error) {
	session, err := s.sessionRepo.FindByID(id)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if session.State != domain.SessionStateBreakRunning &&
		session.State != domain.SessionStateBreakPaused {
		return nil, ErrInvalidStateTransition
	}

	now := time.Now()
	session.State = domain.SessionStateBreakFinished
	session.BreakFinishedAt = &now
	session.UpdatedAt = now

	return session, s.sessionRepo.UpdateSession(session)
}
