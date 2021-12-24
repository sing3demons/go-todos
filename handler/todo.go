package handler

import (
	"mime/multipart"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/go-todos/model"
	"github.com/sing3demons/go-todos/service"
)

type TodoHandler interface {
	AllTodos(c *fiber.Ctx) error
	FindTodo(c *fiber.Ctx) error
	CreateTodo(c *fiber.Ctx) error
	DeleteTodo(c *fiber.Ctx) error
	UpdateTodo(c *fiber.Ctx) error
}

type todoHandler struct {
	service service.TodoService
}

func NewtodoHandler(service service.TodoService) TodoHandler {
	return &todoHandler{service: service}
}

func (h *todoHandler) AllTodos(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "24"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	todos, paging, err := h.service.FindTodos(limit, page)
	if err != nil {
		return c.JSON(err)
	}
	var result Pagination
	copier.Copy(&result, &paging)
	result.Rows = todos
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *todoHandler) FindTodo(c *fiber.Ctx) error {
	id, err := h.findByID(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err)
	}
	todo, err := h.service.FindTodo(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err)
	}
	return c.Status(fiber.StatusOK).JSON(todo)
}

type insertTodo struct {
	Title string                `form:"title" validate:"required"`
	Desc  string                `form:"desc" validate:"required"`
	Image *multipart.FileHeader `form:"image" validate:"required"`
}

func (h *todoHandler) CreateTodo(c *fiber.Ctx) error {
	var form insertTodo
	c.BodyParser(&form)

	var todo model.Todo
	image, err := h.uploadImage(c, "todo")
	if err != nil {

		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}

	copier.Copy(&todo, &form)
	todo.Image = image

	err = h.service.CreateTodo(todo)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Something went wrong"})
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (h *todoHandler) DeleteTodo(c *fiber.Ctx) error {
	id, err := h.findByID(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err)
	}

	todo, err := h.service.FindTodo(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err)
	}

	h.removeImage(todo.Image)

	err = h.service.DeleteTodo(*todo)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

type updateTodo struct {
	Title string `form:"title" json:"title"`
	Desc  string `form:"desc" json:"desc"`
	Image string `form:"image" json:"image"`
}

func (h *todoHandler) UpdateTodo(c *fiber.Ctx) error {

	id, err := h.findByID(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err)
	}

	todo, err := h.service.FindTodo(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err)
	}

	var form updateTodo
	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	copier.Copy(&todo, &form)
	if form.Image == "" {
		image, err := h.uploadImage(c, "todo")
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		h.removeImage(todo.Image)
		todo.Image = image
	}

	if err := h.service.UpdateTodo(*todo); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
