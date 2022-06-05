package chat

import (
	"encoding/json"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/iterator"
)

type Handler struct {
	Service Service
}

type response_type struct {
	UserId int    `json:"userId"`
	Id     int    `json:"Id"`
	Title  string `json:"title"`
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
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
	// body, err := ioutil.ReadAll(resp.Body)
	var message response_type
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
