package menu

import (
	"fmt"
	"kakigori-api/internal/config"
	"kakigori-api/internal/handler/response"
	"kakigori-api/internal/middleware"
	"kakigori-api/internal/usecase"
	"net/http"
)

type Handler struct {
	usecase *usecase.MenuUsecase
	logger  *middleware.Logger
}

func NewHandler(usecase *usecase.MenuUsecase, logger *middleware.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

// ListMenu GET /v1/stores/{storeId}/menu - メニュー一覧を取得
func (h *Handler) ListMenu(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	storeID := r.PathValue("storeId")
	// Store repoから取得すべきだが、環境変数で管理しているので、簡易的な実装
	if !config.IsValidStoreID(storeID) {
		h.logger.LogError(fmt.Sprintf("Store not found: %s", storeID))
		response.WriteError(w, http.StatusNotFound, "Store not found")
		return
	}

	menu := h.usecase.ListMenu()

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"menu": menu,
	})
}
