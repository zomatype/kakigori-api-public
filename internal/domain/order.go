package domain

import "time"

// オペレーション次第であとで使う可能性がある。一旦残してる。
const (
	OrderStatusPending       = "pending"
	OrderStatusWaitingPickup = "waitingPickup"
	OrderStatusCompleted     = "completed"
)

type Order struct {
	ID          string    `json:"id"`
	MenuItemID  string    `json:"menu_item_id"`
	MenuName    string    `json:"menu_name"`
	Status      string    `json:"status"`
	OrderNumber int       `json:"order_number"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewOrder(id, menuItemID, menuName string, orderNumber int) *Order {
	return &Order{
		ID:          id,
		MenuItemID:  menuItemID,
		MenuName:    menuName,
		Status:      OrderStatusPending,
		OrderNumber: orderNumber,
		CreatedAt:   time.Now(),
	}
}

// 更新すべきかチェックすべき
func (o *Order) Complete() {
	o.Status = OrderStatusCompleted
}

// 更新すべきかチェックすべき
func (o *Order) WaitingPickup() {
	o.Status = OrderStatusWaitingPickup
}
