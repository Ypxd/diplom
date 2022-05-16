package v1

import (
	"github.com/Ypxd/diplom/auth/internal/models"
	"github.com/Ypxd/diplom/auth/internal/service"
	"github.com/Ypxd/diplom/auth/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getAllEvents(c *gin.Context) {
	var (
		err    error
		result models.EventsResponse
	)

	result, err = h.services.Events.GetAllEvents(c.Request.Context())
	if err != nil {
		response(c, http.StatusInternalServerError, err, result, nil)
		return
	}

	response(c, http.StatusOK, err, result, nil)
	return
}

func (h *Handler) getEvent(c *gin.Context) {
	var (
		err    error
		token  models.JWTToken
		tokenS string
		result []models.MyEvents
		req    []models.AllTags
	)

	err = c.BindJSON(&req)
	if err != nil {
		response(c, http.StatusBadRequest, err, "", nil)
		return
	}

	jwtToken := c.GetHeader("Token")

	token, tokenS, err = service.ParseJWT(utils.GetConfig().Auth.TokenSecret, jwtToken)
	if err != nil {
		response(c, http.StatusUnauthorized, err, result, nil)
		return
	}
	c.Header("token", tokenS)

	result, err = h.services.Events.GetEvents(c.Request.Context(), req, token.UserID)
	if err != nil {
		response(c, http.StatusInternalServerError, err, result, nil)
		return
	}

	response(c, http.StatusOK, err, result, nil)
	return
}

func (h *Handler) initEventsRoutes(api *gin.RouterGroup) {
	group := api.Group("/events")
	{
		group.POST("/", h.getAllEvents)
		group.POST("/event_by_tag", h.getEvent)
	}
}
