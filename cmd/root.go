package cmd

import (
	"github.com/go-playground/validator/v10"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xbt573/project-example/controllers"
	"github.com/xbt573/project-example/database"
	"github.com/xbt573/project-example/services"
	"log/slog"
	"os"
)

func init() {
	rootCmd.PersistentFlags().String("config", "", "Config location")
	rootCmd.PersistentFlags().String("database_url", "", "PostgreSQL database url")
	rootCmd.PersistentFlags().String("listen_url", "", "host:port to listen on")
	rootCmd.PersistentFlags().String("secret", "", "JWT signing secret")

	viper.BindEnv("config", "CONFIG")
	viper.BindEnv("database_url", "DATABASE_URL")
	viper.BindEnv("listen_url", "LISTEN_URL")
	viper.BindEnv("secret", "JWT_SECRET")

	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("database_url", rootCmd.PersistentFlags().Lookup("database_url"))
	viper.BindPFlag("listen_url", rootCmd.PersistentFlags().Lookup("listen_url"))
	viper.BindPFlag("secret", rootCmd.PersistentFlags().Lookup("secret"))

	viper.SetDefault("listen_url", "localhost:3000")
	viper.SetDefault("secret", "asdihuiqhfqif")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$XDG_CONFIG_HOME/project-example")
	viper.AddConfigPath("$HOME/.config/project-example")
	viper.AddConfigPath("/etc/project-example/")

	cobra.OnInitialize(func() {
		if config := viper.GetString("config"); config != "" {
			viper.SetConfigFile(config)
		}

		if err := viper.ReadInConfig(); err != nil {
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
	Short: "service example",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			databaseUrl = viper.GetString("database_url")
			listenUrl   = viper.GetString("listen_url")
			secret      = viper.GetString("secret")
		)

		if databaseUrl == "" {
			slog.Error("database_url cannot be empty!")
			os.Exit(1)
		}

		slog.Info("Starting project-example")

		db, err := database.NewDB(databaseUrl, nil)
		if err != nil {
			slog.Error("Failed to init db!", slog.String("err", err.Error()))
		}

		todoService := services.NewTodoService(db)
		userService := services.NewUserService(db, secret)
		validatorObj := validator.New(validator.WithRequiredStructEnabled())

		todoController := controllers.NewTodoController(todoService, validatorObj)
		userController := controllers.NewUserController(userService, validatorObj)

		app := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})
		jwtAccessMiddleware := jwtware.New(jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: []byte(secret)},
			Claims:     jwt.MapClaims{"subject": "access_token"},
		})

		jwtRefreshMiddleware := jwtware.New(jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: []byte(secret)},
			Claims:     jwt.MapClaims{"subject": "refresh_token"},
		})

		app.Route("/api/v1", func(router fiber.Router) {
			router.Post("/register", userController.Register)
			router.Post("/login", userController.Login)
			router.Get("/refresh", jwtRefreshMiddleware, userController.Refresh)

			router.Route("tasks", func(router fiber.Router) {
				router.Get("/", jwtAccessMiddleware, todoController.List)
				router.Get("/:id", jwtAccessMiddleware, todoController.Find)
				router.Post("/", jwtAccessMiddleware, todoController.Create)
				router.Patch("/", jwtAccessMiddleware, todoController.Update)
				router.Delete("/:id", jwtAccessMiddleware, todoController.Delete)
			})
		})

		slog.Info("Started listening", slog.String("listen_url", listenUrl))
		if err := app.Listen(listenUrl); err != nil {
			slog.Error("Failed to listen!", slog.String("err", err.Error()))
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}
