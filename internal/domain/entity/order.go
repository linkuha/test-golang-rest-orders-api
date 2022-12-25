package entity

import validation "github.com/go-ozzo/ozzo-validation"

type Order struct {
	ID     string
	UserID string
	Number int
}

type OrderProduct struct {
	ID        int
	OrderID   string
	ProductID string
	Amount    int
}

type OrderProductView struct {
	ID     string
	Amount int
}

// Validate ...
func (m *Order) Validate() error {
	return validation.ValidateStruct(
		m,
		validation.Field(&m.UserID, validation.Required),
		validation.Field(&m.Number, validation.Required),
	)
}

func (m *OrderProduct) Validate() error {
	return validation.ValidateStruct(
		m,
		validation.Field(&m.OrderID, validation.Required),
		validation.Field(&m.ProductID, validation.Required),
		validation.Field(&m.Amount, validation.Min(1)),
	)
}
