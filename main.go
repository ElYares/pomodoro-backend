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
	// Carga variables de entorno desde el archivo .env, si existe.
	// En entornos productivos, lo habitual es depender solo de variables
	// de entorno del sistema u orquestador (Docker, Kubernetes, etc.).
	_ = godotenv.Load()

	// Construye la configuración de la aplicación a partir de las variables
	// de entorno y valores por defecto seguros.
	cfg := config.Load()

	// Inicializa el cliente de MongoDB utilizando la URI configurada.
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("error al conectar a MongoDB: %v", err)
	}

	// Verifica la conectividad contra la base de datos mediante un ping.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB no responde al ping: %v", err)
	}

	// Obtiene una referencia a la base de datos específica para este servicio.
	db := client.Database(cfg.MongoDatabase)

	// Inyección de dependencias:
	// Repositorio (persistencia) → Servicio (negocio) → Handler HTTP (transporte).
	sessionRepo := repository.NewMongoSessionRepository(db)
	sessionService := service.NewSessionService(sessionRepo)
	sessionHandler := httphandler.NewSessionHandler(sessionService)

	// gin.Default configura un router HTTP con middlewares de logging y
	// recuperación de pánicos activados por defecto.
	router := gin.Default()

	// Endpoint de salud que permite a herramientas de monitoreo o balanceadores
	// verificar que el servicio está operativo.
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "pomodoro-backend",
		})
	})

	// Grupo base de la API versionado. A partir de aquí se registran los
	// endpoints específicos del dominio Pomodoro.
	api := router.Group("/api/v1")
	{
		sessionHandler.RegisterRoutes(api)
	}

	// Configuración del servidor HTTP. Se utiliza http.Server para tener
	// mayor control sobre timeouts y apagado ordenado en futuras mejoras.
	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	log.Printf("Pomodoro Service escuchando en puerto %s", cfg.HTTPPort)

	// Inicia el servidor HTTP. Si falla al arrancar, se registra el error
	// y se termina el proceso con un código de salida no exitoso.
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error al iniciar el servidor HTTP: %v", err)
	}
}
