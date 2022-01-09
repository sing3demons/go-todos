package middleware

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/go-todos/database"
	"github.com/sing3demons/go-todos/model"
)

func GenerateToken(userID uint) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   strconv.Itoa(int(userID)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Local().Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func verifyToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			str := fmt.Sprintf("Unexpected signing method: %v", t_.Header["alg"])
			err := errors.New(str)
			return nil, err
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
}

func JwtVerify() fiber.Handler {
	return func(c *fiber.Ctx) error {
		bearerToken := c.Get("Authorization")
		if bearerToken == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		authHeader := strings.Split(bearerToken, " ")[1]
		token, _ := verifyToken(authHeader)
		claims, ok := token.Claims.(jwt.MapClaims)

		if ok && !token.Valid {
			return c.SendStatus(fiber.StatusForbidden)
		}

		var user model.User
		id := claims["sub"]

		db := database.GetDB()
		db.First(&user, id)

		c.Locals("sub", user)

		return c.Next()
	}
}
