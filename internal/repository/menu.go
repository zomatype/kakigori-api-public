package repository

import (
	"fmt"
	"kakigori-api/internal/domain"
)

type MenuRepository interface {
	ListMenu() ([]domain.MenuItem, error)
	GetMenu(id string) (*domain.MenuItem, error)
}

type MemoryMenuRepository struct {
	items []domain.MenuItem
}

func NewMemoryMenuRepository() *MemoryMenuRepository {
	return &MemoryMenuRepository{
		items: []domain.MenuItem{
			{ID: "giiku-sai", Name: "技育祭な いちご味", Description: "技育祭をイメージしたいちご味のかき氷"},
			{ID: "giiku-haku", Name: "技育博な メロン味", Description: "技育博をイメージしたメロン味のかき氷"},
			{ID: "giiku-ten", Name: "技育展な ブルーハワイ味", Description: "技育展をイメージしたブルーハワイ味のかき氷"},
			{ID: "giiku-camp", Name: "技育CAMPな オレンジ味", Description: "技育CAMPをイメージしたオレンジ味のかき氷"},
		},
	}
}

func (r *MemoryMenuRepository) ListMenu() ([]domain.MenuItem, error) {
	return r.items, nil
}

func (r *MemoryMenuRepository) GetMenu(id string) (*domain.MenuItem, error) {
	for _, item := range r.items {
		if item.ID == id {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("menu item not found: %s", id)
}
