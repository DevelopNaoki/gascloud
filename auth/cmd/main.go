package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/DevelopNaoki/gascloud/auth/internal/config"
	"github.com/DevelopNaoki/gascloud/auth/internal/handler"
	"github.com/DevelopNaoki/gascloud/auth/internal/repository"
	"github.com/spf13/cobra"

	"github.com/gin-gonic/gin"
)

var configPath string
var RootCmd = &cobra.Command{
	Use:   "gascloud-auth",
	Short: "gascloud authorized api server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// parse to struct and validate configuration file
		c, err := config.LoadConfigFile(configPath)
		if err != nil {
			return err
		}

		// initialize database and share connection
		db, err := repository.ConnectionDB(c.DB)
		if err != nil {
			return err
		}
		conn := &handler.Handler{
			DB:        db,
			ExpiredAt: time.Duration(c.API.Expire) * time.Hour,
		}

		// setup and run the api server
		gin.SetMode(gin.ReleaseMode)
		g := gin.New()

		// set custom log
		g.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("{\"datetime\":\"%v\", \"status\":\"%3d\", \"latency\":\"%v\", \"from\":\"%s\", \"method\":\"%s\", \"path\":%#v, \"error\":\"%s\"}\n",
				param.TimeStamp.Format("2006-01-02 15:04:05"),
				param.StatusCode,
				param.Latency,
				param.ClientIP,
				param.Method,
				param.Path,
				param.ErrorMessage,
			)
		}))
		g.Use(gin.Recovery())

		api := g.Group(c.API.Prefix)
		api.GET("/health", handler.HealthCheck)

		api.POST("/token/new", conn.IssueToken)

		accounts := g.Group("/account")
		accounts.Use(conn.AuthMiddleware())

		//accountG.Use(middleware.JWT([]byte(conn.Secret)))
		//accountG.POST("/register", conn.Register)
		accounts.GET("/show", conn.Show)

		addr := c.API.Address + ":" + strconv.Itoa(c.API.Port)
		g.Run(addr)

		return nil
	},
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "")
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize()
	RootCmd.Flags().StringVarP(&configPath, "config", "c", "", "Specify Config File Path")
}
