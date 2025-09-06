package domain

type Store struct {
	ID        string `json:"id"`
	MaxOrders int
}

func NewStore(id string, maxOrders int) *Store {
	return &Store{
		ID:        id,
		MaxOrders: maxOrders,
	}
}
