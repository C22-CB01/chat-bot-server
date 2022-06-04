package chat

import (
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
