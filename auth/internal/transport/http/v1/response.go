package v1

import (
	"github.com/Ypxd/diplom/auth/internal/models"
	"github.com/Ypxd/diplom/shared"
	"github.com/gin-gonic/gin"
)

func response(c *gin.Context,
	statusCode int,
	err error,
	message interface{},
	count *int) {
	resp := models.HttpResponse{
		ErrorText: "",
		HasError:  false,
		Message:   message,
		Count:     count,
	}
	if err != nil {
		resp.ErrorText = err.Error()
		resp.HasError = true
		shared.GetLogger().Errorf("%s", err.Error())
	}
	c.JSON(statusCode, resp)

}
