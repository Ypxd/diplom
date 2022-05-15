package v1

import (
	"github.com/Ypxd/diplom/auth/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}
func (h *Handler) InitRoutes(api *gin.RouterGroup) {
	v1 := api.Group("")
	{
		h.initAuthRoutes(v1)
		h.initEventsRoutes(v1)
		h.initTagsRoutes(v1)
	}
}
