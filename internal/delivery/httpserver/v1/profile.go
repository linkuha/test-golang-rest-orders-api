package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/usecase/profile"
	"net/http"
)

// @Summary Create profile
// @Security ApiKeyAuth
// @Tags profile
// @Description Create profile for user
// @ID profile-create
// @Accept  json
// @Produce  json
// @Param input body entity.Profile true "profile data"
// @Success 200 {string} string "id"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /profile [post]
func (ctrl *Controller) createProfile(c *gin.Context) {
	var input entity.Profile

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	input.UserID = userId

	uc := profile.NewProfileUseCase(ctrl.repos.Profiles)
	id, err := uc.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Get profile
// @Security ApiKeyAuth
// @Tags profile
// @Description get profile
// @ID profile-get
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.Profile
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /profile/:id [get]
func (ctrl *Controller) getProfile(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	p, err := ctrl.repos.Profiles.GetByUserID(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
// @Param input body entity.Profile true "profile data"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /profile [put]
func (ctrl *Controller) updateProfile(c *gin.Context) {
	var input entity.Profile

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	input.UserID = userId

	err = ctrl.repos.Profiles.Update(&input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{true})
}
