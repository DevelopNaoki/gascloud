package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

func Health(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}
