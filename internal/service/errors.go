package service

import "errors"

// ErrInvalidStateTransition indica que se intentó aplicar una transición
// de estado que no es válida para la sesión actual.
var ErrInvalidStateTransition = errors.New("transición de estado inválida para la sesión")
