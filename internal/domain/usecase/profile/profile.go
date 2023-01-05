package profile

import (
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/errs"
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
		return 0, errs.NewErrorWrapper(errs.Validation, err, "profile validation error")
	}

	res, err := uc.repo.Store(&profile)
	if err != nil {
		return 0, errs.NewErrorWrapper(errs.Database, err, "error from profile repo")
	}
	return res, nil
}

func (uc *UseCase) Remove(profile entity.Profile) error {
	if err := uc.repo.RemoveByUserID(profile.UserID); err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from profile repo")
	}
	return nil
}

func (uc *UseCase) Update(profile entity.Profile) error {
	if err := profile.Validate(); err != nil {
		return errs.NewErrorWrapper(errs.Validation, err, "profile validation error")
	}

	if err := uc.repo.Update(&profile); err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from profile repo")
	}
	return nil
}
