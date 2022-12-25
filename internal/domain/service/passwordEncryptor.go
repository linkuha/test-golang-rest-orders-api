package service

import (
	"crypto/sha1"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

const (
	salt = "asdkjqw21e8h128hd12sa"
)

type PasswordEncryptor struct {
}

func (p PasswordEncryptor) EncryptString(s string) (string, error) {
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
func (p PasswordEncryptor) CompareHashAndPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
