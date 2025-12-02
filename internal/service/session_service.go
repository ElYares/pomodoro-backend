package service

import (
	"time"

	"pomodoro-backend/internal/domain"
)

// SessionService encapsula la lógica de negocio relacionada con
// la gestión de sesiones Pomodoro.
type SessionService struct {
	repo domain.SessionRepository
}

// NewSessionService crea una nueva instancia de SessionService.
// Recibe una implementación de SessionRepository para mantener
// desacoplada la capa de negocio de la infraestructura.
func NewSessionService(repo domain.SessionRepository) *SessionService {
	return &SessionService{
		repo: repo,
	}
}

// CreateAndStartSession crea una nueva sesión Pomodoro en estado RUNNING.
// Inicializa los campos de auditoría y delega la persistencia al repositorio.
func (s *SessionService) CreateAndStartSession(
	userID string,
	projectID *string,
	taskID *string,
	focusMinutes int,
	breakMinutes int,
) (*domain.Session, error) {
	now := time.Now().UTC()

	session := &domain.Session{
		UserID:       userID,
		ProjectID:    projectID,
		TaskID:       taskID,
		FocusMinutes: focusMinutes,
		BreakMinutes: breakMinutes,
		State:        domain.SessionStateRunning,
		StartedAt:    now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.CreateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

// PauseSession cambia el estado de una sesión a PAUSED si actualmente
// está en RUNNING.
func (s *SessionService) PauseSession(id string) (*domain.Session, error) {
	session, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if session.State != domain.SessionStateRunning {
		return nil, ErrInvalidStateTransition
	}

	now := time.Now().UTC()
	session.State = domain.SessionStatePaused
	session.PausedAt = &now
	session.UpdatedAt = now

	if err := s.repo.UpdateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

// ResumeSession cambia el estado de una sesión de PAUSED a RUNNING.
func (s *SessionService) ResumeSession(id string) (*domain.Session, error) {
	session, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if session.State != domain.SessionStatePaused {
		return nil, ErrInvalidStateTransition
	}

	now := time.Now().UTC()
	session.State = domain.SessionStateRunning
	session.UpdatedAt = now

	if err := s.repo.UpdateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

// FinishSession marca una sesión como finalizada si se encuentra
// en estado RUNNING o PAUSED.
func (s *SessionService) FinishSession(id string) (*domain.Session, error) {
	session, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if session.State != domain.SessionStateRunning && session.State != domain.SessionStatePaused {
		return nil, ErrInvalidStateTransition
	}

	now := time.Now().UTC()
	session.State = domain.SessionStateFinished
	session.FinishedAt = &now
	session.UpdatedAt = now

	if err := s.repo.UpdateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}
