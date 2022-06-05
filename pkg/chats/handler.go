package chat

import (
	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		Service: service,
	}
}

// @Summary Testing Endpoint
// @Description Endpoint for testing by returning Hello World
// @Tags /chat
// @Produce json
// @Success 200 {object} string
// @Failure 400
// @Failure 500
// @Router /api/chat [get]
func (h *Handler) HelloWorld(c *fiber.Ctx) error {
	message, status, err := h.Service.HelloWorld()

	if err != nil {
		return c.Status(status).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
	})
}

// @Summary Set user data for messages
// @Description Endpoint for setting user data
// @Tags /chat/user
// @Produce json
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/chat/user [post]
func (h *Handler) CreateUserData(c *fiber.Ctx) error {
	claims := c.Locals("claims")
	message, status, err := h.Service.CreateUserData(claims.(*auth.Token).UID)

	if err != nil {
		return c.Status(status).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
	})
}
