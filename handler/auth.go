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

// Register godoc
// @Summary Add an user
// @Description add by json User
// @Tags users
// @Accept  json
// @Produce  json
// @Param Register-Form body Register true "register"
// @Success 201
// @Failure 422 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/v1/auth/sign-up [post]
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

// Login godoc
// @Summary login
// @Description login
// @Tags users
// @Accept  json
// @Produce  json
// @Param Login-Form body formLogin true "login"
// @Success 200 {object} string "token"
// @Failure 422 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/v1/auth/sign-in [post]
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

// Profile godoc
// @Summary get an users
// @Description get by json user
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} userResponse
// @Failure 500 {object} map[string]any
// @Router /api/v1/auth/profile [get]
func (h *userHandler) Profile(c *fiber.Ctx) error {
	user := c.Locals("sub").(model.User)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	return c.Status(fiber.StatusOK).JSON(serializedUser)
}

// FindUsers godoc
// @Summary get an users
// @Description get by json users
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} userResponse
// @Failure 500 {object} map[string]any
// @Router /api/v1/users [get]
func (h *userHandler) FindUsers(c *fiber.Ctx) error {
	users, err := h.service.FindByUsers()
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err)
	}

	serializedUser := []userResponse{}
	copier.Copy(&serializedUser, &users)

	return c.Status(fiber.StatusOK).JSON(serializedUser)

}
