package entity

import validation "github.com/go-ozzo/ozzo-validation"

type Profile struct {
	UserID     string
	FirstName  string
	LastName   string
	MiddleName string
	FullName   string
	Sex        rune
	Age        int
}

// Validate ...
func (u *Profile) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.UserID, validation.Required),
		validation.Field(&u.FirstName, validation.Required),
		validation.Field(&u.LastName, validation.Required),
		validation.Field(&u.Sex, validation.Required),
	)
}

func (u *Profile) IsMan() bool {
	return u.Sex == 'm'
}

func (u *Profile) IsWoman() bool {
	return u.Sex == 'w'
}
