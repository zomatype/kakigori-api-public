package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_HealthCheck(t *testing.T) {
	handler := NewHandler()

	tests := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{
			name:       "正常なGETリクエスト",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
		},
		{
			name:       "不正なメソッド",
			method:     http.MethodPost,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(tt.method, "/health", nil)
			w := httptest.NewRecorder()

			handler.HealthCheck(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ステータスコードが一致しません: got %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var response map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("レスポンスのデコードに失敗: %v", err)
				}

				if status, ok := response["status"].(string); !ok || status != "healthy" {
					t.Errorf("statusフィールドが不正: got %v, want healthy", status)
				}

				if service, ok := response["service"].(string); !ok || service != "kakigori-api" {
					t.Errorf("serviceフィールドが不正: got %v, want kakigori-api", service)
				}

				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Content-Typeが不正: got %v, want application/json", contentType)
				}
			}
		})
	}
}
