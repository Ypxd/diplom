package v1

import (
	"errors"
	"github.com/Ypxd/diplom/auth/internal/models"
	"github.com/Ypxd/diplom/auth/internal/service"
	"github.com/Ypxd/diplom/auth/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) changePass(c *gin.Context) {
	var (
		err    error
		token  models.JWTToken
		tokenS string
		req    models.ChangePassReq
	)

	err = c.BindJSON(&req)
	if err != nil {
		response(c, http.StatusBadRequest, err, "", nil)
		return
	}

	switch {
	case req.OldPassword == nil:
		err = errors.New("введите старый пароль")
		response(c, http.StatusBadRequest, err, "", nil)
		return
	case req.NewPassword == nil:
		err = errors.New("введите новый пароль")
		response(c, http.StatusBadRequest, err, "", nil)
		return
	}

	switch {
	case *req.OldPassword == "":
		err = errors.New("введите старый пароль")
		response(c, http.StatusBadRequest, err, "", nil)
		return
	case *req.NewPassword == "":
		err = errors.New("введите новый пароль")
		response(c, http.StatusBadRequest, err, "", nil)
		return
	case *req.NewPassword == *req.OldPassword:
		err = errors.New("пароли должны отличаться")
		response(c, http.StatusBadRequest, err, "", nil)
		return
	}

	jwtToken := c.GetHeader("Token")

	token, tokenS, err = service.ParseJWT(utils.GetConfig().Auth.TokenSecret, jwtToken)
	if err != nil {
		response(c, http.StatusUnauthorized, err, "unauthorized", nil)
		return
	}
	c.Header("token", tokenS)

	err = h.services.Auth.ChangePass(c.Request.Context(), req, token.UserID)
	if err != nil {
		response(c, http.StatusInternalServerError, err, err.Error(), nil)
		return
	}

	response(c, http.StatusOK, err, "", nil)
	return
}

func (h *Handler) me(c *gin.Context) {
	var (
		err      error
		token    models.JWTToken
		tokenS   string
		userInfo models.UserInfo
	)

	jwtToken := c.GetHeader("Token")

	token, tokenS, err = service.ParseJWT(utils.GetConfig().Auth.TokenSecret, jwtToken)
	if err != nil {
		response(c, http.StatusUnauthorized, err, userInfo, nil)
		return
	}
	c.Header("token", tokenS)

	userInfo, err = h.services.Auth.UserInfo(c.Request.Context(), token.UserID)
	if err != nil {
		response(c, http.StatusInternalServerError, err, userInfo, nil)
		return
	}

	response(c, http.StatusOK, err, userInfo, nil)
	return
}

func (h *Handler) auth(c *gin.Context) {
	var (
		err error
		req models.AuthReq
	)

	err = c.BindJSON(&req)
	if err != nil {
		response(c, http.StatusBadRequest, err, "", nil)
		return
	}

	switch {
	case req.Login == nil:
		err = errors.New("введите логин")
		response(c, http.StatusBadRequest, err, "", nil)
		return
	case req.Password == nil:
		err = errors.New("введите пароль")
		response(c, http.StatusBadRequest, err, "", nil)
		return
	}

	switch {
	case *req.Login == "":
		err = errors.New("введите логин")
		response(c, http.StatusBadRequest, err, "", nil)
		return
	case *req.Password == "":
		err = errors.New("введите пароль")
		response(c, http.StatusBadRequest, err, "", nil)
		return
	}

	m, err := h.services.Auth.Auth(c.Request.Context(), req)
	if err != nil {
		response(c, http.StatusInternalServerError, err, m, nil)
		return
	}
	c.Header("token", m)

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
		response(c, http.StatusBadRequest, err, "", nil)
		return
	}
	switch {
	case req.Login == nil:
		err = errors.New("введите логин")
		response(c, http.StatusBadRequest, err, "", nil)
		return
	case req.Password == nil:
		err = errors.New("введите пароль")
		response(c, http.StatusBadRequest, err, "", nil)
		return
	case req.Email == nil:
		err = errors.New("введите почту")
		response(c, http.StatusBadRequest, err, "", nil)
		return
	case req.Name == nil:
		err = errors.New("введите имя")
		response(c, http.StatusBadRequest, err, "", nil)
		return
	case req.Age == nil:
		err = errors.New("введите возраст")
		response(c, http.StatusBadRequest, err, "", nil)
		return
	}

	switch {
	case *req.Login == "":
		err = errors.New("введите логин")
		response(c, http.StatusBadRequest, err, "", nil)
		return
	case *req.Password == "":
		err = errors.New("введите пароль")
		response(c, http.StatusBadRequest, err, "", nil)
		return
	case *req.Email == "":
		err = errors.New("введите почту")
		response(c, http.StatusBadRequest, err, "", nil)
		return
	case *req.Name == "":
		err = errors.New("введите имя")
		response(c, http.StatusBadRequest, err, "", nil)
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
		group.POST("/", h.auth)
		group.POST("/me", h.me)
		group.POST("/change_pass", h.changePass)
	}
}
