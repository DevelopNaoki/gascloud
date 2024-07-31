package main

import (
	"fmt"
	"os"

	"github.com/DevelopNaoki/gascloud/auth/internal/config"
	"github.com/spf13/cobra"
)

var confpath string
var RootCmd = &cobra.Command{
	Use:   "gascloud-auth",
	Short: "gascloud authorized api server",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.LoadConfigFile(confpath)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", config)
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
	// Initializing cobra and setting commands
	cobra.OnInitialize()
	RootCmd.Flags().StringVarP(&confpath, "config-file", "c", "", "Specify Config File Path")
}
