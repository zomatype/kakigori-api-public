package usecase

import (
	"kakigori-api/internal/domain"
	"testing"
)

type dummyMenuRepo struct{}

func (d *dummyMenuRepo) ListMenu() ([]domain.MenuItem, error) {
	return []domain.MenuItem{{ID: "id1", Name: "name1", Description: "desc1"}}, nil
}
func (d *dummyMenuRepo) GetMenu(id string) (*domain.MenuItem, error) {
	if id == "id1" {
		return &domain.MenuItem{ID: "id1", Name: "name1", Description: "desc1"}, nil
	}
	return nil, nil
}

func TestMenuUsecase(t *testing.T) {
	t.Parallel()
	uc := NewMenuUsecase(&dummyMenuRepo{})
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{"ListMenu returns items", func(t *testing.T) {
			items := uc.ListMenu()
			if len(items) != 1 {
				t.Errorf("expected 1 item, got %d", len(items))
			}
			if items[0].ID != "id1" {
				t.Errorf("expected id1, got %s", items[0].ID)
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
