package repository

import (
	"kakigori-api/internal/domain"
	"testing"
)

func TestMemoryStoreRepository(t *testing.T) {
	t.Parallel()
	r := NewMemoryStoreRepository()
	store := &domain.Store{ID: "store1", MaxOrders: 10}
	r.Create(store)

	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{"Create and GetByID", func(t *testing.T) {
			got, err := r.GetByID("store1")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ID != "store1" {
				t.Errorf("expected store1, got %s", got.ID)
			}
		}},
		{"Create duplicate store returns error", func(t *testing.T) {
			err := r.Create(store)
			if err == nil {
				t.Errorf("expected error for duplicate store, got nil")
			}
		}},
		{"GetByID not found returns error", func(t *testing.T) {
			_, err := r.GetByID("notfound")
			if err == nil {
				t.Errorf("expected error for notfound, got nil")
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
