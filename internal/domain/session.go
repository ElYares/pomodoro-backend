package domain

import "time"

// SessionState define el estado de una sesión Pomodoro.
type SessionState string

const (
	SessionStateRunning   SessionState = "RUNNING"
	SessionStatePaused    SessionState = "PAUSED"
	SessionStateFinished  SessionState = "FINISHED"
	SessionStateCancelled SessionState = "CANCELLED"
)

// Session representa una sesión Pomodoro de un usuario.
// Es parte del dominio, por lo que no contiene detalles de infraestructura.
type Session struct {
	ID            string       `json:"id"`
	UserID        string       `json:"user_id"`
	ProjectID     *string      `json:"project_id,omitempty"`
	TaskID        *string      `json:"task_id,omitempty"`
	FocusMinutes  int          `json:"focus_minutes"`
	BreakMinutes  int          `json:"break_minutes"`
	State         SessionState `json:"state"`
	StartedAt     time.Time    `json:"started_at"`
	PausedAt      *time.Time   `json:"paused_at,omitempty"`
	FinishedAt    *time.Time   `json:"finished_at,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	Interruptions int          `json:"interruptions"`
}

// SessionRepository define el contrato de persistencia para las sesiones.
// La lógica de negocio depende de esta interfaz, no de una implementación concreta.
type SessionRepository interface {
	CreateSession(s *Session) error
	UpdateSession(s *Session) error
	FindByID(id string) (*Session, error)
}
