package handler

import (
	"kakigori-api/internal/handler/health"
	"kakigori-api/internal/handler/menu"
	"kakigori-api/internal/handler/order"
	"kakigori-api/internal/middleware"
	"kakigori-api/internal/usecase"
	"net/http"
)

type Router struct {
	healthHandler *health.Handler
	menuHandler   *menu.Handler
	orderHandler  *order.Handler
}

func NewRouter(
	menuUsecase *usecase.MenuUsecase,
	orderUsecase *usecase.OrderUsecase,
	logger *middleware.Logger,
) *Router {
	return &Router{
		healthHandler: health.NewHandler(),
		menuHandler:   menu.NewHandler(menuUsecase, logger),
		orderHandler:  order.NewHandler(orderUsecase, logger),
	}
}

func (r *Router) Setup(mux *http.ServeMux) {
	// ヘルスチェック
	mux.HandleFunc("GET /health", r.healthHandler.HealthCheck)

	// v1 API
	mux.HandleFunc("GET /v1/stores/{storeId}/menu", r.menuHandler.ListMenu)
	mux.HandleFunc("POST /v1/stores/{storeId}/orders", r.orderHandler.CreateOrder)
	mux.HandleFunc("GET /v1/stores/{storeId}/orders", r.orderHandler.ListOrders)
	mux.HandleFunc("GET /v1/stores/{storeId}/orders/{orderId}", r.orderHandler.GetOrder)
	mux.HandleFunc("POST /v1/stores/{storeId}/orders/{orderId}/waiting-pickup", r.orderHandler.WaitingPickupOrder)
	mux.HandleFunc("POST /v1/stores/{storeId}/orders/{orderId}/complete", r.orderHandler.CompleteOrder)
}
