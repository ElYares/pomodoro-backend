package service

import "errors"

// Errores estándar del dominio de servicio
var (
	// Sesiones
	ErrSessionNotFound        = errors.New("session not found")
	ErrInvalidState           = errors.New("invalid session state")
	ErrInvalidStateTransition = errors.New("invalid state transition")

	// Tareas
	ErrTaskNotFound = errors.New("task not found")
)
