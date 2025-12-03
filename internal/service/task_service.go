package service

import (
	"time"

	"pomodoro-backend/internal/domain"
)

// TaskService
//
// Representa la capa de aplicación o casos de uso para tareas.
// Aquí vive la lógica de negocio y reglas de integridad.
//
// La responsabilidad de esta capa NO es persistencia (repositorio) ni transporte
// (HTTP). Solo orquesta las operaciones necesarias para cumplir casos de uso.
type TaskService struct {
	repo domain.TaskRepository
}

// NewTaskService crea una instancia del servicio utilizando un repositorio
// concreto que implemente TaskRepository.
func NewTaskService(repo domain.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// CreateTask
//
// Caso de uso: El usuario crea una nueva tarea.
// Se asignan timestamps y se delega la persistencia al repositorio.
func (s *TaskService) CreateTask(userID, title, desc string, projectID *string) (*domain.Task, error) {
	now := time.Now().UTC()

	task := &domain.Task{
		UserID:      userID,
		Title:       title,
		Description: desc,
		ProjectID:   projectID,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

// UpdateTask
//
// Caso de uso: El usuario modifica los atributos de una tarea existente.
func (s *TaskService) UpdateTask(id, title, desc string, projectID *string) (*domain.Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	task.Title = title
	task.Description = desc
	task.ProjectID = projectID
	task.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(task); err != nil {
		return nil, err
	}

	return task, nil
}

// DeleteTask elimina una tarea por su identificador.
func (s *TaskService) DeleteTask(id string) error {
	return s.repo.Delete(id)
}

// GetTask obtiene una tarea por ID.
func (s *TaskService) GetTask(id string) (*domain.Task, error) {
	return s.repo.FindByID(id)
}

// GetTasksByUser obtiene todas las tareas pertenecientes a un usuario.
func (s *TaskService) GetTasksByUser(userID string) ([]*domain.Task, error) {
	return s.repo.FindByUser(userID)
}
