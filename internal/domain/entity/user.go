package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/service"
	"regexp"
)

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	PasswordHash string `json:"password_hash"`
	//Status       int
	//Roles        string
}

// Validate ...
func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.ID, is.UUIDv4),
		validation.Field(&u.Username, validation.Required, validation.Match(regexp.MustCompile("^[a-zA-Z0-9_-]{3,255}$"))),
		validation.Field(&u.Password, validation.Required.When(u.PasswordHash == ""), validation.Length(6, 100)),
		validation.Field(&u.PasswordHash, validation.Required.When(u.Password == "")),
	)
}

// BeforeCreate ...
func (u *User) BeforeCreate(encryptor service.PasswordEncryptor) error {
	if len(u.Password) > 0 {
		enc, err := encryptor.EncryptString(u.Password)
		if err != nil {
			return err
		}

		u.PasswordHash = enc
	}

	return nil
}

// ComparePassword ...
func (u *User) ComparePassword(password string, encryptor service.PasswordEncryptor) bool {
	return encryptor.CompareHashAndPassword(u.PasswordHash, password)
}

// Sanitize ...
func (u *User) Sanitize() {
	u.Password = ""
}
