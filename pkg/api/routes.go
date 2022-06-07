package api

import (
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gitlab.com/c22-cb01/chat-bot-server/internal/firebase"
	"gitlab.com/c22-cb01/chat-bot-server/internal/logger"
	chat "gitlab.com/c22-cb01/chat-bot-server/pkg/chats"
)

func (s *Server) SetupRouter() {
	s.Router.Use(logger.MiddleWare())

	chatService := chat.NewService(s.FirebaseApp, s.FirebaseAuth, s.Firestore)
	chatHandler := chat.NewHandler(chatService)

	api := s.Router.Group("/api")
	api.Get("/", chatHandler.HelloWorld)

	chat := api.Group("/chat")
	chat.Use(firebase.MiddleWare(s.FirebaseAuth))

	s.Router.Use(recover.New())

	chat.Post("/user", chatHandler.CreateUserData)
	chat.Post("/group", chatHandler.CreateGroup)
	chat.Post("/message", chatHandler.CreateMessage)

}
