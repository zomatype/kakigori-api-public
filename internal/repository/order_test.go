package repository

import (
	"kakigori-api/internal/domain"
	"testing"
	"time"
)

func TestMemoryOrderRepository(t *testing.T) {
	t.Parallel()
	storeID := "store1"
	order := &domain.Order{ID: "order1", MenuItemID: "m1", MenuName: "test", OrderNumber: 1, CreatedAt: time.Now()}

	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{"Create and GetByID", func(t *testing.T) {
			r := NewMemoryOrderRepository()
			r.Create(storeID, order)
			got, err := r.GetByID(storeID, "order1")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ID != "order1" {
				t.Errorf("expected order1, got %s", got.ID)
			}
		}},
		{"GetAllByStore", func(t *testing.T) {
			r := NewMemoryOrderRepository()
			r.Create(storeID, order)
			orders, err := r.GetAllByStore(storeID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(orders) != 1 {
				t.Errorf("expected 1 order, got %d", len(orders))
			}
		}},
		{"Delete order", func(t *testing.T) {
			r := NewMemoryOrderRepository()
			r.Create(storeID, order)
			err := r.Delete(storeID, "order1")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			_, err = r.GetByID(storeID, "order1")
			if err == nil {
				t.Errorf("expected error for deleted order, got nil")
			}
		}},
		{"Update order", func(t *testing.T) {
			r := NewMemoryOrderRepository()
			r.Create(storeID, order)
			order2 := &domain.Order{ID: "order2", MenuItemID: "m2", MenuName: "test2", OrderNumber: 2, CreatedAt: time.Now()}
			r.Create(storeID, order2)
			order2.MenuName = "updated"
			err := r.Update(storeID, order2)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got, _ := r.GetByID(storeID, "order2")
			if got.MenuName != "updated" {
				t.Errorf("expected updated, got %s", got.MenuName)
			}
		}},
		{"GetOldestOrder", func(t *testing.T) {
			r := NewMemoryOrderRepository()
			r.Create(storeID, order)
			order3 := &domain.Order{ID: "order3", MenuItemID: "m3", MenuName: "test3", OrderNumber: 3, CreatedAt: time.Now().Add(-time.Hour)}
			r.Create(storeID, order3)
			oldest, err := r.GetOldestOrder(storeID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if oldest.ID != "order3" {
				t.Errorf("expected order3, got %s", oldest.ID)
			}
		}},
		{"GetNextOrderNumber", func(t *testing.T) {
			r := NewMemoryOrderRepository()
			r.Create(storeID, order)
			num, err := r.GetNextOrderNumber(storeID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if num < 1 {
				t.Errorf("expected >=1, got %d", num)
			}
		}},
		{"GetOrderCount", func(t *testing.T) {
			r := NewMemoryOrderRepository()
			r.Create(storeID, order)
			count, err := r.GetOrderCount(storeID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if count < 1 {
				t.Errorf("expected >=1, got %d", count)
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
