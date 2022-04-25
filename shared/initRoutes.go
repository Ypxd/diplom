package shared

import "github.com/gin-gonic/gin"

func NewSharedRoutes(r *gin.Engine) {
	getMetricsRoute(r)
	getHealthcheckRoute(r)
	getRouteList(r)
	return
}
