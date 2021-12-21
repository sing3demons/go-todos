package handler

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/go-todos/model"
	"github.com/sing3demons/go-todos/service"
)

type TodoHandler interface {
	AllTodos(c *fiber.Ctx) error
	FindTodos(c *fiber.Ctx) error
	CreateTodos(c *fiber.Ctx) error
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

func (h *todoHandler) FindTodos(c *fiber.Ctx) error {
	id, err := h.findTodoByID(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err)
	}
	todo, err := h.service.FindTodo(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err)
	}
	return c.Status(fiber.StatusOK).JSON(todo)
}

func (h *todoHandler) findTodoByID(c *fiber.Ctx) (uint, error) {
	uid, err := strconv.ParseUint(c.Params("id"), 0, 0)
	if err != nil {
		return 0, err
	}
	id := uint(uid)
	return id, nil
}

type insertTodo struct {
	Title string `form:"title"`
	Desc  string `form:"desc"`
	Image string `form:"image"`
}

func (h *todoHandler) CreateTodos(c *fiber.Ctx) error {
	var form insertTodo
	c.BodyParser(&form)

	var todo model.Todo
	image, err := h.uploadImage(c, "todo", &todo)
	if err != nil {

		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}
	form.Image = image
	copier.Copy(&todo, &form)

	err = h.service.CreateTodo(todo)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Something went wrong"})
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (h *todoHandler) uploadImage(c *fiber.Ctx, name string, todo *model.Todo) (string, error) {
	file, err := c.FormFile("image")
	if err != nil || file == nil {
		log.Println(err)
		return "", err
	}

	n := time.Now().Unix()
	s := strconv.FormatInt(n, 10)
	filename := "uploads/" + name + "/" + "images" + "/" + strings.Replace(s, "-", "", -1)
	// extract image extension from original file filename
	fileExt := strings.Split(file.Filename, ".")[1]
	// generate image from filename and extension
	image := fmt.Sprintf("%s.%s", filename, fileExt)

	if err := c.SaveFile(file, image); err != nil {
		return "", err
	}

	return image, nil
}
