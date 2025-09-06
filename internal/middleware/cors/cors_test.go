package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if len(config.AllowOrigins) != 1 || config.AllowOrigins[0] != "*" {
		t.Error("Expected default AllowOrigins to be ['*']")
	}

	expectedMethods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodOptions,
	}
	if !stringSliceEqual(config.AllowMethods, expectedMethods) {
		t.Errorf("Expected default AllowMethods to be %v, got %v", expectedMethods, config.AllowMethods)
	}

	expectedHeaders := []string{
		"Content-Type",
	}
	if !stringSliceEqual(config.AllowHeaders, expectedHeaders) {
		t.Errorf("Expected default AllowHeaders to be %v, got %v", expectedHeaders, config.AllowHeaders)
	}

	if config.MaxAge != 86400 {
		t.Errorf("Expected default MaxAge to be 86400, got %d", config.MaxAge)
	}
}

func TestMiddleware(t *testing.T) {
	config := DefaultConfig()
	middleware := Middleware(config)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	tests := []struct {
		name           string
		method         string
		origin         string
		expectedStatus int
		expectedHeader map[string]string
	}{
		{
			name:           "Simple GET request",
			method:         "GET",
			origin:         "http://example.com",
			expectedStatus: http.StatusOK,
			expectedHeader: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		},
		{
			name:           "OPTIONS request",
			method:         "OPTIONS",
			origin:         "http://example.com",
			expectedStatus: http.StatusOK,
			expectedHeader: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
				"Access-Control-Max-Age":       "86400",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "http://example.com", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			w := httptest.NewRecorder()
			handler := middleware(nextHandler)
			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			for key, expected := range tt.expectedHeader {
				if actual := w.Header().Get(key); actual != expected {
					t.Errorf("Expected header %s to be %s, got %s", key, expected, actual)
				}
			}
		})
	}
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
