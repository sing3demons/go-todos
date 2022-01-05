package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/go-todos/model"
	"github.com/sing3demons/go-todos/service"
)

type Register struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6,max=32"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserHandler interface {
	Register(c *fiber.Ctx) error
}

type userHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) UserHandler {
	return &userHandler{
		service: service,
	}
}

func (h *userHandler) Register(c *fiber.Ctx) error {
	u := new(Register)
	if err := c.BodyParser(&u); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	errors := ValidateStruct(*u)
	if errors != nil {
		return c.JSON(errors)

	}

	var user model.User
	copier.Copy(&user, &u)
	user.Password = user.GenerateEncryptedPassword()

	if err := h.service.Register(user); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Something went wrong"})
	}

	return c.SendStatus(fiber.StatusCreated)
}
