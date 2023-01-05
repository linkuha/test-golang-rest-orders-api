package v1

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func (ctrl *Controller) customLogRequest(c *gin.Context) {
	reqID := requestid.Get(c)
	ipAddr := c.ClientIP()

	backupLogger := log.Logger
	log.Logger = log.Logger.With().Str("request_id", reqID).Logger()

	start := time.Now()
	log.Info().Msgf("started %s - %s [%s]", c.Request.Method, c.Request.URL, ipAddr)

	c.Next()

	statusCode := c.Writer.Status()
	log.Info().Msgf("completed in %v - %d (%s)", time.Now().Sub(start), statusCode, http.StatusText(statusCode))
	log.Logger = backupLogger
}
