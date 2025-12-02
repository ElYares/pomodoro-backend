package http

import (
	"net/http"

	"pomodoro-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// SessionHandler expone endpoints HTTP para operar sobre sesiones
// Pomodoro consumiendo la lógica de negocio del SessionService.
type SessionHandler struct {
	svc *service.SessionService
}

// NewSessionHandler construye una nueva instancia de SessionHandler.
func NewSessionHandler(svc *service.SessionService) *SessionHandler {
	return &SessionHandler{
		svc: svc,
	}
}

// RegisterRoutes registra los endpoints asociados a sesiones Pomodoro
// sobre el grupo de rutas recibido.
func (h *SessionHandler) RegisterRoutes(rg *gin.RouterGroup) {
	sessions := rg.Group("/sessions")
	{
		sessions.POST("", h.createSession)
		sessions.PATCH("/:id/pause", h.pauseSession)
		sessions.PATCH("/:id/resume", h.resumeSession)
		sessions.PATCH("/:id/finish", h.finishSession)
	}
}

// createSessionRequest define el cuerpo esperado para la creación
// de una nueva sesión Pomodoro.
type createSessionRequest struct {
	UserID       string  `json:"user_id" binding:"required"`
	ProjectID    *string `json:"project_id"`
	TaskID       *string `json:"task_id"`
	FocusMinutes int     `json:"focus_minutes" binding:"required,min=1,max=120"`
	BreakMinutes int     `json:"break_minutes" binding:"required,min=0,max=60"`
}

// createSession maneja la creación de una nueva sesión Pomodoro.
func (h *SessionHandler) createSession(c *gin.Context) {
	var req createSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "payload inválido",
			"detail": err.Error(),
		})
		return
	}

	session, err := h.svc.CreateAndStartSession(
		req.UserID,
		req.ProjectID,
		req.TaskID,
		req.FocusMinutes,
		req.BreakMinutes,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo crear la sesión"})
		return
	}

	c.JSON(http.StatusCreated, session)
}

// pauseSession cambia el estado de la sesión a PAUSED.
func (h *SessionHandler) pauseSession(c *gin.Context) {
	id := c.Param("id")

	session, err := h.svc.PauseSession(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrInvalidStateTransition {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

// resumeSession cambia el estado de la sesión a RUNNING.
func (h *SessionHandler) resumeSession(c *gin.Context) {
	id := c.Param("id")

	session, err := h.svc.ResumeSession(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrInvalidStateTransition {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

// finishSession marca la sesión como finalizada.
func (h *SessionHandler) finishSession(c *gin.Context) {
	id := c.Param("id")

	session, err := h.svc.FinishSession(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrInvalidStateTransition {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}
