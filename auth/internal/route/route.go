package route

import (
	"github.com/DevelopNaoki/gascloud/auth/internal/handler"

	"github.com/gin-gonic/gin"
)

func Root(api *gin.RouterGroup, conn *handler.Handler) {
	public := api.Group("/public")
	{
		public.GET("/health", handler.HealthCheck)
		public.POST("/token/issue", conn.IssueToken)
	}
	auth := api.Group("/auth")
	{
		auth.Use(conn.AuthMiddleware())

		accounts := auth.Group("/account")
		{
			accounts.GET("/show", conn.AccountShow)
		}

		roles := auth.Group("/role")
		{
			roles.GET("/show", conn.RoleShow)
		}

		admin := auth.Group("/admin")
		{
			admin.Use(conn.AdminMiddleware())

			accounts := admin.Group("/account")
			{
				accounts.POST("/register", conn.AccountRegister)
				accounts.GET("/list", conn.AccountList)
			}

			roles := admin.Group("/role")
			{
				roles.GET("/list", conn.RoleList)
			}
		}
	}
}
