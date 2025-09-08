package order

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kakigori-api/internal/config"
	"kakigori-api/internal/domain"
	"kakigori-api/internal/middleware"
	"kakigori-api/internal/usecase"
)

type dummyMenuRepo struct{}

func (d *dummyMenuRepo) ListMenu() ([]domain.MenuItem, error) { return nil, nil }
func (d *dummyMenuRepo) GetMenu(id string) (*domain.MenuItem, error) {
	return &domain.MenuItem{ID: "m1", Name: "test", Description: "desc"}, nil
}

type dummyOrderRepo struct {
	order     *domain.Order
	err       error
	deleteErr error
}

func (d *dummyOrderRepo) Create(storeID string, order *domain.Order) error { return d.err }
func (d *dummyOrderRepo) GetByID(storeID, orderID string) (*domain.Order, error) {
	return d.order, d.err
}
func (d *dummyOrderRepo) GetAllByStore(storeID string) ([]*domain.Order, error) {
	if d.err != nil {
		return nil, d.err
	}
	return []*domain.Order{d.order}, nil
}
func (d *dummyOrderRepo) Update(storeID string, order *domain.Order) error { return d.err }
func (d *dummyOrderRepo) Delete(storeID, orderID string) error {
	if d.deleteErr != nil {
		return d.deleteErr
	}
	return d.err
}
func (d *dummyOrderRepo) GetNextOrderNumber(storeID string) (int, error)       { return 1, nil }
func (d *dummyOrderRepo) GetOrderCount(storeID string) (int, error)            { return 1, nil }
func (d *dummyOrderRepo) GetOldestOrder(storeID string) (*domain.Order, error) { return d.order, d.err }

type dummyStoreRepo struct{ err error }

func (d *dummyStoreRepo) Create(s *domain.Store) error { return d.err }
func (d *dummyStoreRepo) GetByID(id string) (*domain.Store, error) {
	return &domain.Store{ID: id, MaxOrders: 10}, d.err
}

func TestOrderHandler(t *testing.T) {
	t.Parallel()
	order := &domain.Order{ID: "store-1-1", MenuItemID: "m1", MenuName: "test", Status: domain.OrderStatusPending, OrderNumber: 1}
	validStoreID := "store-1"
	validOrderID := "store-1-1"
	testCases := []struct {
		name     string
		handler  string
		method   string
		storeID  string
		orderID  string
		body     string
		validID  bool
		order    *domain.Order
		errGet   error
		errDel   error
		wantCode int
		wantBody string
	}{
		{"GetOrder success", "GetOrder", http.MethodGet, validStoreID, validOrderID, "", true, order, nil, nil, http.StatusOK, `{"id":"store-1-1","menu_item_id":"m1","menu_name":"test","order_number":1,"status":"pending"}`},
		{"GetOrder invalid method", "GetOrder", http.MethodPost, validStoreID, validOrderID, "", true, order, nil, nil, http.StatusMethodNotAllowed, "Method not allowed"},
		{"GetOrder store not found", "GetOrder", http.MethodGet, "invalid", validOrderID, "", false, order, nil, nil, http.StatusNotFound, "Store not found"},
		{"GetOrder invalid order id format", "GetOrder", http.MethodGet, validStoreID, "badid", "", true, order, nil, nil, http.StatusBadRequest, "Invalid order ID format"},
		{"GetOrder order not found", "GetOrder", http.MethodGet, validStoreID, validOrderID, "", true, nil, errors.New("not found"), nil, http.StatusNotFound, "Order not found"},

		{"CreateOrder success", "CreateOrder", http.MethodPost, validStoreID, "", `{"menu_item_id":"m1"}`, true, order, nil, nil, http.StatusCreated, "id"},
		{"CreateOrder invalid method", "CreateOrder", http.MethodGet, validStoreID, "", `{"menu_item_id":"m1"}`, true, order, nil, nil, http.StatusMethodNotAllowed, "Method not allowed"},
		{"CreateOrder store not found", "CreateOrder", http.MethodPost, "invalid", "", `{"menu_item_id":"m1"}`, false, order, nil, nil, http.StatusNotFound, "Store not found"},
		{"CreateOrder invalid body", "CreateOrder", http.MethodPost, validStoreID, "", `invalid`, true, order, nil, nil, http.StatusBadRequest, "Invalid request body"},
		{"CreateOrder order create error", "CreateOrder", http.MethodPost, validStoreID, "", `{"menu_item_id":"m1"}`, true, nil, errors.New("fail"), nil, http.StatusBadRequest, "fail"},

		{"ListOrders success", "ListOrders", http.MethodGet, validStoreID, "", "", true, order, nil, nil, http.StatusOK, "orders"},
		{"ListOrders invalid method", "ListOrders", http.MethodPost, validStoreID, "", "", true, order, nil, nil, http.StatusMethodNotAllowed, "Method not allowed"},
		{"ListOrders store not found", "ListOrders", http.MethodGet, "invalid", "", "", false, order, nil, nil, http.StatusNotFound, "Store not found"},
		{"ListOrders list error", "ListOrders", http.MethodGet, validStoreID, "", "", true, nil, errors.New("fail"), nil, http.StatusInternalServerError, "Failed to list orders"},

		{"WaitingPickupOrder success", "WaitingPickupOrder", http.MethodPost, validStoreID, validOrderID, "", true, &domain.Order{ID: "store-1-1", MenuItemID: "m1", MenuName: "test", Status: domain.OrderStatusPending, OrderNumber: 1}, nil, nil, http.StatusOK, "status"},
		{"WaitingPickupOrder invalid method", "WaitingPickupOrder", http.MethodGet, validStoreID, validOrderID, "", true, order, nil, nil, http.StatusMethodNotAllowed, "Method not allowed"},
		{"WaitingPickupOrder store not found", "WaitingPickupOrder", http.MethodPost, "invalid", validOrderID, "", false, order, nil, nil, http.StatusNotFound, "Store not found"},
		{"WaitingPickupOrder invalid order id format", "WaitingPickupOrder", http.MethodPost, validStoreID, "badid", "", true, order, nil, nil, http.StatusBadRequest, "Invalid order ID format"},
		{"WaitingPickupOrder order not found", "WaitingPickupOrder", http.MethodPost, validStoreID, validOrderID, "", true, nil, errors.New("not found"), nil, http.StatusNotFound, "Order not found"},
		{"WaitingPickupOrder invalid status", "WaitingPickupOrder", http.MethodPost, validStoreID, validOrderID, "", true, &domain.Order{ID: "store-1-1", MenuItemID: "m1", MenuName: "test", Status: domain.OrderStatusCompleted, OrderNumber: 1}, nil, nil, http.StatusBadRequest, "order status is not pending"},

		{"CompleteOrder success", "CompleteOrder", http.MethodPost, validStoreID, validOrderID, "", true, &domain.Order{ID: "store-1-1", MenuItemID: "m1", MenuName: "test", Status: domain.OrderStatusWaitingPickup, OrderNumber: 1}, nil, nil, http.StatusOK, "status"},
		{"CompleteOrder invalid method", "CompleteOrder", http.MethodGet, validStoreID, validOrderID, "", true, order, nil, nil, http.StatusMethodNotAllowed, "Method not allowed"},
		{"CompleteOrder store not found", "CompleteOrder", http.MethodPost, "invalid", validOrderID, "", false, order, nil, nil, http.StatusNotFound, "Store not found"},
		{"CompleteOrder invalid order id format", "CompleteOrder", http.MethodPost, validStoreID, "badid", "", true, order, nil, nil, http.StatusBadRequest, "Invalid order ID format"},
		{"CompleteOrder order not found", "CompleteOrder", http.MethodPost, validStoreID, validOrderID, "", true, nil, errors.New("not found"), nil, http.StatusNotFound, "Order not found"},
		{"CompleteOrder invalid status", "CompleteOrder", http.MethodPost, validStoreID, validOrderID, "", true, &domain.Order{ID: "store-1-1", MenuItemID: "m1", MenuName: "test", Status: domain.OrderStatusPending, OrderNumber: 1}, nil, nil, http.StatusBadRequest, "order status is not waiting pickup"},
	}

	origGetStoreIDs := config.GetStoreIDs
	defer func() { config.GetStoreIDs = origGetStoreIDs }()

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.validID {
				config.GetStoreIDs = func() []string { return []string{tt.storeID} }
			} else {
				config.GetStoreIDs = func() []string { return []string{"other"} }
			}
			logger := &middleware.Logger{}
			menuRepo := &dummyMenuRepo{}
			orderRepo := &dummyOrderRepo{order: tt.order, err: tt.errGet, deleteErr: tt.errDel}
			storeRepo := &dummyStoreRepo{}
			if tt.validID {
				storeRepo.Create(&domain.Store{ID: tt.storeID, MaxOrders: 10})
			}
			uc := usecase.NewOrderUsecase(menuRepo, orderRepo, storeRepo)
			h := NewHandler(uc, logger)

			var r *http.Request
			var w *httptest.ResponseRecorder

			switch tt.handler {
			case "GetOrder":
				r = httptest.NewRequest(tt.method, "/v1/stores/"+tt.storeID+"/orders/"+tt.orderID, nil)
				r.SetPathValue("storeId", tt.storeID)
				r.SetPathValue("orderId", tt.orderID)
				w = httptest.NewRecorder()
				h.GetOrder(w, r)
			case "CreateOrder":
				r = httptest.NewRequest(tt.method, "/v1/stores/"+tt.storeID+"/orders", strings.NewReader(tt.body))
				r.SetPathValue("storeId", tt.storeID)
				w = httptest.NewRecorder()
				h.CreateOrder(w, r)
			case "ListOrders":
				r = httptest.NewRequest(tt.method, "/v1/stores/"+tt.storeID+"/orders", nil)
				r.SetPathValue("storeId", tt.storeID)
				w = httptest.NewRecorder()
				h.ListOrders(w, r)
			case "WaitingPickupOrder":
				r = httptest.NewRequest(tt.method, "/v1/stores/"+tt.storeID+"/orders/"+tt.orderID+"/waiting-pickup", nil)
				r.SetPathValue("storeId", tt.storeID)
				r.SetPathValue("orderId", tt.orderID)
				w = httptest.NewRecorder()
				h.WaitingPickupOrder(w, r)
			case "CompleteOrder":
				r = httptest.NewRequest(tt.method, "/v1/stores/"+tt.storeID+"/orders/"+tt.orderID+"/complete", nil)
				r.SetPathValue("storeId", tt.storeID)
				r.SetPathValue("orderId", tt.orderID)
				w = httptest.NewRecorder()
				h.CompleteOrder(w, r)
			}

			res := w.Result()
			if res.StatusCode != tt.wantCode {
				t.Errorf("got status %d, want %d", res.StatusCode, tt.wantCode)
			}
			body := w.Body.String()
			if tt.wantBody != "" && !strings.Contains(body, tt.wantBody) {
				t.Errorf("body = %q, want %q", body, tt.wantBody)
			}
		})
	}
}
