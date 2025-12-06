package service

import (
	"time"
)

// GenerateID genera un ID único basado en timestamp.
// Puedes reemplazarlo luego por UUID v4 sin romper nada.
func GenerateID() string {
	return time.Now().Format("20060102150405.000000")
}
