package handler

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Pagination struct {
	Limit      int         `json:"limit,omitempty"`
	Page       int         `json:"page,omitempty"`
	Sort       string      `json:"sort,omitempty"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Rows       interface{} `json:"rows"`
}

func (h *todoHandler) removeImage(path string) error {
	if path != "" {
		pwd, _ := os.Getwd()
		os.Remove(pwd + "/" + path)
	}
	return nil
}

func (h *todoHandler) findByID(c *fiber.Ctx) (uint, error) {
	uid, err := strconv.ParseUint(c.Params("id"), 0, 0)
	if err != nil {
		return 0, err
	}
	id := uint(uid)
	return id, nil
}

func (h *todoHandler) uploadImage(c *fiber.Ctx, name string) (string, error) {
	file, err := c.FormFile("image")
	if err != nil || file == nil {
		log.Println(err)
		return "", err
	}
	m := time.Now().UnixMilli()
	n := time.Now().Unix() + m
	s := strconv.FormatInt(n, 12)
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

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateStruct(user interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

type userResponse struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
	Role      string `json:"role"`
}
