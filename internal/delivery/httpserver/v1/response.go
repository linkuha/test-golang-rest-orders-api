package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
)

type errorResponse struct {
	Success bool   `json:"ok"`
	Message string `json:"message"`
}

type statusResponse struct {
	Success bool `json:"ok"`
}

type dataResponse struct {
	Success bool        `json:"ok"`
	Data    interface{} `json:"data"`
}

func newErrorResponse(c *gin.Context, err error) {
	errDetails := handleDomainError(err)

	var lvl zerolog.Level
	switch {
	case errDetails.Code >= 500:
		lvl = zerolog.ErrorLevel
	case errDetails.Code >= 400:
		lvl = zerolog.WarnLevel
	default:
		lvl = zerolog.DebugLevel
	}
	log.WithLevel(lvl).Msgf("client error: %s", errDetails.ClientError)
	log.Info().Msgf("internal error: %s", errDetails.DebugError)

	c.AbortWithStatusJSON(errDetails.Code, errorResponse{Success: false, Message: errDetails.ClientError})
}

func newDataResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, dataResponse{Success: true, Data: data})
}
