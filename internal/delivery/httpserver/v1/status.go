package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/linkuha/test-golang-rest-orders-api/pkg/status"
	"net/http"
)

func (ctrl *Controller) status(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, status.GetStatusInfo())
}
