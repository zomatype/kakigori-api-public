package domain

import (
	"reflect"
	"testing"
)

func TestNewStore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        string
		maxOrders int
		want      *Store
	}{
		{"normal", "1", 10, &Store{ID: "1", MaxOrders: 10}},
		{"zero", "", 0, &Store{ID: "", MaxOrders: 0}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewStore(tt.id, tt.maxOrders)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}
