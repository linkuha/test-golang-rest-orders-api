package user_test

import (
	"github.com/golang/mock/gomock"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	mockUser "github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository/user/mocks"
	mock_service "github.com/linkuha/test-golang-rest-orders-api/internal/domain/service/mocks"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/usecase/user"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockUser.NewMockRepository(ctrl)

	encryptor := mock_service.NewPasswordEncryptor()

	userName := "qwerty"
	mockResp := &entity.User{
		ID:           "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7",
		Username:     userName,
		PasswordHash: "testpassword",
	}
	expected := mockResp
	repo.EXPECT().GetByUsername(userName).Return(mockResp, nil).Times(1)

	useCase := user.NewUserUseCase(repo, encryptor)
	u, err := useCase.GetUserIfCredentialsValid(userName, "password")
	require.NoError(t, err)
	require.Equal(t, expected, u)
}
