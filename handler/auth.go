package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/go-todos/middleware"
	"github.com/sing3demons/go-todos/model"
	"github.com/sing3demons/go-todos/service"
)

type UserHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Profile(c *fiber.Ctx) error
	FindUsers(c *fiber.Ctx) error
}

type userHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) UserHandler {
	return &userHandler{
		service: service,
	}
}

type Register struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6,max=32"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
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

type formLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

func (h *userHandler) Login(c *fiber.Ctx) error {
	u := new(formLogin)
	if err := c.BodyParser(&u); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	errors := ValidateStruct(*u)
	if errors != nil {
		return c.JSON(errors)

	}

	user, err := h.service.FindByEmail(u.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}
	if err := compareHashAndPassword(user.Password, u.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

func (h *userHandler) Profile(c *fiber.Ctx) error {
	user := c.Locals("sub").(model.User)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	return c.Status(fiber.StatusOK).JSON(serializedUser)
}



func (h *userHandler) FindUsers(c *fiber.Ctx) error {
	users, err := h.service.FindByUsers()
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err)
	}

	serializedUser := []userResponse{}
	copier.Copy(&serializedUser, &users)

	return c.Status(fiber.StatusOK).JSON(serializedUser)

}
