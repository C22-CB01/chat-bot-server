package api

import (
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gitlab.com/c22-cb01/chat-bot-server/internal/firebase"
	"gitlab.com/c22-cb01/chat-bot-server/internal/logger"
	chat "gitlab.com/c22-cb01/chat-bot-server/pkg/chats"
)

func (s *Server) SetupRouter() {
	s.Router.Use(logger.MiddleWare())

	chatService := chat.NewService()
	chatHandler := chat.NewHandler(chatService)

	api := s.Router.Group("/api")
	chat := api.Group("/chat")
	res := api.Group("/res")
	res.Use(firebase.MiddleWare(s.FirebaseAuth))

	s.Router.Use(recover.New())

	chat.Get("/", chatHandler.HelloWorld)
	res.Get("/", chatHandler.HelloWorld)

}
