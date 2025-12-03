package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"pomodoro-backend/internal/config"
	"pomodoro-backend/internal/repository"
	"pomodoro-backend/internal/service"
	httphandler "pomodoro-backend/internal/transport/http"
)

// main es el punto de entrada del servicio Pomodoro.
// Su responsabilidad es orquestar la inicialización de configuración,
// dependencias e infraestructura de red (servidor HTTP).
func main() {
	// Cargar variables de entorno
	_ = godotenv.Load()

	// Cargar configuración centralizada
	cfg := config.Load()

	// Inicializar cliente MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("error al conectar a MongoDB: %v", err)
	}

	// Ping para validar conectividad
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB no responde al ping: %v", err)
	}

	// Seleccionar base de datos
	db := client.Database(cfg.MongoDatabase)

	// ---------------------------
	// Inyección de dependencias
	// ---------------------------

	// Sessions
	sessionRepo := repository.NewMongoSessionRepository(db)
	sessionService := service.NewSessionService(sessionRepo)
	sessionHandler := httphandler.NewSessionHandler(sessionService)

	// Tasks
	taskRepo := repository.NewMongoTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := httphandler.NewTaskHandler(taskService)

	// ---------------------------
	// Configuración del router
	// ---------------------------

	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "pomodoro-backend",
		})
	})

	// API versionada
	api := router.Group("/api/v1")
	{
		// Registro de módulos
		sessionHandler.RegisterRoutes(api)
		taskHandler.RegisterRoutes(api)
	}

	// Configuración del servidor HTTP
	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	log.Printf("Pomodoro Service escuchando en puerto %s", cfg.HTTPPort)

	// Iniciar servidor
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error al iniciar el servidor HTTP: %v", err)
	}
}
