package user

import (
	"context"
	"errors"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/errs"
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

func (uc *UseCase) GetUserIfCredentialsValid(ctx context.Context, username, password string) (*entity.User, error) {
	u, err := uc.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errs.NewErrorWrapper(errs.Database, err, "error from user repo")
	}

	if !u.ComparePassword(password, uc.encryptor) {
		return nil, errs.NewErrorWrapper(errs.UserCredentials, errs.InvalidPassword, "invalid credentials")
	}
	return u, nil
}

func (uc *UseCase) Create(ctx context.Context, user entity.User) (string, error) {
	if err := user.Validate(); err != nil {
		return "", errs.NewErrorWrapper(errs.Validation, err, "user validation error")
	}

	exist, _ := uc.repo.GetByUsername(ctx, user.Username)
	if exist != nil {
		return "", errs.NewErrorWrapper(errs.Exist, nil, "username is already taken")
	}

	if err := user.BeforeCreate(uc.encryptor); err != nil {
		return "", errs.NewErrorWrapper(errs.Internal, err, "encryptor error")
	}

	res, err := uc.repo.Store(ctx, &user)
	if err != nil {
		return "", errs.NewErrorWrapper(errs.Database, err, "error from user repo")
	}
	return res, nil
}

func (uc *UseCase) Remove(ctx context.Context, user entity.User) error {
	if err := uc.repo.Remove(ctx, user.ID); err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from user repo")
	}
	return nil
}

func (uc *UseCase) Update(ctx context.Context, user entity.User) error {
	if err := user.Validate(); err != nil {
		return errs.NewErrorWrapper(errs.Validation, err, "user validation error")
	}
	if user.PasswordHash == "" {
		return errs.NewErrorWrapper(errs.Validation, errors.New("empty password hash"), "user validation error")
	}

	if err := uc.repo.Update(ctx, &user); err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from user repo")
	}
	return nil
}

func (uc *UseCase) AddFollower(ctx context.Context, follower entity.Follower) error {
	if err := follower.Validate(); err != nil {
		return errs.NewErrorWrapper(errs.Validation, err, "follower validation error")
	}

	if err := uc.repo.AddFollower(ctx, follower.UserID, follower.FollowerID); err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from user repo")
	}
	return nil
}
