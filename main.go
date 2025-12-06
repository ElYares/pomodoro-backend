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

func main() {

	// Cargar .env si existe (no obligatorio en Docker)
	_ = godotenv.Load()

	cfg := config.Load()

	// ---------------------------
	// Conexión a MongoDB
	// ---------------------------

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("error al conectar a MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB no responde al ping: %v", err)
	}

	db := client.Database(cfg.MongoDatabase)

	// ---------------------------
	// Inyección de Repositorios
	// ---------------------------

	sessionRepo := repository.NewMongoSessionRepository(db)
	taskRepo := repository.NewMongoTaskRepository(db)
	cycleRepo := repository.NewMongoCycleRepository(db)

	// ---------------------------
	// Inyección de Servicios
	// ---------------------------

	sessionService := service.NewSessionService(sessionRepo, taskRepo) // <- corregido
	taskService := service.NewTaskService(taskRepo)
	cycleService := service.NewCycleService(cycleRepo)

	// ---------------------------
	// Inyección de Handlers
	// ---------------------------

	sessionHandler := httphandler.NewSessionHandler(sessionService)
	taskHandler := httphandler.NewTaskHandler(taskService)
	_ = cycleService // Pendiente: aún no tienes endpoints para cycles

	// ---------------------------
	// Router
	// ---------------------------

	router := gin.Default()

	// Middleware CORS
	router.Use(func(c *gin.Context) {
		// Permitir llamadas desde tu frontend en http://localhost:3000
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Health Check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "pomodoro-backend",
		})
	})

	api := router.Group("/api/v1")
	{
		sessionHandler.RegisterRoutes(api)
		taskHandler.RegisterRoutes(api)
	}

	// ---------------------------
	// Servidor HTTP
	// ---------------------------

	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	log.Printf("Pomodoro backend escuchando en puerto %s", cfg.HTTPPort)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error al iniciar servidor: %v", err)
	}
}
