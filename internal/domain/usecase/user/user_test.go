package user_test

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
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

	repo := mockUser.NewMockRepository(ctrl)

	userName := "qwerty"
	mockResp := &entity.User{
		ID:           "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7",
		Username:     userName,
		PasswordHash: "testpassword",
	}
	expected := mockResp
	repo.EXPECT().GetByUsername(userName).Return(mockResp, nil).Times(1)

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	u, err := useCase.GetUserIfCredentialsValid(userName, "password")
	require.NoError(t, err)
	require.Equal(t, expected, u)
}

func TestGetError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockUser.NewMockRepository(ctrl)

	userName := "qwerty"
	mockResp := &entity.User{
		ID:           "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7",
		Username:     userName,
		PasswordHash: "testpassword",
	}
	repo.EXPECT().GetByUsername(userName).Return(mockResp, nil).Times(1)
	expectedErr := errors.New("invalid password")

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	u, err := useCase.GetUserIfCredentialsValid(userName, "wrong")
	require.Error(t, err)
	require.EqualError(t,
		expectedErr,
		err.Error(),
	)
	require.Nil(t, u)
}

func TestGetDbError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockUser.NewMockRepository(ctrl)
	repoErr := errors.New("db is down")

	userName := "qwerty"
	repo.EXPECT().GetByUsername(userName).Return(nil, repoErr).Times(1)
	expectedErr := fmt.Errorf("err from user repository: %w", repoErr)

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	u, err := useCase.GetUserIfCredentialsValid(userName, "password")
	require.Error(t, err)
	require.EqualError(t,
		expectedErr,
		err.Error(),
	)
	require.Nil(t, u)
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockUser.NewMockRepository(ctrl)

	u := entity.TestUser(t)
	u.PasswordHash = "testpassword"

	expected := "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"
	repo.EXPECT().Store(u).Return(expected, nil).Times(1)

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	id, err := useCase.Create(*u)
	require.NoError(t, err)
	require.Equal(t, expected, id)
}

func TestCreateValidateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockUser.NewMockRepository(ctrl)

	u := entity.User{
		Username: "f",
		Password: "password",
	}

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	id, err := useCase.Create(u)
	require.Error(t, err)
	require.ErrorContains(t, err, "validate error")
	require.Empty(t, id)
}

func TestCreateDbError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockUser.NewMockRepository(ctrl)
	dbErr := errors.New("db is down")

	u := entity.User{
		Username: "qwerty",
		Password: "password",
	}
	u2 := u
	u2.PasswordHash = "testpassword"
	repo.EXPECT().Store(&u2).Return("", dbErr).Times(1)
	expectedErr := fmt.Errorf("err from user repository: %w", dbErr)

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	id, err := useCase.Create(u)
	require.Error(t, err)
	require.EqualError(t, err, expectedErr.Error())
	require.Empty(t, id)
}

func TestRemove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockUser.NewMockRepository(ctrl)

	u := entity.TestUser(t)
	u.ID = "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"
	repo.EXPECT().Remove(u.ID).Return(nil).Times(1)

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	err := useCase.Remove(*u)
	require.NoError(t, err)
}

func TestRemoveDbError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockUser.NewMockRepository(ctrl)
	dbErr := errors.New("db is down")

	u := entity.TestUser(t)
	u.ID = "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"
	repo.EXPECT().Remove(u.ID).Return(dbErr).Times(1)
	expectedErr := fmt.Errorf("err from user repository: %w", dbErr)

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	err := useCase.Remove(*u)
	require.Error(t, err)
	require.EqualError(t, err, expectedErr.Error())
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockUser.NewMockRepository(ctrl)

	u := entity.TestUser(t)
	u.ID = "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"
	u.PasswordHash = "testpassword"
	repo.EXPECT().Update(u).Return(nil).Times(1)

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	err := useCase.Update(*u)
	require.NoError(t, err)
}

func TestUpdateValidateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockUser.NewMockRepository(ctrl)

	u := entity.TestUser(t)
	u.ID = "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	err := useCase.Update(*u)
	require.Error(t, err)
	require.ErrorContains(t, err, "validate error")

}

func TestUpdateDbError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockUser.NewMockRepository(ctrl)
	dbErr := errors.New("db is down")

	u := entity.TestUser(t)
	u.ID = "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"
	u.PasswordHash = "testpassword"

	repo.EXPECT().Update(u).Return(dbErr).Times(1)
	expectedErr := fmt.Errorf("err from user repository: %w", dbErr)

	encryptor := mockService.NewPasswordEncryptor()
	useCase := user.NewUserUseCase(repo, encryptor)
	err := useCase.Update(*u)
	require.Error(t, err)
	require.EqualError(t, err, expectedErr.Error())
}
