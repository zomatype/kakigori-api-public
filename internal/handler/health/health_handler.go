package health

import (
	"kakigori-api/internal/handler/response"
	"net/http"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// HealthCheck GET /health - ヘルスチェック
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": "kakigori-api",
	})
}
