package version1

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/xbt573/project-example/models"
	"github.com/xbt573/project-example/services"
	"strconv"
)

func Hello(ctx *fiber.Ctx) error {
	return ctx.SendString("Hello World!")
}

func CreateTask(todoService *services.TodoService) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var task models.TODO
		if err := ctx.BodyParser(&task); err != nil {
			return err
		}

		task, err := todoService.Create(task)
		if err != nil {
			return err
		}

		return ctx.JSON(task)
	}
}

func UpdateTask(todoService *services.TodoService) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var task models.TODO
		if err := ctx.BodyParser(&task); err != nil {
			return err
		}

		task, err := todoService.Update(task)
		if err != nil {
			return err
		}

		return ctx.JSON(task)
	}
}

func ListTasks(todoService *services.TodoService) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		tasks, err := todoService.List()
		if err != nil {
			return err
		}

		return ctx.JSON(tasks)
	}
}

func SearchTask(todoService *services.TodoService) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		strid := ctx.Params("id")
		if strid == "" {
			return errors.New("invalid id")
		}

		id, err := strconv.ParseUint(strid, 10, 64)
		if err != nil {
			return err
		}

		task, err := todoService.Search(uint(id))
		if err != nil {
			return err
		}

		return ctx.JSON(task)
	}
}

func DeleteTask(todoService *services.TodoService) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		strid := ctx.Params("id")
		if strid == "" {
			return errors.New("invalid id")
		}

		id, err := strconv.ParseUint(strid, 10, 64)
		if err != nil {
			return err
		}

		task, err := todoService.Delete(uint(id))
		if err != nil {
			return err
		}

		return ctx.JSON(task)
	}
}
