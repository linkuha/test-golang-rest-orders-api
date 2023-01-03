package user

import (
	"errors"
	"fmt"
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
		return nil, fmt.Errorf("err from user repository: %w", err)
	}

	if !u.ComparePassword(password, uc.encryptor) {
		return nil, errors.New("invalid password")
	}
	return u, nil
}

func (uc *UseCase) Create(user entity.User) (string, error) {
	if err := user.Validate(); err != nil {
		return "", fmt.Errorf("validate error: %w", err)
	}

	if err := user.BeforeCreate(uc.encryptor); err != nil {
		return "", fmt.Errorf("encryptor error: %w", err)
	}

	res, err := uc.repo.Store(&user)
	if err != nil {
		return "", fmt.Errorf("err from user repository: %w", err)
	}
	return res, nil
}

func (uc *UseCase) Remove(user entity.User) error {
	if err := uc.repo.Remove(user.ID); err != nil {
		return fmt.Errorf("err from user repository: %w", err)
	}
	return nil
}

func (uc *UseCase) Update(user entity.User) error {
	if err := user.Validate(); err != nil {
		return fmt.Errorf("validate error: %w", err)
	}
	if user.PasswordHash == "" {
		return fmt.Errorf("validate error: %w", errors.New("empty password hash"))
	}

	if err := uc.repo.Update(&user); err != nil {
		return fmt.Errorf("err from user repository: %w", err)
	}
	return nil
}
