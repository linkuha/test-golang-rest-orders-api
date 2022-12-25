package profile

import (
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
		return 0, err
	}

	return uc.repo.Store(&profile)
}

func (uc *UseCase) Remove(profile entity.Profile) error {
	return uc.repo.RemoveByUserID(profile.UserID)
}

func (uc *UseCase) Update(profile entity.Profile) error {
	if err := profile.Validate(); err != nil {
		return err
	}

	return uc.repo.Update(&profile)
}
