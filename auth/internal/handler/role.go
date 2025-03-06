package handler

import (
	"net/http"

	"github.com/DevelopNaoki/gascloud/auth/internal/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (conn *Handler) Register(c *gin.Context) {
	if c.MustGet("role").(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
		c.Abort()
		return
	}

	var request struct {
		Account     string `json:"account" binding:"required"`
		Passwd      string `json:"password" binding:"required"`
		MailAddr    string `json:"mail" binding:"required"`
		Description string `json:"description"`
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
		Description: request.Description,
	}
	err = conn.DB.Take(&account).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "already exists"})
		c.Abort()
		return
	}

	result := conn.DB.Create(&account)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, &account)
}

func (conn *Handler) RoleShow(c *gin.Context) {
	role := &model.Role{
		Name: c.MustGet("role").(string),
	}

	err := conn.DB.Take(&role).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, &role)
}

func (conn *Handler) RoleList(c *gin.Context) {
	var roles []model.Role
	result := conn.DB.Find(&roles)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, &roles)
}
