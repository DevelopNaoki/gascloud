package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/DevelopNaoki/gascloud/auth/internal/config"
	"github.com/spf13/cobra"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var confpath string
var RootCmd = &cobra.Command{
	Use:   "gascloud-auth",
	Short: "gascloud authorized api server",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := config.LoadConfigFile(confpath)
		if err != nil {
			return err
		}

		e := echo.New()

		e.Use(middleware.Logger())
		e.Use(middleware.Recover())

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
