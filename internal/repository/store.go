package repository

import (
	"fmt"
	"kakigori-api/internal/domain"
	"sync"
)

var (
	ErrStoreNotFound      = fmt.Errorf("store not found")
	ErrStoreAlreadyExists = fmt.Errorf("store already exists")
)

type StoreRepository interface {
	Create(s *domain.Store) error
	GetByID(id string) (*domain.Store, error)
}

type MemoryStoreRepository struct {
	stores map[string]*domain.Store
	mu     sync.RWMutex
}

func NewMemoryStoreRepository() *MemoryStoreRepository {
	return &MemoryStoreRepository{
		stores: make(map[string]*domain.Store),
	}
}

func (r *MemoryStoreRepository) Create(s *domain.Store) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.stores[s.ID]; exists {
		return ErrStoreAlreadyExists
	}
	r.stores[s.ID] = s
	return nil
}

func (r *MemoryStoreRepository) GetByID(id string) (*domain.Store, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, exists := r.stores[id]
	if !exists {
		return nil, ErrStoreNotFound
	}
	return s, nil
}
