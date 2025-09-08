package order

import (
	"encoding/json"
	"fmt"
	"kakigori-api/internal/config"
	"kakigori-api/internal/handler/response"
	"kakigori-api/internal/middleware"
	"kakigori-api/internal/usecase"
	"net/http"
)

type Handler struct {
	usecase *usecase.OrderUsecase
	logger  *middleware.Logger
}

func NewHandler(usecase *usecase.OrderUsecase, logger *middleware.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

// CreateOrder POST /v1/stores/{storeId}/orders - 注文を作成
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.LogError(fmt.Sprintf("Method not allowed: %s", r.Method))
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

	var req struct {
		MenuItemID string `json:"menu_item_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.LogError(fmt.Sprintf("Invalid request body: %v", err))
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	order, err := h.usecase.CreateOrder(storeID, req.MenuItemID)
	if err != nil {
		h.logger.LogError(fmt.Sprintf("Failed to create order for store %s: %v", storeID, err))
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.LogInfo(fmt.Sprintf("Order created: %s for store %s", order.ID, storeID))
	response.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"id":           order.ID,
		"menu_item_id": order.MenuItemID,
		"menu_name":    order.MenuName,
		"order_number": order.OrderNumber,
		"status":       order.Status,
	})
}

// ListOrders GET /v1/stores/{storeId}/orders - 注文一覧を取得
func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.LogError(fmt.Sprintf("Method not allowed: %s", r.Method))
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	storeID := r.PathValue("storeId")
	if !config.IsValidStoreID(storeID) {
		h.logger.LogError(fmt.Sprintf("Store not found: %s", storeID))
		response.WriteError(w, http.StatusNotFound, "Store not found")
		return
	}

	orders, err := h.usecase.ListOrders(storeID)
	if err != nil {
		h.logger.LogError(fmt.Sprintf("Failed to list orders for store %s: %v", storeID, err))
		response.WriteError(w, http.StatusInternalServerError, "Failed to list orders")
		return
	}

	result := make([]map[string]interface{}, 0, len(orders))
	for _, o := range orders {
		result = append(result, map[string]interface{}{
			"id":           o.ID,
			"menu_item_id": o.MenuItemID,
			"menu_name":    o.MenuName,
			"order_number": o.OrderNumber,
			"status":       o.Status,
		})
	}

	h.logger.LogInfo(fmt.Sprintf("Retrieved %d orders for store %s", len(result), storeID))
	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"orders": result,
	})
}

// GetOrders GET /v1/stores/{storeId}/orders/{orderId} - 注文を取得
// 本来は認証を入れるべきだが、今回は省略
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.LogError(fmt.Sprintf("Method not allowed: %s", r.Method))
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	storeID := r.PathValue("storeId")
	if !config.IsValidStoreID(storeID) {
		h.logger.LogError(fmt.Sprintf("Store not found: %s", storeID))
		response.WriteError(w, http.StatusNotFound, "Store not found")
		return
	}

	orderID := r.PathValue("orderId")
	// 注文ID形式チェック: store-id-order-number
	if len(orderID) == 0 || len(storeID) == 0 || len(orderID) <= len(storeID)+1 || orderID[:len(storeID)] != storeID || orderID[len(storeID)] != '-' {
		h.logger.LogError(fmt.Sprintf("Invalid order ID format: %s", orderID))
		response.WriteError(w, http.StatusBadRequest, "Invalid order ID format")
		return
	}

	order, err := h.usecase.GetOrder(storeID, orderID)
	if err != nil {
		if err.Error() == "not found" {
			h.logger.LogError(fmt.Sprintf("Order not found: %s", orderID))
			response.WriteError(w, http.StatusNotFound, "Order not found")
			return
		}
		h.logger.LogError(fmt.Sprintf("Failed to get order: %v", err))
		response.WriteError(w, http.StatusInternalServerError, "Failed to get order")
		return
	}

	h.logger.LogInfo(fmt.Sprintf("Retrieved order %s for store %s", orderID, storeID))
	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id":           order.ID,
		"menu_item_id": order.MenuItemID,
		"menu_name":    order.MenuName,
		"order_number": order.OrderNumber,
		"status":       order.Status,
	})
}

// WaitingPickupOrder POST /v1/stores/{storeId}/orders/{orderId}/waiting-pickup - 注文を受付済み状態に更新
func (h *Handler) WaitingPickupOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.LogError(fmt.Sprintf("Method not allowed: %s", r.Method))
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	storeID := r.PathValue("storeId")
	if !config.IsValidStoreID(storeID) {
		h.logger.LogError(fmt.Sprintf("Store not found: %s", storeID))
		response.WriteError(w, http.StatusNotFound, "Store not found")
		return
	}

	orderID := r.PathValue("orderId")
	// 注文ID形式チェック: store-id-order-number
	if len(orderID) == 0 || len(storeID) == 0 || len(orderID) <= len(storeID)+1 || orderID[:len(storeID)] != storeID || orderID[len(storeID)] != '-' {
		h.logger.LogError(fmt.Sprintf("Invalid order ID format: %s", orderID))
		response.WriteError(w, http.StatusBadRequest, "Invalid order ID format")
		return
	}

	order, err := h.usecase.WaitingPickup(storeID, orderID)
	if err != nil {
		if err.Error() == "not found" {
			h.logger.LogError(fmt.Sprintf("Order not found: %s", orderID))
			response.WriteError(w, http.StatusNotFound, "Order not found")
			return
		}
		h.logger.LogError(fmt.Sprintf("Failed to update order status: %v", err))
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.LogInfo(fmt.Sprintf("Order %s updated to waiting pickup for store %s", orderID, storeID))
	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id":           order.ID,
		"menu_item_id": order.MenuItemID,
		"menu_name":    order.MenuName,
		"order_number": order.OrderNumber,
		"status":       order.Status,
	})
}

// CompleteOrder POST /v1/stores/{storeId}/orders/{orderId}/complete - 注文を完了状態に更新
func (h *Handler) CompleteOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.LogError(fmt.Sprintf("Method not allowed: %s", r.Method))
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	storeID := r.PathValue("storeId")
	if !config.IsValidStoreID(storeID) {
		h.logger.LogError(fmt.Sprintf("Store not found: %s", storeID))
		response.WriteError(w, http.StatusNotFound, "Store not found")
		return
	}

	orderID := r.PathValue("orderId")
	// 注文ID形式チェック: store-id-order-number
	if len(orderID) == 0 || len(storeID) == 0 || len(orderID) <= len(storeID)+1 || orderID[:len(storeID)] != storeID || orderID[len(storeID)] != '-' {
		h.logger.LogError(fmt.Sprintf("Invalid order ID format: %s", orderID))
		response.WriteError(w, http.StatusBadRequest, "Invalid order ID format")
		return
	}

	order, err := h.usecase.Complete(storeID, orderID)
	if err != nil {
		if err.Error() == "not found" {
			h.logger.LogError(fmt.Sprintf("Order not found: %s", orderID))
			response.WriteError(w, http.StatusNotFound, "Order not found")
			return
		}
		h.logger.LogError(fmt.Sprintf("Failed to update order status: %v", err))
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.LogInfo(fmt.Sprintf("Order %s completed for store %s", orderID, storeID))
	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id":           order.ID,
		"menu_item_id": order.MenuItemID,
		"menu_name":    order.MenuName,
		"order_number": order.OrderNumber,
		"status":       order.Status,
	})
}
