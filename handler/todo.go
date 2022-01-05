package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/go-todos/model"
	"github.com/sing3demons/go-todos/cache"
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
	cache   *cache.Cacher
}

func NewtodoHandler(service service.TodoService, cache *cache.Cacher) TodoHandler {
	return &todoHandler{service: service, cache: cache}
}

func (h *todoHandler) AllTodos(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "24"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	str := fmt.Sprintf("::%d::%d", limit, page)
	query1CacheKey := "todo::all" + str
	query2CacheKey := "todo::page" + str

	cacheItems, err := h.cache.MGet([]string{query1CacheKey, query2CacheKey})
	if err != nil {
		return c.JSON(err)
	}

	todoJS := cacheItems[0]
	pageJS := cacheItems[1]

	var todos []model.Todo
	var paging *model.Pagination

	if todoJS != nil && len(todoJS.(string)) > 0 {
		err := json.Unmarshal([]byte(todoJS.(string)), &todos)
		if err != nil {
			h.cache.Del()
			log.Printf("redis: %v", err)
		}
	}

	itemToCaches := map[string]interface{}{}

	if todoJS == nil {
		todos, paging, err = h.service.FindTodos(limit, page)
		if err != nil {
			return c.JSON(err)
		}
		itemToCaches[query1CacheKey] = todos
	}

	if pageJS != nil && len(pageJS.(string)) > 0 {
		err := json.Unmarshal([]byte(pageJS.(string)), &paging)
		if err != nil {
			h.cache.Del(query2CacheKey)
			log.Println(err.Error())
		}
	}

	if pageJS == nil {
		itemToCaches[query2CacheKey] = paging
	}

	if len(itemToCaches) > 0 {
		fmt.Println("MSet")
		timeToExpire := 10 * time.Second
		err := h.cache.MSet(itemToCaches)
		if err != nil {
			log.Println(err.Error())
		}

		// Set time to expire
		keys := []string{}
		for k := range itemToCaches {
			keys = append(keys, k)
		}
		err = h.cache.Expires(keys, timeToExpire)
		if err != nil {
			log.Println(err.Error())
		}
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
