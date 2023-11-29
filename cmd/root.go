package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xbt573/project-example/api"
	"github.com/xbt573/project-example/database"
	"github.com/xbt573/project-example/services"
	"log/slog"
	"os"
)

var (
	config string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&config, "config", "", "Config file location")
	rootCmd.PersistentFlags().String("listen_url", "localhost:3000", "host:port for http server to listen on")
	rootCmd.PersistentFlags().String("database_url", "", "PostgreSQL database url")

	viper.BindEnv("listen_url", "LISTEN_URL")
	viper.BindEnv("database_url", "DATABASE_URL")

	viper.BindPFlag("listen_url", rootCmd.PersistentFlags().Lookup("listen_url"))
	viper.BindPFlag("database_url", rootCmd.PersistentFlags().Lookup("database_url"))

	cobra.OnInitialize(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/project-example")
		viper.AddConfigPath("$XDG_CONFIG_HOME/project-example")
		viper.AddConfigPath("$HOME/.config/project-example")

		if config != "" {
			viper.SetConfigFile(config)
		}

		err := viper.ReadInConfig()
		if err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				return
			}

			slog.Error("Failed to read config!", slog.String("err", err.Error()))
			os.Exit(1)
		}
	})
}

var rootCmd = &cobra.Command{
	Use:   "project-example",
	Short: "project-example is a simple todo service",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			listenUrl   = viper.GetString("listen_url")
			databaseUrl = viper.GetString("database_url")
		)

		if databaseUrl == "" {
			slog.Error("Invalid database url!")
			os.Exit(1)
		}

		err := database.InitDatabase(databaseUrl)
		if err != nil {
			slog.Error("Failed to init db!", slog.String("err", err.Error()))
			os.Exit(1)
		}

		services.InitTodoService(database.GetInstance())

		server := api.NewAPI()
		if err := server.Listen(listenUrl); err != nil {
			slog.Error("Failed serving!", slog.String("err", err.Error()))
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Failed to start cmd!", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
