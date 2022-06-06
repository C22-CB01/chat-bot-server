package chat

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/iterator"
)

type Handler struct {
	Service Service
}

type Text_message struct {
	Text string `json:"message,omitempty"`
	Tag  string `json:"tag"`
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
	// message, status, err := h.Service.HelloWorld()
	ml_server := os.Getenv("ML_SERVER_URL")
	post_body, _ := json.Marshal(map[string]string{
		"text": "Hello",
	})
	response_body := bytes.NewBuffer(post_body)
	resp, err := http.Post(ml_server, "application/json", response_body)
	var message Text_message
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&message)
	status := fiber.StatusOK

	if err != nil {
		return c.Status(status).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	defer resp.Body.Close()

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

// @Summary Set group for messages
// @Description Endpoint for setting group to store messages
// @Tags /chat/group
// @Produce json
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/chat/group [post]
func (h *Handler) CreateGroup(c *fiber.Ctx) error {
	claims := c.Locals("claims")
	message, status, err := h.Service.CreateGroup(claims.(*auth.Token).UID)

	if err != nil && err != iterator.Done {
		return c.Status(status).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
	})
}

// @Summary Add message according to group
// @Description Endpoint for sending messages according to group
// @Tags /chat/message
// @Produce json
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/chat/message [post]
func (h *Handler) CreateMessage(c *fiber.Ctx) error {
	claims := c.Locals("claims")

	// Send client user message
	payload := Text_message{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	message, status, err, groupId := h.Service.CreateMessageUser(claims.(*auth.Token).UID, payload.Text)

	if err != nil && err != iterator.Done {
		return c.Status(status).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// bot message
	payload, status, err = h.Service.ProcessedML(payload.Text)

	message, status, err = h.Service.CreateMessageBot(groupId, payload.Text)

	if err != nil && err != iterator.Done {
		return c.Status(status).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"payload": payload,
	})
}
