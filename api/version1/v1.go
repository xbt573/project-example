package version1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xbt573/project-example/services"
)

func Register(app *fiber.App) {
	todoService := services.GetTodoServiceInstance()

	app.Get("/hello", Hello)
	app.Get("/api/v1/tasks", ListTasks(todoService))
	app.Get("/api/v1/tasks/:id", SearchTask(todoService))
	app.Post("/api/v1/tasks", CreateTask(todoService))
	app.Patch("/api/v1/tasks", UpdateTask(todoService))
	app.Delete("/api/v1/tasks/:id", DeleteTask(todoService))
}
