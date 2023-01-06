package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/usecase/profile"
	"net/http"
)

// @Summary Create my profile
// @Security ApiKeyAuth
// @Tags profile
// @Description Create profile for logged user
// @ID profile-create-my
// @Accept  json
// @Produce  json
// @Param input body entity.Profile true "profile data"
// @Success 200 {string} string "id"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /profiles/my [post]
func (ctrl *Controller) createMyProfile(c *gin.Context) {
	var input entity.Profile

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, newJSONBindingErrorWrapper(err))
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, err)
		return
	}
	input.UserID = userId

	uc := profile.NewProfileUseCase(ctrl.repos.Profiles)
	if err := uc.Create(ctrl.ctx, input); err != nil {
		newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{true})
}

// @Summary Get profile
// @Security ApiKeyAuth
// @Tags profile
// @Description get profile
// @ID profile-get
// @Accept  json
// @Produce  json
// @Param id path string true "Profile ID"
// @Success 200 {object} entity.Profile
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /profiles/{id} [get]
func (ctrl *Controller) getProfile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newErrorResponse(c, emptyParameterID)
		return
	}

	p, err := ctrl.repos.Profiles.GetByUserID(ctrl.ctx, id)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, p)
}

// @Summary Get my profile
// @Security ApiKeyAuth
// @Tags profile
// @Description get profile of logged user
// @ID profile-get-my
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.Profile
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /profiles/my [get]
func (ctrl *Controller) getMyProfile(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	p, err := ctrl.repos.Profiles.GetByUserID(ctrl.ctx, userId)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, p)
}

// @Summary Update profile
// @Security ApiKeyAuth
// @Tags profile
// @Description update profile
// @ID profile-update
// @Accept  json
// @Produce  json
// @Param id path string true "Profile ID"
// @Param input body entity.Profile true "profile data"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /profiles/{id} [put]
func (ctrl *Controller) updateProfile(c *gin.Context) {
	var input entity.Profile

	id := c.Param("id")
	if id == "" {
		newErrorResponse(c, emptyParameterID)
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, newJSONBindingErrorWrapper(err))
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, err)
		return
	}
	input.UserID = userID

	if userID != id {
		newErrorResponse(c, forbiddenError)
		return
	}

	err = ctrl.repos.Profiles.Update(ctrl.ctx, &input)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{true})
}
