package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Order struct {
	ID     string `json:"id"`
	UserID string `json:"user_id" binding:"required"`
	Number int    `binding:"required"`
}

type OrderProduct struct {
	ID        int    `json:"id"`
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Amount    int    `json:"amount"`
}

type OrderProductView struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"`
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
