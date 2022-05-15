package v1

import (
	"bufio"
	"fmt"
	"github.com/Ypxd/diplom/auth/internal/models"
	"github.com/Ypxd/diplom/auth/internal/service"
	"github.com/Ypxd/diplom/auth/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func (h *Handler) getAllEvents(c *gin.Context) {
	var (
		err    error
		req    models.EventsRequest
		tokenS string
	)

	err = c.BindJSON(&req)
	if err != nil {
		response(c, http.StatusBadRequest, err, "", nil)
		return
	}

	_, tokenS, err = service.ParseJWT(utils.GetConfig().Auth.TokenSecret, req.JWTToken)
	if err != nil {
		response(c, http.StatusUnauthorized, err, "unauthorized", nil)
		return
	}
	c.Header("token", tokenS)

	res, err := h.services.Events.GetAllEvents(c.Request.Context())
	if err != nil {
		response(c, http.StatusInternalServerError, err, "", nil)
		return
	}

	/*
		file, err := os.Open(res[0].PNG)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		fileInfo, _ := file.Stat()
		var size int64 = fileInfo.Size()
		by := make([]byte, size)
		buffer := bufio.NewReader(file)
		_, err = buffer.Read(by)
		img, _, _ := image.Decode(bytes.NewReader(by))
		out, err := os.Create("qwerty.png")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = png.Encode(out, img)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	*/
	var result []models.EventsResponse
	for _, r := range res {
		var rs models.EventsResponse
		rs.Title = r.Title
		rs.Address = r.Address

		file, err := os.Open(res[0].PNG)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		fileInfo, _ := file.Stat()
		var size int64 = fileInfo.Size()
		by := make([]byte, size)
		buffer := bufio.NewReader(file)
		_, err = buffer.Read(by)

		rs.PNG = by
		result = append(result, rs)
	}
	response(c, http.StatusOK, err, result, nil)
	return
}

func (h *Handler) initEventsRoutes(api *gin.RouterGroup) {
	group := api.Group("/events")
	{
		group.POST("/", h.getAllEvents)
	}
}
