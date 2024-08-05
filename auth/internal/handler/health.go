package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func (conn *Handler) Login(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}
