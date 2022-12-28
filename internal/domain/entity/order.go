package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Order struct {
	ID     string
	UserID string `json:"user_id"`
	Number int
}

type OrderProduct struct {
	ID        int
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
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
		validation.Field(&m.ID, is.UUIDv4),
		validation.Field(&m.UserID, validation.Required, is.UUIDv4),
		validation.Field(&m.Number, validation.Required),
	)
}

func (m *OrderProduct) Validate() error {
	return validation.ValidateStruct(
		m,
		validation.Field(&m.OrderID, validation.Required, is.UUIDv4),
		validation.Field(&m.ProductID, validation.Required, is.UUIDv4),
		validation.Field(&m.Amount, validation.Min(1)),
	)
}
