package middleware

import (
	"kakigori-api/internal/middleware/cors"
	"kakigori-api/internal/middleware/logging"
	"log"
	"net/http"
)

type Logger = logging.Logger

func NewLogger(infoLogger, errorLogger *log.Logger) *Logger {
	return logging.NewLogger(infoLogger, errorLogger)
}

func CORS(config cors.Config) func(http.Handler) http.Handler {
	return cors.Middleware(config)
}

func DefaultCORSConfig() cors.Config {
	return cors.DefaultConfig()
}
