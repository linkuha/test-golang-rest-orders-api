package profile

import (
	"context"
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

func (uc *UseCase) Create(ctx context.Context, profile entity.Profile) error {
	if err := profile.Validate(); err != nil {
		return errs.NewErrorWrapper(errs.Validation, err, "profile validation error")
	}

	if _, err := uc.repo.GetByUserID(ctx, profile.UserID); err == nil {
		return errs.NewErrorWrapper(errs.Exist, nil, "already exist")
	}

	_, err := uc.repo.Store(ctx, &profile)
	if err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from profile repo")
	}
	return nil
}

func (uc *UseCase) Remove(ctx context.Context, profile entity.Profile) error {
	if err := uc.repo.RemoveByUserID(ctx, profile.UserID); err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from profile repo")
	}
	return nil
}

func (uc *UseCase) Update(ctx context.Context, profile entity.Profile) error {
	if err := profile.Validate(); err != nil {
		return errs.NewErrorWrapper(errs.Validation, err, "profile validation error")
	}

	if err := uc.repo.Update(ctx, &profile); err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from profile repo")
	}
	return nil
}
