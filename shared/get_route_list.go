package shared

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type route struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Handler string `json:"handler"`
}

func getRouteList(r *gin.Engine) {
	r.GET("/routes", func(c *gin.Context) {
		resp := make([]route, 0)
		for _, r := range r.Routes() {
			resp = append(resp, route{
				Method:  r.Method,
				Path:    r.Path,
				Handler: r.Handler,
			})
		}

		c.JSON(http.StatusOK, resp)
	})
}
