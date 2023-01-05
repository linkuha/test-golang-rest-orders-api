package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/service"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/usecase/user"
	"net/http"
)

// @Summary Add follower
// @Security ApiKeyAuth
// @Tags follower
// @Description add follower
// @ID follower-add
// @Accept  json
// @Produce  json
// @Param input body entity.Follower true "follower data"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /followers [post]
func (ctrl *Controller) addFollower(c *gin.Context) {
	var input entity.Follower

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, newJSONBindingErrorWrapper(err))
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	if userID != input.UserID {
		newErrorResponse(c, forbiddenError)
		return
	}

	encryptor := service.NewPasswordEncryptor()
	uc := user.NewUserUseCase(ctrl.repos.Users, encryptor)

	err = uc.AddFollower(ctrl.ctx, input)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{true})
}
