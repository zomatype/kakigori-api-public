package usecase

import (
	"kakigori-api/internal/domain"
	"testing"
	"time"
)

type dummyOrderRepo struct {
	orders    map[string]*domain.Order
	createErr error
	getErr    error
	listErr   error
	deleteErr error
}

func (d *dummyOrderRepo) Create(storeID string, order *domain.Order) error { return d.createErr }
func (d *dummyOrderRepo) GetByID(storeID, orderID string) (*domain.Order, error) {
	if d.getErr != nil {
		return nil, d.getErr
	}
	return d.orders[orderID], nil
}
func (d *dummyOrderRepo) GetAllByStore(storeID string) ([]*domain.Order, error) {
	if d.listErr != nil {
		return nil, d.listErr
	}
	var res []*domain.Order
	for _, o := range d.orders {
		res = append(res, o)
	}
	return res, nil
}
func (d *dummyOrderRepo) Update(storeID string, order *domain.Order) error { return nil }
func (d *dummyOrderRepo) Delete(storeID, orderID string) error             { return d.deleteErr }
func (d *dummyOrderRepo) GetNextOrderNumber(storeID string) (int, error)   { return 1, nil }
func (d *dummyOrderRepo) GetOrderCount(storeID string) (int, error)        { return len(d.orders), nil }
func (d *dummyOrderRepo) GetOldestOrder(storeID string) (*domain.Order, error) {
	for _, o := range d.orders {
		return o, nil
	}
	return nil, nil
}

type dummyStoreRepo struct {
	stores map[string]*domain.Store
	getErr error
}

func (d *dummyStoreRepo) Create(s *domain.Store) error { return nil }
func (d *dummyStoreRepo) GetByID(id string) (*domain.Store, error) {
	if d.getErr != nil {
		return nil, d.getErr
	}
	return d.stores[id], nil
}

func TestOrderUsecase(t *testing.T) {
	t.Parallel()
	menuRepo := &dummyMenuRepo{}
	orderRepo := &dummyOrderRepo{orders: map[string]*domain.Order{"order1": {ID: "order1", MenuItemID: "id1", MenuName: "name1", Status: domain.OrderStatusPending, OrderNumber: 1, CreatedAt: time.Now()}}}
	storeRepo := &dummyStoreRepo{stores: map[string]*domain.Store{"store1": {ID: "store1", MaxOrders: 10}}}
	uc := NewOrderUsecase(menuRepo, orderRepo, storeRepo)
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{"CreateOrder success", func(t *testing.T) {
			order, err := uc.CreateOrder("store1", "id1")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if order.MenuItemID != "id1" {
				t.Errorf("expected id1, got %s", order.MenuItemID)
			}
		}},
		{"ListOrders success", func(t *testing.T) {
			// Create a new order repo with a single order
			orderRepo := &dummyOrderRepo{orders: map[string]*domain.Order{
				"order1": {
					ID:          "order1",
					MenuItemID:  "id1",
					MenuName:    "name1",
					Status:      domain.OrderStatusPending,
					OrderNumber: 1,
					CreatedAt:   time.Now(),
				},
			}}
			uc := NewOrderUsecase(menuRepo, orderRepo, storeRepo)

			orders, err := uc.ListOrders("store1")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(orders) != 1 {
				t.Errorf("expected 1 order, got %d", len(orders))
			}
		}},
		{"GetOrder success", func(t *testing.T) {
			order, err := uc.GetOrder("store1", "order1")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if order.ID != "order1" {
				t.Errorf("expected order1, got %s", order.ID)
			}
		}},
		{"WaitingPickup success", func(t *testing.T) {
			order, err := uc.WaitingPickup("store1", "order1")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if order.Status != domain.OrderStatusWaitingPickup {
				t.Errorf("expected status %s, got %s", domain.OrderStatusWaitingPickup, order.Status)
			}
		}},
		{"Complete success", func(t *testing.T) {
			// Create a new order repo with a waiting pickup order
			orderRepo := &dummyOrderRepo{orders: map[string]*domain.Order{
				"order1": {
					ID:          "order1",
					MenuItemID:  "id1",
					MenuName:    "name1",
					Status:      domain.OrderStatusWaitingPickup,
					OrderNumber: 1,
					CreatedAt:   time.Now(),
				},
			}}
			uc := NewOrderUsecase(menuRepo, orderRepo, storeRepo)

			order, err := uc.Complete("store1", "order1")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if order.Status != domain.OrderStatusCompleted {
				t.Errorf("expected status %s, got %s", domain.OrderStatusCompleted, order.Status)
			}
		}},
		{"Complete fails when not waiting pickup", func(t *testing.T) {
			// Create a new order in pending state
			newOrder := &domain.Order{
				ID:          "order2",
				MenuItemID:  "id1",
				MenuName:    "name1",
				Status:      domain.OrderStatusPending,
				OrderNumber: 2,
				CreatedAt:   time.Now(),
			}
			orderRepo.orders["order2"] = newOrder

			_, err := uc.Complete("store1", "order2")
			if err == nil {
				t.Error("expected error, got nil")
			}
		}},
		{"WaitingPickup fails when not pending", func(t *testing.T) {
			// Create a new order in completed state
			newOrder := &domain.Order{
				ID:          "order3",
				MenuItemID:  "id1",
				MenuName:    "name1",
				Status:      domain.OrderStatusCompleted,
				OrderNumber: 3,
				CreatedAt:   time.Now(),
			}
			orderRepo.orders["order3"] = newOrder

			_, err := uc.WaitingPickup("store1", "order3")
			if err == nil {
				t.Error("expected error, got nil")
			}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t)
		})
	}
}
