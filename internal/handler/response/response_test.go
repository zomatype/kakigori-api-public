package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name       string
		data       interface{}
		statusCode int
	}{
		{
			name: "正常なJSONデータ",
			data: map[string]interface{}{
				"message": "success",
				"code":    200,
			},
			statusCode: http.StatusOK,
		},
		{
			name:       "空のデータ",
			data:       map[string]interface{}{},
			statusCode: http.StatusOK,
		},
		{
			name: "エラーステータス",
			data: map[string]interface{}{
				"error": "not found",
			},
			statusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteJSON(w, tt.statusCode, tt.data)

			if w.Code != tt.statusCode {
				t.Errorf("WriteJSON() status = %v, want %v", w.Code, tt.statusCode)
			}

			if w.Code == tt.statusCode {
				return
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("WriteJSON() Content-Type = %v, want application/json", contentType)
			}

			var response map[string]interface{}
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			for key, expectedValue := range tt.data.(map[string]interface{}) {
				if actualValue := response[key]; actualValue != expectedValue {
					t.Errorf("WriteJSON() %s = %v, want %v", key, actualValue, expectedValue)
				}
			}
		})
	}
}

func TestWriteError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
	}{
		{
			name:       "Bad Request",
			statusCode: http.StatusBadRequest,
			message:    "Invalid request",
		},
		{
			name:       "Not Found",
			statusCode: http.StatusNotFound,
			message:    "Resource not found",
		},
		{
			name:       "Internal Server Error",
			statusCode: http.StatusInternalServerError,
			message:    "Server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteError(w, tt.statusCode, tt.message)

			if w.Code != tt.statusCode {
				t.Errorf("WriteError() status = %v, want %v", w.Code, tt.statusCode)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("WriteError() Content-Type = %v, want application/json", contentType)
			}

			var response ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if response.Message != tt.message {
				t.Errorf("WriteError() message = %v, want %v", response.Message, tt.message)
			}

			expectedError := http.StatusText(tt.statusCode)
			if response.Error != expectedError {
				t.Errorf("WriteError() error = %v, want %v", response.Error, expectedError)
			}
		})
	}
}
