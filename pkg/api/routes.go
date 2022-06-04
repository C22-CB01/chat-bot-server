package api

import (
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gitlab.com/c22-cb01/chat-bot-server/internal/logger"
	chat "gitlab.com/c22-cb01/chat-bot-server/pkg/chats"
)

func (s *Server) SetupRouter() {
	s.Router.Use(logger.MiddleWare())
	s.Router.Use(recover.New())

	chatService := chat.NewService()
	chatHandler := chat.NewHandler(chatService)

	api := s.Router.Group("/api")
	chat := api.Group("/chat")

	chat.Get("/", chatHandler.HelloWorld)


}
