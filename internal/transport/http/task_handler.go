package http

import (
	"net/http"

	"pomodoro-backend/internal/domain"
	"pomodoro-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// TaskHandler
//
// Controlador HTTP responsable de exponer los casos de uso de tareas
// a través de la capa de transporte. Recibe y valida payloads JSON,
// delega operaciones al servicio y formatea las respuestas HTTP.
type TaskHandler struct {
	svc *service.TaskService
}

// NewTaskHandler construye una instancia del controlador HTTP.
func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

// RegisterRoutes registra todos los endpoints relacionados con tareas.
func (h *TaskHandler) RegisterRoutes(rg *gin.RouterGroup) {
	tasks := rg.Group("/tasks")
	{
		tasks.POST("", h.createTask)
		tasks.GET("/user/:userID", h.getTasksByUser)
		tasks.GET("/:id", h.getTask)
		tasks.PUT("/:id", h.updateTask)
		tasks.DELETE("/:id", h.deleteTask)

		tasks.PATCH("/:id/complete", h.markCompleted)

		tasks.PATCH("/:id/start", h.markInProgress)
		tasks.PATCH("/:id/pause", h.markPaused)
		tasks.PATCH("/:id/reopen", h.reopenTask)
	}
}

// createTaskRequest
//
// Estructura utilizada para validar el cuerpo de la petición POST.
type createTaskRequest struct {
	UserID      string  `json:"user_id" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	ProjectID   *string `json:"project_id"`
}

func (h *TaskHandler) markCompleted(c *gin.Context) {
	id := c.Param("id")

	err := h.svc.MarkCompleted(id)
	if err != nil {
		if err == service.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "tarea no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo completar la tarea"})
		return
	}

	// Opcional: devolver la tarea actualizada
	updated, _ := h.svc.GetTask(id)
	c.JSON(http.StatusOK, updated)
}

func (h *TaskHandler) markInProgress(c *gin.Context) {
	id := c.Param("id")

	err := h.svc.UpdateStatus(id, domain.TaskStatusInProgress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo actualizar el estado"})
		return
	}

	updated, _ := h.svc.GetTask(id)
	c.JSON(http.StatusOK, updated)
}

func (h *TaskHandler) markPaused(c *gin.Context) {
	id := c.Param("id")

	err := h.svc.UpdateStatus(id, domain.TaskStatusPaused)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo pausar la tarea"})
		return
	}

	updated, _ := h.svc.GetTask(id)
	c.JSON(http.StatusOK, updated)
}

func (h *TaskHandler) reopenTask(c *gin.Context) {
	id := c.Param("id")

	err := h.svc.UpdateStatus(id, domain.TaskStatusPending)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo reabrir la tarea"})
		return
	}

	updated, _ := h.svc.GetTask(id)
	c.JSON(http.StatusOK, updated)
}

// createTask maneja la creación de una nueva tarea.
func (h *TaskHandler) createTask(c *gin.Context) {
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.svc.CreateTask(req.UserID, req.Title, req.Description, req.ProjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear la tarea"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// getTasksByUser devuelve todas las tareas pertenecientes a un usuario.
func (h *TaskHandler) getTasksByUser(c *gin.Context) {
	userID := c.Param("userID")

	tasks, err := h.svc.GetTasksByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error obteniendo tareas"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// getTask devuelve una tarea por ID.
func (h *TaskHandler) getTask(c *gin.Context) {
	id := c.Param("id")

	task, err := h.svc.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tarea no encontrada"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// updateTask actualiza los datos de una tarea.
type updateTaskRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	ProjectID   *string `json:"project_id"`
}

func (h *TaskHandler) updateTask(c *gin.Context) {
	id := c.Param("id")

	var req updateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.svc.UpdateTask(id, req.Title, req.Description, req.ProjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar tarea"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// deleteTask elimina una tarea por ID.
func (h *TaskHandler) deleteTask(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.DeleteTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al eliminar tarea"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": true})
}
