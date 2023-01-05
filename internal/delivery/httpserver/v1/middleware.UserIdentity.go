package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/errs"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/service"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (ctrl *Controller) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, errs.NewErrorWrapper(errs.APIAuthorization, errors.New("empty auth header"), "userIdentity failure"))
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(c, errs.NewErrorWrapper(errs.APIAuthorization, errors.New("invalid auth header"), "userIdentity failure"))
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(c, errs.NewErrorWrapper(errs.APIAuthorization, errors.New("token is empty"), "userIdentity failure"))
		return
	}

	authTokenGenerator := service.AuthTokenGenerator{}
	userId, err := authTokenGenerator.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, errs.NewErrorWrapper(errs.APIAuthorization, err, "userIdentity failure"))
		return
	}

	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (string, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return "", errs.NewErrorWrapper(errs.APIAuthorization, errors.New("user id not found"), "userIdentity failure")
	}

	idStr, ok := id.(string)
	if !ok {
		return "", errs.NewErrorWrapper(errs.APIAuthorization, errors.New("user id is invalid type"), "userIdentity failure")
	}

	return idStr, nil
}
