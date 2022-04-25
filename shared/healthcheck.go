package shared

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func getHealthcheckRoute(r *gin.Engine) {
	r.GET("/healthcheck", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
}
