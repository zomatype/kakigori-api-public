package repository

import (
	"fmt"
	"kakigori-api/internal/domain"
	"sync"
)

type OrderRepository interface {
	Create(storeID string, order *domain.Order) error
	GetByID(storeID, orderID string) (*domain.Order, error)
	GetAllByStore(storeID string) ([]*domain.Order, error)
	Update(storeID string, order *domain.Order) error
	Delete(storeID, orderID string) error
	GetNextOrderNumber(storeID string) (int, error)
	GetOrderCount(storeID string) (int, error)
	GetOldestOrder(storeID string) (*domain.Order, error)
}

func (r *MemoryOrderRepository) GetAllByStore(storeID string) ([]*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	storeOrders, exists := r.orders[storeID]
	if !exists {
		return nil, fmt.Errorf("store not found: %s", storeID)
	}
	orders := make([]*domain.Order, 0, len(storeOrders))
	for _, order := range storeOrders {
		orders = append(orders, order)
	}
	return orders, nil
}

type MemoryOrderRepository struct {
	mu           sync.RWMutex
	orders       map[string]map[string]*domain.Order
	orderNumbers map[string]int
}

func NewMemoryOrderRepository() *MemoryOrderRepository {
	return &MemoryOrderRepository{
		orders:       make(map[string]map[string]*domain.Order),
		orderNumbers: make(map[string]int),
	}
}

func (r *MemoryOrderRepository) InitializeStore(store *domain.Store) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[store.ID] = make(map[string]*domain.Order)
}

func (r *MemoryOrderRepository) Create(storeID string, order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[storeID]; !exists {
		r.orders[storeID] = make(map[string]*domain.Order)
	}
	r.orders[storeID][order.ID] = order
	return nil
}

func (r *MemoryOrderRepository) GetByID(storeID, orderID string) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	storeOrders, exists := r.orders[storeID]
	if !exists {
		return nil, fmt.Errorf("store not found: %s", storeID)
	}
	order, ok := storeOrders[orderID]
	if !ok {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}

	return order, nil
}

func (r *MemoryOrderRepository) GetOldestOrder(storeID string) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	storeOrders, exists := r.orders[storeID]
	if !exists || len(storeOrders) == 0 {
		return nil, nil
	}
	var oldest *domain.Order
	for _, order := range storeOrders {
		if oldest == nil || order.CreatedAt.Before(oldest.CreatedAt) {
			oldest = order
		}
	}
	return oldest, nil
}

func (r *MemoryOrderRepository) GetNextOrderNumber(storeID string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orderNumbers[storeID]++
	return r.orderNumbers[storeID], nil
}

func (r *MemoryOrderRepository) GetOrderCount(storeID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	storeOrders, exists := r.orders[storeID]
	if !exists {
		return 0, nil
	}
	return len(storeOrders), nil
}

func (r *MemoryOrderRepository) Update(storeID string, order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	storeOrders, exists := r.orders[storeID]
	if !exists {
		return fmt.Errorf("store not found: %s", storeID)
	}
	if _, ok := storeOrders[order.ID]; !ok {
		return fmt.Errorf("order not found: %s", order.ID)
	}
	storeOrders[order.ID] = order
	return nil
}

func (r *MemoryOrderRepository) Delete(storeID, orderID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	storeOrders, exists := r.orders[storeID]
	if !exists {
		return fmt.Errorf("store not found: %s", storeID)
	}
	if _, ok := storeOrders[orderID]; !ok {
		return fmt.Errorf("order not found: %s", orderID)
	}
	delete(storeOrders, orderID)
	return nil
}
