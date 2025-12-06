package domain

import "time"

// PomodoroCycle representa un ciclo completado de un Pomodoro.
// Es un registro HISTÓRICO que no se modifica una vez guardado.
type PomodoroCycle struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	TaskID     string    `json:"task_id"`
	Duration   int       `json:"duration"` // Duración del ciclo en minutos
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	BreakUsed  bool      `json:"break_used"` // Si hubo descanso tras el ciclo
}

// CycleRepository define el contrato mínimo para almacenar y consultar ciclos.
type CycleRepository interface {
	Save(cycle *PomodoroCycle) error
	GetByTask(taskID string) ([]*PomodoroCycle, error)
}
