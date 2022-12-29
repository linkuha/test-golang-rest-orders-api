package service

import (
	"crypto/sha1"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

const (
	salt = "asdkjqw21e8h128hd12sa"
)

type PasswordEncryptor interface {
	EncryptString(s string) (string, error)
	CompareHashAndPassword(hash, password string) bool
}

func NewPasswordEncryptor() PasswordEncryptorBcrypt {
	return PasswordEncryptorBcrypt{}
}

type PasswordEncryptorBcrypt struct {
}

func (p PasswordEncryptorBcrypt) EncryptString(s string) (string, error) {
	return encryptWithBcrypt(s)
}

func encryptWithSha1(s string) (string, error) {
	hash := sha1.New()
	hash.Write([]byte(s))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt))), nil
}

func encryptWithBcrypt(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// CompareWithHash ...
func (p PasswordEncryptorBcrypt) CompareHashAndPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
