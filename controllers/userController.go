package controllers

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xbt573/project-example/models"
	"github.com/xbt573/project-example/services"
)

type UserController interface {
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Refresh(ctx *fiber.Ctx) error
}

func NewUserController(userService services.UserService, validatorObj *validator.Validate) UserController {
	return &concreteUserController{userService, validatorObj}
}

type concreteUserController struct {
	userService services.UserService
	validator   *validator.Validate
}

func (c *concreteUserController) Register(ctx *fiber.Ctx) error {
	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := c.validator.Struct(user); err != nil {
		var errors map[string]string

		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Error()
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	token, err := c.userService.Register(user.Login, user.Password)
	if err != nil {
		if errors.Is(err, services.ErrUserServiceAlreadyExists) {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "user already exists",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"token": token,
	})
}

func (c *concreteUserController) Login(ctx *fiber.Ctx) error {
	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := c.validator.Struct(user); err != nil {
		var errors map[string]string

		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Error()
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	token, err := c.userService.Login(user.Login, user.Password)
	if err != nil {
		if errors.Is(err, services.ErrUserServiceUnauthorized) {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "wrong password or username",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"token": token,
	})
}

func (c *concreteUserController) Refresh(ctx *fiber.Ctx) error {
	token := ctx.Locals("user").(*jwt.Token)

	tokens, err := c.userService.Refresh(token.Raw)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"token": tokens,
	})
}
