package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/DevelopNaoki/gascloud/auth/internal/config"
	"github.com/DevelopNaoki/gascloud/auth/internal/handler"
	"github.com/DevelopNaoki/gascloud/auth/internal/model"
	"github.com/DevelopNaoki/gascloud/auth/internal/repository"
	"github.com/spf13/cobra"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var configPath string
var RootCmd = &cobra.Command{
	Use:   "gascloud-auth",
	Short: "gascloud authorized api server",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("config loading...")
		// parse to struct and validate configuration file
		c, err := config.LoadConfigFile(configPath)
		if err != nil {
			return err
		}
		fmt.Printf("complete\n")

		// initialize database and share connection
		fmt.Printf("connecting database...")
		db, err := repository.ConnectionDB(c.DB)
		if err != nil {
			return err
		}
		fmt.Printf("complete\n")
		conn := &handler.Handler{DB: db, Secret: c.API.Secret}
		db.AutoMigrate(&model.Account{}, &model.Role{}, &model.RoleBind{}, &model.Permission{}, &model.PermissionBind{}, &model.ServiceCatalog{})
		fmt.Printf("initialize database success\n")

		// setup and run the api server
		e := echo.New()

		e.Use(middleware.Logger())
		e.Use(middleware.Recover())

		api := e.Group(c.API.Prefix)
		api.GET("/health", handler.Health)
		api.POST("/account/login", conn.Login)

		apiAccount := api.Group("/account")
		apiAccount.Use(middleware.JWT([]byte(conn.Secret)))
		apiAccount.POST("/register", conn.Register)

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
	RootCmd.Flags().StringVarP(&configPath, "config", "c", "", "Specify Config File Path")
}
