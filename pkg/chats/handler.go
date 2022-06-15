package chat

// TODO: Make response for places (get from ML server)

import (
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

type response_message struct {
	Message string `json:"message,omitempty"`
}

func NewHandler(service Service) *Handler {
	return &Handler{
		Service: service,
	}
}

// @Summary      Testing Server
// @Description  Endpoint for testing by returning Hello World
// @Tags         Test
// @Produce      json
// @Success      200  {object}  response_message
// @Failure      400  {object}  response_message
// @Failure      500  {object}  response_message
// @Router       / [post]
func (h *Handler) HelloWorld(c *fiber.Ctx) error {
	// message, status, err := h.Service.HelloWorld()
	status := http.StatusOK

	return c.Status(status).JSON(fiber.Map{
		"message": "Hello World",
	})
}

// @Summary                     Set user data for messages
// @Description                 Endpoint for setting user data
// @Tags                        User
// @Produce                     json
// @Success                     200  {object}  response_message
// @Failure                     400  {object}  response_message
// @Failure                     500  {object}  response_message
// @Router                      /chat/user [post]
// @Security JWT
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

// @Summary      Set group for messages
// @Description  Endpoint for setting group to store messages
// @Tags         Group
// @Produce      json
// @Success      200  {object}  response_message
// @Failure      400  {object}  response_message
// @Failure      500  {object}  response_message
// @Router       /chat/group [post]
// @Security JWT
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

// @Summary      Add message according to group
// @Description  Endpoint for sending messages according to group
// @Tags         Message
// @Produce      json
// @Accept       json
// @Param        tag  body      Text_message  true  "Text tag"
// @Success      200  {object}  Text_message
// @Failure      400  {object}  response_message
// @Failure      500  {object}  response_message
// @Router       /chat/message [post]
// @Security JWT
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
	ml_server := os.Getenv("ML_SERVER_URL")

	// determined which response
	switch payload.Tag {
	case "restaurant_recommendation":
		ml_server = ml_server + "/resto-recommendation"
	case "hotel_recommendation":
		ml_server = ml_server + "/hotel-recommendation"
	}
	payload, status, err = h.Service.ProcessedML(payload.Text, ml_server)

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
