package domain

import "time"

// SessionState define el estado de una sesión Pomodoro.
type SessionState string

const (
	// Estados del tiempo de enfoque
	SessionStateRunning   SessionState = "RUNNING"
	SessionStatePaused    SessionState = "PAUSED"
	SessionStateFinished  SessionState = "FINISHED"
	SessionStateCancelled SessionState = "CANCELLED"

	// Estados del tiempo de descanso (break)
	SessionStateBreakRunning  SessionState = "BREAK_RUNNING"
	SessionStateBreakPaused   SessionState = "BREAK_PAUSED"
	SessionStateBreakFinished SessionState = "BREAK_FINISHED"
)

// Session representa una sesión Pomodoro completa (focus + break).
type Session struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	TaskID    *string `json:"task_id,omitempty"`
	ProjectID *string `json:"project_id,omitempty"`

	FocusMinutes int          `json:"focus_minutes"`
	BreakMinutes int          `json:"break_minutes"`
	State        SessionState `json:"state"`

	StartedAt  time.Time  `json:"started_at"`
	PausedAt   *time.Time `json:"paused_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`

	// Campos para Break
	BreakStartedAt  *time.Time `json:"break_started_at,omitempty"`
	BreakFinishedAt *time.Time `json:"break_finished_at,omitempty"`

	Interruptions int `json:"interruptions"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SessionRepository define el contrato de persistencia para las sesiones.
type SessionRepository interface {
	CreateSession(s *Session) error
	UpdateSession(s *Session) error
	FindByID(id string) (*Session, error)
}
