package health

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	services []Service
}

func NewHandler(services []Service) Handler {
	return Handler{
		services: services,
	}
}

func (h *Handler) Check(c echo.Context) (err error) {
	for i := 0; i < len(h.services); i++ {
		h.services[i].Status, err = request(h.services[i].URL)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, h.services)
}

func request(url string) (status bool, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		status = true
	}

	return
}
