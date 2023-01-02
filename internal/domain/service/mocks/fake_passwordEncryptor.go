package mock_service

const (
	salt = "test"
)

func NewPasswordEncryptor() PasswordEncryptorFake {
	return PasswordEncryptorFake{}
}

type PasswordEncryptorFake struct {
}

func (p PasswordEncryptorFake) EncryptString(s string) (string, error) {
	return salt + s, nil
}

func (p PasswordEncryptorFake) CompareHashAndPassword(hash, password string) bool {
	return hash == salt+password
}
