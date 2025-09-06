package logging

import (
	"fmt"
	"kakigori-api/internal/common"
	"log"
	"net/http"
	"time"
)

// この世には２種類のログしかない。INFOとERRORだ。
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func NewLogger(infoLogger, errorLogger *log.Logger) *Logger {
	return &Logger{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
	}
}

func (l *Logger) LogInfo(message string) {
	if l.infoLogger != nil {
		l.infoLogger.Println("[INFO]", message)
	}
}

func (l *Logger) LogError(message string) {
	if l.errorLogger != nil {
		l.errorLogger.Println("[ERROR]", message)
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (l *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		clientIP := common.GetClientIP(r)

		message := fmt.Sprintf("%s %s %s %d %v",
			clientIP,
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration,
		)

		if wrapped.statusCode >= 400 {
			l.LogError(message)
		} else {
			l.LogInfo(message)
		}
	})
}
