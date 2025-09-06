package menu

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	config "kakigori-api/internal/config"
	"kakigori-api/internal/domain"
	"kakigori-api/internal/middleware"
	"kakigori-api/internal/usecase"
)

// dummyMenuRepo implements repository.MenuRepository for testing
type dummyMenuRepo struct{}

func (d *dummyMenuRepo) ListMenu() ([]domain.MenuItem, error) {
	return []domain.MenuItem{{ID: "1", Name: "test", Description: "desc"}}, nil
}
func (d *dummyMenuRepo) GetMenu(id string) (*domain.MenuItem, error) {
	return &domain.MenuItem{ID: "1", Name: "test", Description: "desc"}, nil
}

func TestHandler_ListMenu(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		method   string
		storeID  string
		validID  bool
		wantCode int
		wantBody string
	}{
		{"success", http.MethodGet, "valid", true, http.StatusOK, "menu"},
		{"invalid method", http.MethodPost, "valid", true, http.StatusMethodNotAllowed, "Method not allowed"},
		{"store not found", http.MethodGet, "invalid", false, http.StatusNotFound, "Store not found"},
	}

	origGetStoreIDs := config.GetStoreIDs
	defer func() { config.GetStoreIDs = origGetStoreIDs }()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.validID {
				config.GetStoreIDs = func() []string { return []string{tt.storeID} }
			} else {
				config.GetStoreIDs = func() []string { return []string{"other"} }
			}
			logger := &middleware.Logger{}
			repo := &dummyMenuRepo{}
			h := NewHandler(usecase.NewMenuUsecase(repo), logger)

			r := httptest.NewRequest(tt.method, "/v1/stores/"+tt.storeID+"/menu", nil)
			r.SetPathValue("storeId", tt.storeID)
			w := httptest.NewRecorder()

			h.ListMenu(w, r)
			res := w.Result()
			if res.StatusCode != tt.wantCode {
				t.Errorf("got status %d, want %d", res.StatusCode, tt.wantCode)
			}
			body := w.Body.String()
			if !strings.Contains(body, tt.wantBody) {
				t.Errorf("body = %q, want %q", body, tt.wantBody)
			}
		})
	}
}
