package service

import (
	"time"

	"pomodoro-backend/internal/domain"
)

type TaskService struct {
	repo domain.TaskRepository
}

func NewTaskService(r domain.TaskRepository) *TaskService {
	return &TaskService{repo: r}
}

//
// ──────────────────────────────────────────────
//   CREAR TAREA
// ──────────────────────────────────────────────
//

func (s *TaskService) CreateTask(userID, title, desc string, projectID *string) (*domain.Task, error) {

	now := time.Now()

	task := &domain.Task{
		ID:                 GenerateID(),
		UserID:             userID,
		Title:              title,
		Description:        desc,
		ProjectID:          projectID,
		Status:             domain.TaskStatusPending,
		Completed:          false,
		PomodorosCompleted: 0,
		TotalFocusMinutes:  0,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.repo.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

//
// ──────────────────────────────────────────────
//   OBTENER TAREA POR ID
// ──────────────────────────────────────────────
//

func (s *TaskService) GetTask(id string) (*domain.Task, error) {
	return s.repo.FindByID(id)
}

//
// ──────────────────────────────────────────────
//   OBTENER TAREAS POR USUARIO
// ──────────────────────────────────────────────
//

func (s *TaskService) GetTasksByUser(userID string) ([]*domain.Task, error) {
	return s.repo.FindByUser(userID)
}

//
// ──────────────────────────────────────────────
//   ACTUALIZAR TAREA
// ──────────────────────────────────────────────
//

func (s *TaskService) UpdateTask(id, title, desc string, projectID *string) (*domain.Task, error) {

	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	task.Title = title
	task.Description = desc
	task.ProjectID = projectID
	task.UpdatedAt = time.Now()

	if err := s.repo.Update(task); err != nil {
		return nil, err
	}

	return task, nil
}

//
// ──────────────────────────────────────────────
//   ELIMINAR TAREA
// ──────────────────────────────────────────────
//

func (s *TaskService) DeleteTask(id string) error {
	return s.repo.Delete(id)
}

//
// ──────────────────────────────────────────────
//   MARCAR COMPLETADA
// ──────────────────────────────────────────────
//

func (s *TaskService) MarkCompleted(id string) error {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return ErrTaskNotFound
	}

	now := time.Now()
	task.Status = domain.TaskStatusCompleted
	task.Completed = true
	task.CompletedAt = &now
	task.UpdatedAt = now

	return s.repo.Update(task)
}

//
// ──────────────────────────────────────────────
//   ACTUALIZAR STATUS
// ──────────────────────────────────────────────
//

func (s *TaskService) UpdateStatus(id string, status domain.TaskStatus) error {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return ErrTaskNotFound
	}

	task.Status = status
	task.UpdatedAt = time.Now()

	return s.repo.Update(task)
}

//
// ──────────────────────────────────────────────
//   MÉTRICAS POMODORO DESDE SESSION SERVICE
// ──────────────────────────────────────────────
//

func (s *TaskService) AddRealMinutes(id string, minutes int) error {
	return s.repo.AddRealMinutes(id, minutes)
}

func (s *TaskService) IncrementPomodoroCount(id string) error {
	return s.repo.IncrementPomodoroCount(id)
}
