package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/service"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/usecase/user"
	"net/http"
)

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body signInInput true "account info"
// @Success 200 {integer} string "id"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (ctrl *Controller) signUp(ctx *gin.Context) {
	var input signInInput

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, newJSONBindingErrorWrapper(err))
		return
	}

	encryptor := service.NewPasswordEncryptor()
	uc := user.NewUserUseCase(ctrl.repos.Users, encryptor)
	u := entity.User{Username: input.Username, Password: input.Password}
	id, err := uc.Create(u)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body signInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (ctrl *Controller) signIn(ctx *gin.Context) {
	var input signInInput

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, newJSONBindingErrorWrapper(err))
		return
	}

	encryptor := service.NewPasswordEncryptor()
	uc := user.NewUserUseCase(ctrl.repos.Users, encryptor)
	u, err := uc.GetUserIfCredentialsValid(input.Username, input.Password)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	authTokenGenerator := service.AuthTokenGenerator{}
	token, err := authTokenGenerator.GenerateToken(u.ID)
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
