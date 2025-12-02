package config

import "os"

// Config agrupa la configuración necesaria para inicializar el servicio.
// Esto permite desacoplar la lógica de negocio de los detalles de entorno.
type Config struct {
	MongoURI      string
	MongoDatabase string
	HTTPPort      string
}

// Load construye una instancia de Config leyendo variables de entorno.
// Si alguna variable no existe, se aplica un valor por defecto seguro
// para entorno local de desarrollo.
func Load() Config {
	return Config{
		MongoURI:      getEnv("MONGO_URI", "mongodb://root:secret@localhost:27017/?authSource=admin"),
		MongoDatabase: getEnv("MONGO_DB", "pomodoro_service_db"),
		HTTPPort:      getEnv("HTTP_PORT", "8080"),
	}
}

// getEnv obtiene el valor de una variable de entorno o devuelve el valor
// por defecto en caso de que no esté definida.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
