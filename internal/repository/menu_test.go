package repository

import (
	"testing"
)

func TestMemoryMenuRepository(t *testing.T) {
	t.Parallel()
	r := NewMemoryMenuRepository()
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{"ListMenu returns all items", func(t *testing.T) {
			items, err := r.ListMenu()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(items) != 4 {
				t.Errorf("expected 4 items, got %d", len(items))
			}
		}},
		{"GetMenu returns item by ID", func(t *testing.T) {
			item, err := r.GetMenu("giiku-sai")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if item.ID != "giiku-sai" {
				t.Errorf("expected ID giiku-sai, got %s", item.ID)
			}
		}},
		{"GetMenu returns error for notfound", func(t *testing.T) {
			_, err := r.GetMenu("notfound")
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
