package handler

import (
	"net/http"
	"time"

	"github.com/DevelopNaoki/gascloud/auth/internal/model"
	//"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (conn *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "request does not contain an access token"})
			c.Abort()
			return
		}

		// search tokens
		var sessions []model.Session
		conn.DB.Where("expired_at > ? AND token = ?", time.Now(), tokenString).Find(&sessions)
		switch len(sessions) {
		case 0:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token signature"})
			c.Abort()
			return
		case 1:
			var accounts []model.Account
			conn.DB.Where("id = ?", sessions[0].Account).Find(&accounts)
			switch len(accounts) {
			case 1:
				var roles []model.Role
				conn.DB.Where("id = ?", accounts[0].Role).Find(&roles)
				switch len(roles) {
				case 1:
					c.Set("uid", sessions[0].Account)
					c.Set("account", accounts[0].Name)
					c.Set("role", roles[0].Name)
					return
				default:
					c.JSON(http.StatusInternalServerError, &model.ErrMsg{Message: "Internal Server Error"})
					c.Abort()
					return
				}
			default:
				c.JSON(http.StatusInternalServerError, &model.ErrMsg{Message: "Internal Server Error"})
				c.Abort()
				return
			}
			return
		default:
			c.JSON(http.StatusInternalServerError, &model.ErrMsg{Message: "Internal Server Error"})
			c.Abort()
			return
		}
		return
	}
}
