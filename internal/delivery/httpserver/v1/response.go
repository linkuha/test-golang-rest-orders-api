package v1

import (
	"github.com/gin-gonic/gin"
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

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	log.Error().Msg(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{Success: false, Message: message})
}

func newDataResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, dataResponse{Success: true, Data: data})
}
