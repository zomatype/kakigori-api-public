package domain

type MenuItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewMenuItem(id, name, desc string) *MenuItem {
	return &MenuItem{
		ID:          id,
		Name:        name,
		Description: desc,
	}
}
