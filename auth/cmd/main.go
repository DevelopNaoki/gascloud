package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/DevelopNaoki/gascloud/auth/internal/config"
	"github.com/DevelopNaoki/gascloud/auth/internal/handler"
	"github.com/DevelopNaoki/gascloud/auth/internal/repository"
	"github.com/spf13/cobra"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var confpath string
var RootCmd = &cobra.Command{
	Use:   "gascloud-auth",
	Short: "gascloud authorized api server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// parse to struct and validate configuration
		c, err := config.LoadConfigFile(confpath)
		if err != nil {
			return err
		}

		// initialize database and share connection
		db, err := repository.ConnectionDB(c.DB)
		if err != nil {
			return err
		}
		conn := &handler.Handler{DB: db}

		// setup and run the api server
		e := echo.New()

		e.Use(middleware.Logger())
		e.Use(middleware.Recover())

		e.GET(c.API.Prefix+"account/login", conn.Login)

		e.Logger.Fatal(e.Start(c.API.Address + ":" + strconv.Itoa(c.API.Port)))

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
	RootCmd.Flags().StringVarP(&confpath, "config", "c", "", "Specify Config File Path")
}
