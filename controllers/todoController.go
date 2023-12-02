package controllers

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xbt573/project-example/models"
	"github.com/xbt573/project-example/services"
	"strconv"
)

type TodoController interface {
	Create(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
	Find(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}

type concreteTodoController struct {
	todoService services.TodoService
	validator   *validator.Validate
}

func NewTodoController(todoService services.TodoService, validatorObj *validator.Validate) TodoController {
	return &concreteTodoController{todoService, validatorObj}
}

func (c *concreteTodoController) Create(ctx *fiber.Ctx) error {
	var todo models.TODO
	if err := ctx.BodyParser(&todo); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := c.validator.Struct(todo); err != nil {
		errors := map[string]string{}

		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Error()
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	userid := uint(claims["id"].(float64))
	todo.UserID = userid

	todo, err := c.todoService.Create(todo)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(todo)
}

func (c *concreteTodoController) List(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	userid := uint(claims["id"].(float64))

	todos, err := c.todoService.List(userid)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(todos)
}

func (c *concreteTodoController) Find(ctx *fiber.Ctx) error {
	strid := ctx.Params("id")
	id64, err := strconv.ParseUint(strid, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad id",
		})
	}

	id := uint(id64)

	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userid := uint(claims["id"].(float64))

	todo, err := c.todoService.Find(userid, id)
	if err != nil {
		if errors.Is(err, services.ErrTodoServiceNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "not found",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(todo)
}

func (c *concreteTodoController) Update(ctx *fiber.Ctx) error {
	var todo models.TODO
	if err := ctx.BodyParser(&todo); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := c.validator.Struct(todo); err != nil {
		errors := map[string]string{}

		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Error()
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	userid := uint(claims["id"].(float64))
	todo.UserID = userid

	todo, err := c.todoService.Update(todo)
	if err != nil {
		if errors.Is(err, services.ErrTodoServiceNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "not found",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(todo)
}

func (c *concreteTodoController) Delete(ctx *fiber.Ctx) error {
	strid := ctx.Params("id")
	id64, err := strconv.ParseUint(strid, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad id",
		})
	}

	id := uint(id64)

	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userid := uint(claims["id"].(float64))

	todo, err := c.todoService.Delete(userid, id)
	if err != nil {
		if errors.Is(err, services.ErrTodoServiceNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "not found",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(todo)
}
