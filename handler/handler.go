package handler

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *todoHandler) removeImage(path string) error {
	if path != "" {
		pwd, _ := os.Getwd()
		os.Remove(pwd + "/" + path)
	}
	return nil
}

func (h *todoHandler) uploadImage(c *fiber.Ctx, name string) (string, error) {
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
