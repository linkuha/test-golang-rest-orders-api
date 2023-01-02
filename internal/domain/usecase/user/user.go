package user

import (
	"errors"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository/user"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/service"
)

type UseCase struct {
	repo      user.Repository
	encryptor service.PasswordEncryptor
}

func NewUserUseCase(repo user.Repository, encryptor service.PasswordEncryptor) *UseCase {
	return &UseCase{repo, encryptor}
}

func (uc *UseCase) GetUserIfCredentialsValid(username, password string) (*entity.User, error) {
	u, err := uc.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	if !u.ComparePassword(password, uc.encryptor) {
		return nil, errors.New("invalid password")
	}
	return u, nil
}

func (uc *UseCase) Create(user entity.User) (string, error) {
	if err := user.Validate(); err != nil {
		return "", err
	}

	if err := user.BeforeCreate(uc.encryptor); err != nil {
		return "", err
	}

	return uc.repo.Store(&user)
}

func (uc *UseCase) Remove(user entity.User) error {
	return uc.repo.Remove(user.ID)
}

func (uc *UseCase) Update(user entity.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	return uc.repo.Update(&user)
}
