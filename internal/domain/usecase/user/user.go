package user

import (
	"errors"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository/user"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/service"
)

type UseCase struct {
	repo user.Repository
}

func NewUserUseCase(repo user.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetUserIfCredentialsValid(username, password string) (*entity.User, error) {
	u, err := uc.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	encryptor := service.PasswordEncryptor{}
	if !u.ComparePassword(password, &encryptor) {
		return nil, errors.New("invalid password")
	}
	return u, nil
}

func (uc *UseCase) Create(user entity.User) (int, error) {
	if err := user.Validate(); err != nil {
		return 0, err
	}

	encryptor := service.PasswordEncryptor{}
	if err := user.BeforeCreate(&encryptor); err != nil {
		return 0, err
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
