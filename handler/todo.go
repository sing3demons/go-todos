package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/go-todos/service"
)

type TodoHandler interface {
	AllTodos(c *fiber.Ctx) error
}

type todoHandler struct {
	service service.TodoService
}

func NewtodoHandler(service service.TodoService) TodoHandler {
	return &todoHandler{service: service}
}

func (h *todoHandler) AllTodos(c *fiber.Ctx) error {
	todos, err := h.service.FindTodos()
	if err != nil {

		return c.JSON(err)
	}
	return c.Status(fiber.StatusOK).JSON(todos)
}
