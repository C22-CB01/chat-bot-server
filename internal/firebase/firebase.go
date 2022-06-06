package firebase

import (
	"context"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
)

func MiddleWare(auth *auth.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authorization := c.Get("Authorization")
		if authorization == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Missing Authorization",
			})

		}
		response_token := strings.Split(authorization, " ")
		if response_token[0] != "Bearer" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Missing Bearer",
			})
		}

		id_token := response_token[1]

		token, err := auth.VerifyIDTokenAndCheckRevoked(context.Background(), id_token)

		if err != nil || token == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Token is malformed",
			})
		}

        c.Locals("claims", token)

		return c.Next()
	}
}
