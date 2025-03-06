package handler

import (
	"net/http"
	"time"

	"github.com/DevelopNaoki/gascloud/auth/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (conn *Handler) IssueToken(c *gin.Context) {
	var request struct {
		Account string `json:"account" binding:"required"`
		Passwd  string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		c.Abort()
		return
	}

	account := &model.Account{
		Name: request.Account,
	}
	err := conn.DB.Take(&account).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "authentication failure"})
		c.Abort()
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Passwd), []byte(request.Passwd))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "authentication failure"})
		c.Abort()
		return
	}

	var sessions []model.Session
	conn.DB.Where("expired_at > ? AND account = ?", time.Now(), account.ID).Find(&sessions)
	if len(sessions) > 0 {
		c.JSON(http.StatusOK, gin.H{"token": sessions[0].Token})
		return
	}

	session := &model.Session{
		Token:     uuid.New().String(),
		Account:   account.ID,
		ExpiredAt: time.Now().Add(conn.ExpiredAt),
	}

	result := conn.DB.Create(&session)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": session.Token})
}

func (conn *Handler) AccountRegister(c *gin.Context) {
	var request struct {
		Account     string     `json:"account" binding:"required"`
		Passwd      string     `json:"password" binding:"required"`
		MailAddr    string     `json:"mail_address"`
		Role        model.UUID `json:"role_id" binding:"required"`
		Description string     `json:"description"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		c.Abort()
		return
	}

	hash, err := bcrypt.GenerateFromPassword(([]byte(request.Passwd)), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		c.Abort()
		return
	}

	account := &model.Account{
		Name:        request.Account,
		Passwd:      string(hash),
		MailAddr:    request.MailAddr,
		Role:        request.Role,
		Description: request.Description,
	}

	result := conn.DB.FirstOrCreate(&account, model.Account{Name: request.Account})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "already exists"})
		c.Abort()
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, &account)
}

func (conn *Handler) AccountShow(c *gin.Context) {
	account := &model.Account{
		Name: c.MustGet("account").(string),
	}

	err := conn.DB.Take(&account).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, &account)
}

func (conn *Handler) AccountList(c *gin.Context) {
	var accounts []model.Account
	result := conn.DB.Find(&accounts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, &accounts)
}
