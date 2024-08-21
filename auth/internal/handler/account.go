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

	err := conn.DB.First(&account, c.FormValue("account")).Error
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &model.ErrMsg{Message: "Authentication Failure"})
	}

	// Password Verification
	err = bcrypt.CompareHashAndPassword([]byte(account.Passwd), []byte(c.FormValue("passwd")))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &model.ErrMsg{Message: "Authentication Failure"})
	}

	// Issue JWT
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

func (conn *Handler) Register(c echo.Context) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(c.FormValue("passwd")), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &model.ErrMsg{Message: "Password Hashing Failed"})
	}

	account := &model.Account{
		Name:        c.FormValue("name"),
		Passwd:      string(hash),
		Description: c.FormValue("Description"),
	}
	err = conn.DB.First(&account, c.FormValue("account")).Error
	if err == nil {
		return c.JSON(http.StatusBadRequest, &model.ErrMsg{Message: "Account Already Exist"})
	}

	result := conn.DB.Create(&account)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, &model.ErrMsg{Message: "Account Register Failed"})
	}

	return c.JSON(http.StatusOK, "ok")
}
