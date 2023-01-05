package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	LeftInStock int     `json:"left_in_stock" binding:"required"`
	Prices      []Price `json:"prices"`
}

type Price struct {
	Currency string `json:"currency" binding:"required"`
	Price    string `json:"price" binding:"required"` // TODO float64?
	//Status   string
}

// Validate ...
func (m *Product) Validate() error {
	if len(m.Prices) > 0 {
		for _, price := range m.Prices {
			if err := price.Validate(); err != nil {
				return err
			}
		}
	}
	return validation.ValidateStruct(
		m,
		validation.Field(&m.ID, is.UUIDv4),
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.LeftInStock, validation.Min(0)),
	)
}

func (m *Price) Validate() error {
	return validation.ValidateStruct(
		m,
		validation.Field(&m.Currency, validation.Required, validation.Length(3, 3)),
		validation.Field(&m.Price, validation.Required),
	)
}

type ProductUpdateInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	LeftInStock *bool   `json:"left_in_stock"`
}
