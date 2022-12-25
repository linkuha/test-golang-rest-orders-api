package entity

import validation "github.com/go-ozzo/ozzo-validation"

type Product struct {
	ID          string
	Name        string
	Description string
	LeftInStock int `json:"left_in_stock"`
}

type Price struct {
	Currency string
	Price    string
	//Status   string
}

// Validate ...
func (m *Product) Validate() error {
	return validation.ValidateStruct(
		m,
		validation.Field(&m.Name, validation.Required),
	)
}

func (m *Price) Validate() error {
	return validation.ValidateStruct(
		m,
		validation.Field(&m.Currency, validation.Required, validation.Length(3, 3)),
		validation.Field(&m.Price, validation.Required, validation.Min(0)),
	)
}

type ProductUpdateInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	LeftInStock *bool   `json:"left_in_stock"`
}
