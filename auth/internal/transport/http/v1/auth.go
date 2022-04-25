package v1

import (
	"errors"
	"github.com/Ypxd/diplom/auth/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) auth(c *gin.Context) {
	var (
		err error
		req models.AuthReq
	)

	err = c.BindJSON(&req)
	if err != nil {
		response(c, http.StatusBadRequest, err, nil, nil)
		return
	}

	switch {
	case req.Login == nil:
		err = errors.New("empty Login")
		response(c, http.StatusBadRequest, err, nil, nil)
		return
	case req.Password == nil:
		err = errors.New("empty Password")
		response(c, http.StatusBadRequest, err, nil, nil)
		return
	}

	m, err := h.services.Auth.Auth(c.Request.Context(), req)
	if err != nil {
		response(c, http.StatusInternalServerError, err, m, nil)
		return
	}

	response(c, http.StatusOK, err, m, nil)
	return
}

func (h *Handler) register(c *gin.Context) {
	var (
		err error
		req models.AuthReq
	)

	err = c.BindJSON(&req)
	if err != nil {
		response(c, http.StatusBadRequest, err, nil, nil)
		return
	}
	switch {
	case req.Login == nil:
		err = errors.New("empty Login")
		response(c, http.StatusBadRequest, err, nil, nil)
		return
	case req.Password == nil:
		err = errors.New("empty Password")
		response(c, http.StatusBadRequest, err, nil, nil)
		return
	case req.Email == nil:
		err = errors.New("empty Email")
		response(c, http.StatusBadRequest, err, nil, nil)
		return
	case req.Name == nil:
		err = errors.New("empty Name")
		response(c, http.StatusBadRequest, err, nil, nil)
		return
	case req.Age == nil:
		err = errors.New("empty Age")
		response(c, http.StatusBadRequest, err, nil, nil)
		return
	}

	m, err := h.services.Auth.Register(c.Request.Context(), req)
	if err != nil {
		response(c, http.StatusInternalServerError, err, m, nil)
		return
	}

	response(c, http.StatusOK, err, m, nil)
	return
}

func (h *Handler) initAuthRoutes(api *gin.RouterGroup) {
	group := api.Group("/auth")
	{
		group.POST("/register", h.register)
	}
}
