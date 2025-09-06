package logging

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	var infoBuf, errorBuf bytes.Buffer
	infoLogger := log.New(&infoBuf, "", 0)
	errorLogger := log.New(&errorBuf, "", 0)

	logger := NewLogger(infoLogger, errorLogger)

	if logger.infoLogger != infoLogger {
		t.Error("Expected infoLogger to be set correctly")
	}

	if logger.errorLogger != errorLogger {
		t.Error("Expected errorLogger to be set correctly")
	}
}

func TestLogger_LogInfo(t *testing.T) {
	var buf bytes.Buffer
	infoLogger := log.New(&buf, "", 0)
	logger := NewLogger(infoLogger, nil)

	message := "test info message"
	logger.LogInfo(message)

	output := buf.String()
	if !strings.Contains(output, "[INFO]") {
		t.Errorf("Expected output to contain '[INFO]', got: %s", output)
	}

	if !strings.Contains(output, message) {
		t.Errorf("Expected output to contain message '%s', got: %s", message, output)
	}
}

func TestLogger_LogError(t *testing.T) {
	var buf bytes.Buffer
	errorLogger := log.New(&buf, "", 0)
	logger := NewLogger(nil, errorLogger)

	message := "test error message"
	logger.LogError(message)

	output := buf.String()
	if !strings.Contains(output, "[ERROR]") {
		t.Errorf("Expected output to contain '[ERROR]', got: %s", output)
	}

	if !strings.Contains(output, message) {
		t.Errorf("Expected output to contain message '%s', got: %s", message, output)
	}
}

func TestLogger_Middleware(t *testing.T) {
	var infoBuf, errorBuf bytes.Buffer
	infoLogger := log.New(&infoBuf, "", 0)
	errorLogger := log.New(&errorBuf, "", 0)
	logger := NewLogger(infoLogger, errorLogger)

	tests := []struct {
		name          string
		handler       func(w http.ResponseWriter, r *http.Request)
		expectedInfo  bool
		expectedError bool
	}{
		{
			name: "Successful request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			},
			expectedInfo:  true,
			expectedError: false,
		},
		{
			name: "Error request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error"))
			},
			expectedInfo:  false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			infoBuf.Reset()
			errorBuf.Reset()

			handler := http.HandlerFunc(tt.handler)
			middleware := logger.Middleware(handler)

			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.100:12345"
			w := httptest.NewRecorder()

			middleware.ServeHTTP(w, req)

			infoLog := infoBuf.String()
			errorLog := errorBuf.String()

			if tt.expectedInfo && infoLog == "" {
				t.Error("Expected info log, got nothing")
			}
			if !tt.expectedInfo && infoLog != "" {
				t.Error("Expected no info log, got log")
			}
			if tt.expectedError && errorLog == "" {
				t.Error("Expected error log, got nothing")
			}
			if !tt.expectedError && errorLog != "" {
				t.Error("Expected no error log, got log")
			}
		})
	}
}

func TestLogger_Middleware_Duration(t *testing.T) {
	var buf bytes.Buffer
	infoLogger := log.New(&buf, "", 0)
	logger := NewLogger(infoLogger, nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	middleware := logger.Middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.100:12345"
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	output := buf.String()
	if !strings.Contains(output, "ms") {
		t.Error("Expected duration to be logged in milliseconds")
	}
}
