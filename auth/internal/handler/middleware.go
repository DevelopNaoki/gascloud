package handler

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/DevelopNaoki/gascloud/auth/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (conn *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			c.Abort()
			return
		}

		var session model.Session
		accountID, err := conn.Cache.Get(tokenString)
		if err != nil {
			err := conn.DB.Where("expired_at > ? AND token = ?", time.Now(), tokenString).Take(&session).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
				c.Abort()
				return
			}
			conn.Cache.SetWithTTL(session.Token, session.Account.String(), time.Until(session.ExpiredAt)/2)
		} else {
			parsedUUID, err := uuid.Parse(accountID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
				c.Abort()
				return
			}
			session.Account = model.UUID(parsedUUID)
			session.Token = tokenString
		}

		value, err := conn.Cache.Get(session.Account.String())
		if err != nil {
			account := &model.Account{
				Common: model.Common{
					ID: session.Account,
				},
			}
			err := conn.DB.Take(&account).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
				c.Abort()
				return
			}
			role := &model.Role{
				Common: model.Common{
					ID: account.Role,
				},
			}
			err = conn.DB.Take(&role).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
				c.Abort()
				return
			}
			value := base64.StdEncoding.EncodeToString([]byte(account.Name + ":" + role.Name))
			conn.Cache.SetWithTTL(session.Account.String(), value, time.Until(session.ExpiredAt)/2)
			c.Set("uid", account.ID)
			c.Set("account", account.Name)
			c.Set("role", role.Name)
		} else {
			decode, err := base64.StdEncoding.DecodeString(value)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
				c.Abort()
				return
			}
			res := strings.SplitN(string(decode), ":", 2)
			c.Set("uid", session.Account.String())
			c.Set("account", res[0])
			c.Set("role", res[1])
		}

		c.Next()
	}
}

func (conn *Handler) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.MustGet("role").(string) != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}
