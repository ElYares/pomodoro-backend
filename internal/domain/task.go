package domain

import "time"

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
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	ProjectID   *string    `json:"project_id,omitempty"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
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
}
