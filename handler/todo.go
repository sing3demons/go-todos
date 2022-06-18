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
	"github.com/sing3demons/go-todos/cache"
	"github.com/sing3demons/go-todos/model"
	"github.com/sing3demons/go-todos/service"
)

type insertTodo struct {
	Title string                `form:"title" validate:"required"`
	Desc  string                `form:"desc" validate:"required"`
	Image *multipart.FileHeader `form:"image"`
}

type TodoHandler interface {
	AllTodos(c *fiber.Ctx) error
	All_Todos(c *fiber.Ctx) error
	FindTodo(c *fiber.Ctx) error
	CreateTodo(c *fiber.Ctx) error
	DeleteTodo(c *fiber.Ctx) error
	UpdateTodo(c *fiber.Ctx) error
}

type todoHandler struct {
	service service.TodoService
	cache   *cache.Cacher
}

func NewTodoHandler(service service.TodoService, cache *cache.Cacher) TodoHandler {
	return &todoHandler{service: service, cache: cache}
}

// AllTodos godoc
// @Summary Show an todos
// @Tags todos
// @Accept  json
// @Produce  json
// @Param page query uint false "page"
// @Param limit query uint false "limit"
// @Success 200 {object} Pagination
// @Router /api/v1/todos [get]
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

// FindTodo godoc
// @Summary Show an todo
// @Description get string by id
// @Tags todos
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path int true "Todo ID"
// @Success 200 {object} todoResponse
// @Failure 404 {object} responseError
// @Router /api/v1/todos/{id} [get]
func (h *todoHandler) FindTodo(c *fiber.Ctx) error {
	id, err := findByID(c)

	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(responseError{
			Status:  fiber.StatusNotFound,
			Message: "ID: invalid",
		})
	}

	todo, err := h.service.FindTodo(id)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(responseError{
			Status:  fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	resp := todoResponse{
		ID:     todo.ID,
		Title:  todo.Title,
		Desc:   todo.Desc,
		Image:  todo.Image,
		UserID: todo.UserID,
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

// CreateTodo godoc
// @Summary Add an todo
// @Description add by form Todo
// @Tags todos
// @Accept  json
// @Produce  json
// @Param title formData string true "title"
// @Param desc formData string true "desc"
// @Param image formData file true "image"
// @Security BearerAuth
// @Success 201
// @Failure 422 {object} responseError
// @Router /api/v1/todos [post]
func (h *todoHandler) CreateTodo(c *fiber.Ctx) error {
	var form insertTodo
	if err := c.BodyParser(&form); err != nil {
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(responseError{
			Status:  fiber.StatusUnprocessableEntity,
			Message: err.Error(),
		})
	}

	errors := ValidateStruct(&form)
	if errors != nil {
		fmt.Println("validate")
		return c.JSON(errors)

	}

	var todo model.Todo
	image, err := uploadImage(c, "todo")
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			responseError{
				Status:  fiber.StatusUnprocessableEntity,
				Message: err.Error(),
			})
	}

	copier.Copy(&todo, &form)
	todo.Image = image

	err = h.service.CreateTodo(todo)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			responseError{
				Status:  fiber.StatusUnprocessableEntity,
				Message: "Something went wrong",
			})
	}
	return c.SendStatus(fiber.StatusCreated)
}

// DeleteTodo godoc
// @Summary delete an todo
// @Description delete by json Todo
// @Tags todos
// @Accept  json
// @Produce  json
// @Param id path int true "Todo ID"
// @Security BearerAuth
// @Success 204
// @Failure 404 {object} responseError
// @Router /api/v1/todos/{id} [delete]
func (h *todoHandler) DeleteTodo(c *fiber.Ctx) error {
	id, err := findByID(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			responseError{
				Status:  fiber.StatusNotFound,
				Message: err.Error(),
			})
	}

	todo, err := h.service.FindTodo(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			responseError{
				Status:  fiber.StatusNotFound,
				Message: err.Error(),
			})
	}

	if err := h.removeImage(todo.Image); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			responseError{
				Status:  fiber.StatusUnprocessableEntity,
				Message: err.Error(),
			})
	}

	if err = h.service.DeleteTodo(*todo); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			responseError{
				Status:  fiber.StatusNotFound,
				Message: err.Error(),
			})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// UpdateTodo godoc
// @Summary update an todo
// @Description update by json Todo
// @Tags todos
// @Accept  json
// @Produce  json
// @Param id path int true "Todo ID"
// @Param title formData string false "title"
// @Param desc formData string false "desc"
// @Param image formData file false "image"
// @Security BearerAuth
// @Success 204
// @Failure 404 {object} responseError
// @Router /api/v1/todos/{id} [put]
func (h *todoHandler) UpdateTodo(c *fiber.Ctx) error {

	id, err := findByID(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			responseError{
				Status:  fiber.StatusNotFound,
				Message: err.Error(),
			})
	}

	todo, err := h.service.FindTodo(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			responseError{
				Status:  fiber.StatusNotFound,
				Message: err.Error(),
			})
	}

	var form updateTodo
	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responseError{
			Status:  fiber.StatusUnprocessableEntity,
			Message: err.Error(),
		})
	}
	copier.Copy(&todo, &form)
	if form.Image == "" {
		image, err := uploadImage(c, "todo")
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responseError{
				Status:  fiber.StatusUnprocessableEntity,
				Message: err.Error(),
			})
		}
		h.removeImage(todo.Image)
		todo.Image = image
	}

	if err := h.service.UpdateTodo(*todo); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			responseError{
				Status:  fiber.StatusUnprocessableEntity,
				Message: err.Error(),
			})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
