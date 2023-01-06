package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Profile struct {
	UserID     string `json:"user_id"`
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	MiddleName string `json:"middle_name"`
	FullName   string `json:"full_name"`
	Sex        string `json:"sex" binding:"required"`
	Age        int    `json:"age" binding:"required"`
}

// Validate ...
func (u *Profile) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.UserID, validation.Required, is.UUIDv4),
		validation.Field(&u.FirstName, validation.Required),
		validation.Field(&u.LastName, validation.Required),
		validation.Field(&u.Sex, validation.Required, validation.Length(1, 1)),
		validation.Field(&u.Age, validation.Required),
	)
}

func (u *Profile) IsMan() bool {
	return u.Sex == "m"
}

func (u *Profile) IsWoman() bool {
	return u.Sex == "w"
}
