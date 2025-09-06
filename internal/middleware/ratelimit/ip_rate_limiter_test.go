package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIpRateLimiter_SecondWindow(t *testing.T) {
	limiter := NewIpRateLimiter(2, 1*time.Second, 1*time.Second)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := limiter.Middleware(handler)

	req := httptest.NewRequest("GET", "/v1/stores/test-store-001/orders", nil)
	req.RemoteAddr = "192.168.1.100:12345"

	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("1st request: expected 200, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("2nd request: expected 200, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("3rd request: expected 403, got %d", w.Code)
	}

	time.Sleep(1100 * time.Millisecond)
	w = httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("After reset: expected 200, got %d", w.Code)
	}
}
