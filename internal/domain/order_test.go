package domain

import (
	"testing"
	"time"
)

func TestNewOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		id          string
		menuItemID  string
		menuName    string
		orderNumber int
	}{
		{"normal", "1", "2", "かき氷", 10},
		{"zero", "", "", "", 0},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewOrder(tt.id, tt.menuItemID, tt.menuName, tt.orderNumber)
			if got.ID != tt.id || got.MenuItemID != tt.menuItemID || got.MenuName != tt.menuName || got.OrderNumber != tt.orderNumber {
				t.Errorf("got %+v, want id=%s, menuItemID=%s, menuName=%s, orderNumber=%d", got, tt.id, tt.menuItemID, tt.menuName, tt.orderNumber)
			}
			if got.Status != OrderStatusPending {
				t.Errorf("got status %s, want %s", got.Status, OrderStatusPending)
			}
			if time.Since(got.CreatedAt) > time.Second {
				t.Errorf("CreatedAt too old: %v", got.CreatedAt)
			}
		})
	}
}

func TestOrder_Complete(t *testing.T) {
	t.Parallel()
	order := NewOrder("1", "2", "かき氷", 10)
	order.Complete()
	if order.Status != OrderStatusCompleted {
		t.Errorf("got %s, want %s", order.Status, OrderStatusCompleted)
	}
}
