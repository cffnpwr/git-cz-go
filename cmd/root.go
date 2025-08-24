package cmd

import (
	"fmt"
	"os"

	"github.com/cffnpwr/git-cz-go/config"
	"github.com/cffnpwr/git-cz-go/internal/app"
	"github.com/spf13/cobra"
)

var (
	configPath string
)
var rootCmd = &cobra.Command{
	Use:   "git-cz",
	Short: "Git CZ is Conventional Commit Tool in CLI",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		// Configの読み込み
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %s\n", err)
			os.Exit(1)
		}

		err = app.Run(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running app: %s\n", err)
			os.Exit(1)
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "./config/config.yaml", "config file path")
}
