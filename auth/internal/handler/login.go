package handler

import (
	"net/http"
	"time"

	"github.com/DevelopNaoki/gascloud/auth/internal/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	Secret string
}

func (conn *Handler) Login(c echo.Context) error {
	// User Search
	var account model.Account

	err := conn.DB.First(&account, c.FormValue("user")).Error
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &model.ErrMsg{Message: " authentication failure"})
	}

	// Password Verification
	err = bcrypt.CompareHashAndPassword([]byte(account.Passwd), []byte(c.FormValue("pass")))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &model.ErrMsg{Message: " authentication failure"})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = account.Name
	claims["uid"] = account.ID
	claims["expire"] = time.Now().Add(time.Hour * 24).Unix()
	t, err := token.SignedString([]byte(conn.Secret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &model.ErrMsg{Message: "Issue Token Failed"})
	}
	return c.JSON(http.StatusOK, map[string]string{"Token": t})
}
