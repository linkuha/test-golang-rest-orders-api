package service_mocks

const (
	salt = "asdkjqw21e8h128hd12sa"
)

func NewPasswordEncryptor() PasswordEncryptorStub {
	return PasswordEncryptorStub{}
}

type PasswordEncryptorStub struct {
}

func (p PasswordEncryptorStub) EncryptString(s string) (string, error) {
	return salt + s, nil
}

func (p PasswordEncryptorStub) CompareHashAndPassword(hash, password string) bool {
	return hash == salt+password
}
