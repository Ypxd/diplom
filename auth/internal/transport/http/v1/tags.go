package v1

import (
	"github.com/Ypxd/diplom/auth/internal/models"
	"github.com/Ypxd/diplom/auth/internal/service"
	"github.com/Ypxd/diplom/auth/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getAllTags(c *gin.Context) {
	var (
		err  error
		tags []models.AllTags
	)

	tags, err = h.services.Tags.GetAllTags(c.Request.Context())
	if err != nil {
		response(c, http.StatusInternalServerError, err, "", nil)
		return
	}

	response(c, http.StatusOK, err, tags, nil)
	return
}

func (h *Handler) updateUnfavoriteTags(c *gin.Context) {
	var (
		err    error
		token  models.JWTToken
		tokenS string
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
		response(c, http.StatusUnauthorized, err, "unauthrorized", nil)
		return
	}
	c.Header("token", tokenS)

	err = h.services.Tags.UpdateUnfavoriteTags(c.Request.Context(), req, token.UserID)
	if err != nil {
		response(c, http.StatusInternalServerError, err, "", nil)
		return
	}

	response(c, http.StatusOK, err, "", nil)
	return
}

func (h *Handler) getUnfavoriteTags(c *gin.Context) {
	var (
		err    error
		token  models.JWTToken
		tokenS string
		tags   []models.AllTags
	)

	jwtToken := c.GetHeader("Token")

	token, tokenS, err = service.ParseJWT(utils.GetConfig().Auth.TokenSecret, jwtToken)
	if err != nil {
		response(c, http.StatusUnauthorized, err, tags, nil)
		return
	}
	c.Header("token", tokenS)

	tags, err = h.services.Tags.GetUnfavoriteTags(c.Request.Context(), token.UserID)
	if err != nil {
		response(c, http.StatusInternalServerError, err, "", nil)
		return
	}

	response(c, http.StatusOK, err, tags, nil)
	return
}

func (h *Handler) updateFavoriteTags(c *gin.Context) {
	var (
		err    error
		token  models.JWTToken
		tokenS string
		req    models.MyEvents
	)

	err = c.BindJSON(&req)
	if err != nil {
		response(c, http.StatusBadRequest, err, "", nil)
		return
	}

	jwtToken := c.GetHeader("Token")

	token, tokenS, err = service.ParseJWT(utils.GetConfig().Auth.TokenSecret, jwtToken)
	if err != nil {
		response(c, http.StatusUnauthorized, err, "", nil)
		return
	}
	c.Header("token", tokenS)

	err = h.services.Tags.UpdateFavoriteTags(c.Request.Context(), req, token.UserID)
	if err != nil {
		response(c, http.StatusInternalServerError, err, "", nil)
		return
	}

	response(c, http.StatusOK, err, "", nil)
	return
}

func (h *Handler) initTagsRoutes(api *gin.RouterGroup) {
	group := api.Group("/tags")
	{
		group.POST("/", h.getAllTags)
		group.POST("/get_unfavorite_tags", h.getUnfavoriteTags)
		group.POST("/update_unfavorite_tags", h.updateUnfavoriteTags)
		group.POST("/update_favorite_tags", h.updateFavoriteTags)
	}
}
