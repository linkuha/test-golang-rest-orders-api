package user_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/errs"
	mockUser "github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository/user/mocks"
	mockService "github.com/linkuha/test-golang-rest-orders-api/internal/domain/service/mocks"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/usecase/user"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	repo := mockUser.NewMockRepository(ctrl)

	userName := "qwerty"
	mockResp := &entity.User{
		ID:           "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7",
		Username:     userName,
		PasswordHash: "testpassword",
	}
	expected := mockResp
	repo.EXPECT().GetByUsername(ctx, userName).Return(mockResp, nil).Times(1)

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	u, err := useCase.GetUserIfCredentialsValid(ctx, userName, "password")
	require.NoError(t, err)
	require.Equal(t, expected, u)
}

func TestGetError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	repo := mockUser.NewMockRepository(ctrl)

	userName := "qwerty"
	mockResp := &entity.User{
		ID:           "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7",
		Username:     userName,
		PasswordHash: "testpassword",
	}
	repo.EXPECT().GetByUsername(ctx, userName).Return(mockResp, nil).Times(1)
	expectedErr := errors.New("invalid password")

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	u, err := useCase.GetUserIfCredentialsValid(ctx, userName, "wrong")
	require.Error(t, err)
	require.EqualError(t,
		err,
		expectedErr.Error(),
	)
	require.Nil(t, u)
}

func TestGetDbError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	repo := mockUser.NewMockRepository(ctrl)
	repoErr := errors.New("db is down")

	userName := "qwerty"
	repo.EXPECT().GetByUsername(ctx, userName).Return(nil, repoErr).Times(1)

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	u, err := useCase.GetUserIfCredentialsValid(ctx, userName, "password")
	require.Error(t, err)
	require.IsType(t, errs.CustomErrorWrapper{}, err)
	require.EqualError(t,
		err,
		repoErr.Error(),
	)
	require.Nil(t, u)
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	repo := mockUser.NewMockRepository(ctrl)

	u := entity.TestUser(t)
	u.PasswordHash = "testpassword"

	expected := "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"
	repo.EXPECT().Store(ctx, u).Return(expected, nil).Times(1)

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	id, err := useCase.Create(ctx, *u)
	require.NoError(t, err)
	require.Equal(t, expected, id)
}

func TestCreateValidateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	repo := mockUser.NewMockRepository(ctrl)

	u := entity.User{
		Username: "f",
		Password: "password",
	}

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	id, err := useCase.Create(ctx, u)
	require.Error(t, err)
	require.IsType(t, errs.CustomErrorWrapper{}, err)
	var tmp errs.CustomErrorWrapper
	errors.As(err, &tmp)
	require.Equal(t, tmp.Code, errs.Validation)
	require.Empty(t, id)
}

func TestCreateDbError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	repo := mockUser.NewMockRepository(ctrl)
	dbErr := errors.New("db is down")

	u := entity.User{
		Username: "qwerty",
		Password: "password",
	}
	u2 := u
	u2.PasswordHash = "testpassword"
	repo.EXPECT().Store(ctx, &u2).Return("", dbErr).Times(1)

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	id, err := useCase.Create(ctx, u)
	require.Error(t, err)
	require.IsType(t, errs.CustomErrorWrapper{}, err)
	require.EqualError(t, err, dbErr.Error())
	require.Empty(t, id)
}

func TestRemove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	repo := mockUser.NewMockRepository(ctrl)

	u := entity.TestUser(t)
	u.ID = "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"
	repo.EXPECT().Remove(ctx, u.ID).Return(nil).Times(1)

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	err := useCase.Remove(ctx, *u)
	require.NoError(t, err)
}

func TestRemoveDbError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	repo := mockUser.NewMockRepository(ctrl)
	dbErr := errors.New("db is down")

	u := entity.TestUser(t)
	u.ID = "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"
	repo.EXPECT().Remove(ctx, u.ID).Return(dbErr).Times(1)

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	err := useCase.Remove(ctx, *u)
	require.Error(t, err)
	require.IsType(t, errs.CustomErrorWrapper{}, err)
	require.EqualError(t, err, dbErr.Error())
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	repo := mockUser.NewMockRepository(ctrl)

	u := entity.TestUser(t)
	u.ID = "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"
	u.PasswordHash = "testpassword"
	repo.EXPECT().Update(ctx, u).Return(nil).Times(1)

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	err := useCase.Update(ctx, *u)
	require.NoError(t, err)
}

func TestUpdateValidateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	repo := mockUser.NewMockRepository(ctrl)

	u := entity.TestUser(t)
	u.ID = "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	err := useCase.Update(ctx, *u)
	require.Error(t, err)
	require.IsType(t, errs.CustomErrorWrapper{}, err)
	var tmp errs.CustomErrorWrapper
	errors.As(err, &tmp)
	require.Equal(t, tmp.Code, errs.Validation)
}

func TestUpdateDbError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	repo := mockUser.NewMockRepository(ctrl)
	dbErr := errors.New("db is down")

	u := entity.TestUser(t)
	u.ID = "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"
	u.PasswordHash = "testpassword"

	repo.EXPECT().Update(ctx, u).Return(dbErr).Times(1)

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	err := useCase.Update(ctx, *u)
	require.Error(t, err)
	require.IsType(t, errs.CustomErrorWrapper{}, err)
	require.EqualError(t, err, dbErr.Error())
}

func TestAddFollower(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	repo := mockUser.NewMockRepository(ctrl)

	follower := entity.TestFollower(t)

	repo.EXPECT().AddFollower(ctx, follower.UserID, follower.FollowerID).Return(nil).Times(1)

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	err := useCase.AddFollower(ctx, *follower)
	require.NoError(t, err)
}

func TestAddFollowerValidateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	repo := mockUser.NewMockRepository(ctrl)

	follower := entity.Follower{
		UserID:     "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7",
		FollowerID: "",
	}

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	err := useCase.AddFollower(ctx, follower)
	require.Error(t, err)
	require.IsType(t, errs.CustomErrorWrapper{}, err)
	var tmp errs.CustomErrorWrapper
	errors.As(err, &tmp)
	require.Equal(t, tmp.Code, errs.Validation)
}

func TestAddFollowerSameIDError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	repo := mockUser.NewMockRepository(ctrl)

	follower := entity.Follower{
		UserID:     "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7",
		FollowerID: "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7",
	}

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	err := useCase.AddFollower(ctx, follower)
	require.Error(t, err)
	require.IsType(t, errs.CustomErrorWrapper{}, err)
	var tmp errs.CustomErrorWrapper
	errors.As(err, &tmp)
	require.Equal(t, tmp.Code, errs.Validation)
}

func TestAddFollowerDbError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	repo := mockUser.NewMockRepository(ctrl)
	dbErr := errors.New("db is down")

	follower := entity.TestFollower(t)

	repo.EXPECT().AddFollower(ctx, follower.UserID, follower.FollowerID).Return(dbErr).Times(1)

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	err := useCase.AddFollower(ctx, *follower)
	require.Error(t, err)
	require.IsType(t, errs.CustomErrorWrapper{}, err)
	require.EqualError(t, err, dbErr.Error())
}
