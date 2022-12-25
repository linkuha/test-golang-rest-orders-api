package entity

import validation "github.com/go-ozzo/ozzo-validation"

type Profile struct {
	UserID     string `json:"user_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	FullName   string `json:"full_name"`
	Sex        string
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
	return u.Sex == "m"
}

func (u *Profile) IsWoman() bool {
	return u.Sex == "w"
}
