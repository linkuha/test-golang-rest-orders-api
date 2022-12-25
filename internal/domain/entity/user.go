package entity

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type User struct {
	ID           string
	Username     string
	Password     string
	PasswordHash string `json:"password_hash"`
	//Status       int
	//Roles        string
}

type PasswordEncryptor interface {
	EncryptString(s string) (string, error)
	CompareHashAndPassword(hash, password string) bool
}

// Validate ...
func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Username, validation.Required),
		validation.Field(&u.Password, validation.By(requiredIf(u.PasswordHash == "")), validation.Length(6, 100)),
	)
}

// BeforeCreate ...
func (u *User) BeforeCreate(encryptor PasswordEncryptor) error {
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
func (u *User) ComparePassword(password string, encryptor PasswordEncryptor) bool {
	return encryptor.CompareHashAndPassword(u.PasswordHash, password)
}

// Sanitize ...
func (u *User) Sanitize() {
	u.Password = ""
}
