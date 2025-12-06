package service

import (
	"time"

	"pomodoro-backend/internal/domain"
)

// CycleService maneja la lógica de negocio para ciclos pomodoro.
type CycleService struct {
	repo domain.CycleRepository
}

// NewCycleService crea el servicio.
func NewCycleService(repo domain.CycleRepository) *CycleService {
	return &CycleService{repo: repo}
}

// RegisterCycle guarda un ciclo terminado.
func (s *CycleService) RegisterCycle(userID string, taskID string, duration int, startedAt time.Time, finishedAt time.Time, breakUsed bool) error {
	cycle := &domain.PomodoroCycle{
		ID:         GenerateID(),
		UserID:     userID,
		TaskID:     taskID,
		Duration:   duration,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		BreakUsed:  breakUsed,
	}

	return s.repo.Save(cycle)
}

// GetCyclesByTask devuelve los ciclos anteriores de una tarea.
func (s *CycleService) GetCyclesByTask(taskID string) ([]*domain.PomodoroCycle, error) {
	return s.repo.GetByTask(taskID)
}
