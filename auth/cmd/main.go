package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/DevelopNaoki/gascloud/auth/internal/cache"
	"github.com/DevelopNaoki/gascloud/auth/internal/config"
	"github.com/DevelopNaoki/gascloud/auth/internal/handler"
	"github.com/DevelopNaoki/gascloud/auth/internal/repository"
	"github.com/DevelopNaoki/gascloud/auth/internal/route"
	"github.com/spf13/cobra"

	"github.com/gin-gonic/gin"
)

var configPath string
var RootCmd = &cobra.Command{
	Use:   "gascloud-auth",
	Short: "gascloud authorized api server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// parse to struct and validate configuration file
		conf, err := config.LoadConfigFile(configPath)
		if err != nil {
			return err
		}

		// initialize database
		db, err := repository.ConnectionDB(conf.DB)
		if err != nil {
			return err
		}

		// connecting cache server
		cacheClient, err := cache.NewCache(conf.Cache)
		if err != nil {
			return err
		}

		// register db and chache connection
		conn := handler.NewHandler(db, cacheClient, time.Duration(conf.API.Expire)*time.Hour)

		gin.SetMode(gin.ReleaseMode)
		g := gin.New()

		// json format log
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

		// load routing
		api := g.Group(conf.API.Prefix)
		route.Root(api, conn)

		g.Run(conf.API.Address + ":" + strconv.Itoa(conf.API.Port))

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
