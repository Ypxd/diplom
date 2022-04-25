package handler

import (
	"github.com/Ypxd/diplom/auth/internal/service"
	v1 "github.com/Ypxd/diplom/auth/internal/transport/http/v1"
	"github.com/Ypxd/diplom/shared"
	"github.com/gin-gonic/gin"
)

func InitHandlers(service *service.Service) *gin.Engine {
	router := gin.Default()
	router.Use(
		gin.Recovery(),
	)

	//auth := shared.NewJWT(utils.GetConfig().Auth.CookieName, utils.GetConfig().Auth.Debug, rdb)

	handlerV1 := v1.NewHandler(service)
	api := router.Group("api")
	{
		handlerV1.InitRoutes(api)
		//InitRoutes(api, auth)
	}

	shared.NewSharedRoutes(router)

	return router
}
