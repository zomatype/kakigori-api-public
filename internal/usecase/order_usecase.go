package usecase

import (
	"fmt"
	"kakigori-api/internal/domain"
	"kakigori-api/internal/repository"
)

type OrderUsecase struct {
	menuRepo  repository.MenuRepository
	orderRepo repository.OrderRepository
	storeRepo repository.StoreRepository
}

func NewOrderUsecase(
	menuRepo repository.MenuRepository,
	orderRepo repository.OrderRepository,
	storeRepo repository.StoreRepository,
) *OrderUsecase {
	return &OrderUsecase{
		menuRepo:  menuRepo,
		orderRepo: orderRepo,
		storeRepo: storeRepo,
	}
}

func (ou *OrderUsecase) CreateOrder(storeId, menuItemId string) (*domain.Order, error) {
	store, err := ou.storeRepo.GetByID(storeId)
	if err != nil {
		return nil, err
	}

	menuItem, err := ou.menuRepo.GetMenu(menuItemId)
	if err != nil {
		return nil, err
	}

	count, err := ou.orderRepo.GetOrderCount(store.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order count: %v", err)
	}

	if count >= store.MaxOrders {
		oldestOrder, err := ou.orderRepo.GetOldestOrder(store.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get oldest order: %v", err)
		}
		if oldestOrder != nil {
			if err := ou.orderRepo.Delete(store.ID, oldestOrder.ID); err != nil {
				return nil, fmt.Errorf("failed to delete oldest order: %v", err)
			}
		}
	}

	orderNumber, err := ou.orderRepo.GetNextOrderNumber(store.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order number: %v", err)
	}

	newOrder := domain.NewOrder(
		fmt.Sprintf("%s-%d", store.ID, orderNumber),
		menuItemId,
		menuItem.Name,
		orderNumber,
	)

	if err := ou.orderRepo.Create(store.ID, newOrder); err != nil {
		return nil, fmt.Errorf("failed to create order: %v", err)
	}

	return newOrder, nil
}

func (ou *OrderUsecase) ListOrders(storeId string) ([]*domain.Order, error) {
	orders, err := ou.orderRepo.GetAllByStore(storeId)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (ou *OrderUsecase) GetOrder(storeId, orderID string) (*domain.Order, error) {
	return ou.orderRepo.GetByID(storeId, orderID)
}

func (ou *OrderUsecase) Complete(storeId, orderID string) (*domain.Order, error) {
	order, err := ou.orderRepo.GetByID(storeId, orderID)
	if err != nil {
		return nil, err
	}

	if order.Status != domain.OrderStatusWaitingPickup {
		return nil, fmt.Errorf("order status is not waiting pickup")
	}

	order.Complete()
	if err := ou.orderRepo.Update(storeId, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (ou *OrderUsecase) WaitingPickup(storeId, orderID string) (*domain.Order, error) {
	order, err := ou.orderRepo.GetByID(storeId, orderID)
	if err != nil {
		return nil, err
	}

	if order.Status != domain.OrderStatusPending {
		return nil, fmt.Errorf("order status is not pending")
	}

	order.WaitingPickup()
	if err := ou.orderRepo.Update(storeId, order); err != nil {
		return nil, err
	}

	return order, nil
}
