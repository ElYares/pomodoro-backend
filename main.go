package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// main es el punto de entrada del servicio Pomodoro.
// En esta primera iteración solo exponemos un endpoint de salud (/health)
// para validar que el servidor HTTP funciona correctamente.
func main() {
	// gin.Default() crea un router con middlewares por defecto
	// (logger y recuperación de pánicos).
	router := gin.Default()

	// GET /health
	// Este endpoint permite a clientes o herramientas de monitoreo
	// comprobar que el servicio está en ejecución.
	router.GET("/health", func(c *gin.Context) {
		// Respondemos con un JSON simple y código 200 OK.
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "pomodoro-backend",
		})
	})

	// Inicia el servidor HTTP escuchando en el puerto 8080.
	// En una siguiente fase extraeremos este valor a configuración.
	if err := router.Run(":8080"); err != nil {
		// En un servicio real, este log sería capturado por un sistema de observabilidad.
		panic("no fue posible iniciar el servidor HTTP: " + err.Error())
	}
}
