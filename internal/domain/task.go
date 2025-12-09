package domain

import "time"

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "PENDING"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusPaused     TaskStatus = "PAUSED"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
)

// Task
//
// Representa una tarea creada por un usuario dentro del sistema de productividad.
// El modelo pertenece estrictamente a la capa de dominio, por lo que se mantiene
// libre de dependencias de infraestructura (Mongo, HTTP, frameworks, etc.).
//
// Este objeto puede ser manipulado por la capa de servicios (casos de uso) y
// persistido por la capa de repositorios.
//
// Atributos clave:
// - UserID: propietario de la tarea
// - Title / Description: contenido básico de la tarea
// - ProjectID: relación opcional con un proyecto
// - Completed / CompletedAt: permiten saber si está finalizada
// - CreatedAt / UpdatedAt: auditoría aplicada por el servicio
type Task struct {
	ID          string  `bson:"_id,omitempty" json:"id"`
	UserID      string  `bson:"user_id" json:"user_id"`
	Title       string  `bson:"title" json:"title"`
	Description string  `bson:"description" json:"description"`
	ProjectID   *string `bson:"project_id,omitempty" json:"project_id,omitempty"`

	Status      TaskStatus `bson:"status" json:"status"`
	Completed   bool       `bson:"completed" json:"completed"`
	CompletedAt *time.Time `bson:"completed_at,omitempty" json:"completed_at,omitempty"`

	// Métricas
	PomodorosCompleted int `bson:"pomodoros_completed" json:"pomodoros_completed"`
	TotalFocusMinutes  int `bson:"total_focus_minutes" json:"total_focus_minutes"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// TaskRepository
//
// Interfaz que define los métodos necesarios para manipular tareas desde la capa
// de persistencia. La capa de servicios depende únicamente de esta interfaz,
// permitiendo intercambiar tecnologías (MongoDB, SQL, archivos, mocks, etc.).
type TaskRepository interface {
	Create(task *Task) error
	Update(task *Task) error
	Delete(id string) error
	FindByID(id string) (*Task, error)
	FindByUser(userID string) ([]*Task, error)

	// Nuevos metodos para metricas pomodoro
	UpdateStatus(id string, status TaskStatus) error
	AddRealMinutes(id string, minutes int) error
	IncrementPomodoroCount(id string) error
}
