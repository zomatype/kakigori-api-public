package domain

import (
	"reflect"
	"testing"
)

func TestNewMenuItem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		id       string
		menuName string
		desc     string
		want     *MenuItem
	}{
		{"normal", "1", "かき氷", "いちご", &MenuItem{ID: "1", Name: "かき氷", Description: "いちご"}},
		{"empty", "", "", "", &MenuItem{ID: "", Name: "", Description: ""}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewMenuItem(tt.id, tt.menuName, tt.desc)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}
