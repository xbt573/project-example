package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xbt573/project-example/api/version1"
)

func NewAPI() *fiber.App {
	app := fiber.New()

	version1.Register(app)

	return app
}
