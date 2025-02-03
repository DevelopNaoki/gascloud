package handler

import (
	"net/http"
	"time"

	"github.com/DevelopNaoki/gascloud/auth/internal/model"
	//"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	//"github.com/harakeishi/gats"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	DB        *gorm.DB
	ExpiredAt time.Duration
}

func (conn *Handler) IssueToken(c *gin.Context) {
	// Requet Parameter
	var request struct {
		Account string `json:"account" binding:"required"`
		Passwd  string `json:"password" binding:"required"`
	}

	// Bind Request Parameter
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Account Search
	account := &model.Account{
		Name: request.Account,
	}
	err := conn.DB.First(&account).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, &model.ErrMsg{Message: "Authentication Failure"})
		return
	}

	// Password Verification
	err = bcrypt.CompareHashAndPassword([]byte(account.Passwd), []byte(request.Passwd))
	if err != nil {
		c.JSON(http.StatusUnauthorized, &model.ErrMsg{Message: "Authentication Failure"})
		return
	}

	// search tokens
	var sessions []model.Session
	conn.DB.Where("expired_at > ? AND account = ?", time.Now(), account.ID).Find(&sessions)
	if len(sessions) > 0 {
		c.JSON(http.StatusOK, gin.H{"token": sessions[0].Token})
		return
	}

	now := time.Now()
	token := uuid.New()
	session := &model.Session{
		Token:     token.String(),
		Account:   account.ID,
		ExpiredAt: now.Add(conn.ExpiredAt),
	}

	result := conn.DB.Create(&session)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, &model.ErrMsg{Message: "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": session.Token})
	return
}

//	func (conn *Handler) Register(c *gin.Context) error {
//		// Password Hashing
//		hash, err := bcrypt.GenerateFromPassword([]byte(c.FormValue("passwd")), bcrypt.DefaultCost)
//		if err != nil {
//			return c.JSON(http.StatusInternalServerError, &model.ErrMsg{Message: "Password Hashing Failed"})
//		}
//
//		// Search Account
//		account := &model.Account{
//			Name:        c.FormValue("name"),
//			Passwd:      string(hash),
//			Description: c.FormValue("description"),
//		}
//		err = conn.DB.First(&account, c.FormValue("account")).Error
//		if err == nil {
//			return c.JSON(http.StatusBadRequest, &model.ErrMsg{Message: "Account Already Exist"})
//		}
//
//		// Register Account
//		result := conn.DB.Create(&account)
//		if result.Error != nil {
//			return c.JSON(http.StatusInternalServerError, &model.ErrMsg{Message: "Account Register Failed"})
//		}
//
//		return c.JSON(http.StatusOK, "ok")
//	}
func (conn *Handler) Show(c *gin.Context) {
	// Requet Parameter
	var request struct {
		Account string `json:"account" binding:"required"`
	}

	var account model.Account
	if c.MustGet("role").(string) == "admin" && request.Account != "" {
		// Bind Request Parameter
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		account.Name = request.Account
	} else {
		account.Name = c.MustGet("account").(string)
	}

	// Account Search
	err := conn.DB.First(&account).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.ErrMsg{Message: "Failed To Retrieve Account Information"})
		return
	}

	c.JSON(http.StatusOK, &account)
	return
}
