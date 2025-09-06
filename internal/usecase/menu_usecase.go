package usecase

import (
	"kakigori-api/internal/domain"
	"kakigori-api/internal/repository"
)

type MenuUsecase struct {
	repo repository.MenuRepository
}

func NewMenuUsecase(repo repository.MenuRepository) *MenuUsecase {
	return &MenuUsecase{
		repo: repo,
	}
}

func (s *MenuUsecase) ListMenu() []domain.MenuItem {
	items, err := s.repo.ListMenu()
	if err != nil {
		return []domain.MenuItem{}
	}
	return items
}
