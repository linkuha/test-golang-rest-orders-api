package profile

import (
	"fmt"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository/profile"
)

type UseCase struct {
	repo profile.Repository
}

func NewProfileUseCase(repo profile.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Create(profile entity.Profile) (int, error) {
	if err := profile.Validate(); err != nil {
		return 0, fmt.Errorf("validate error: %w", err)
	}

	res, err := uc.repo.Store(&profile)
	if err != nil {
		return 0, fmt.Errorf("err from profile repository: %w", err)
	}
	return res, nil
}

func (uc *UseCase) Remove(profile entity.Profile) error {
	if err := uc.repo.RemoveByUserID(profile.UserID); err != nil {
		return fmt.Errorf("err from profile repository: %w", err)
	}
	return nil
}

func (uc *UseCase) Update(profile entity.Profile) error {
	if err := profile.Validate(); err != nil {
		return fmt.Errorf("validate error: %w", err)
	}

	if err := uc.repo.Update(&profile); err != nil {
		return fmt.Errorf("err from profile repository: %w", err)
	}
	return nil
}
