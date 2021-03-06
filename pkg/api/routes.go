package api

import (
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gitlab.com/c22-cb01/chat-bot-server/internal/firebase"
	"gitlab.com/c22-cb01/chat-bot-server/internal/logger"
	chat "gitlab.com/c22-cb01/chat-bot-server/pkg/chats"
)

func (s *Server) SetupRouter() {
	s.Router.Use(logger.MiddleWare())
	s.Router.Use(recover.New())

	chatService := chat.NewService(s.FirebaseApp, s.FirebaseAuth, s.Firestore)
	chatHandler := chat.NewHandler(chatService)

	api := s.Router.Group("/api")
	api.Post("/", chatHandler.HelloWorld)

	chat := api.Group("/chat")
	chat.Use(firebase.MiddleWare(s.FirebaseAuth))

	chat.Post("/user", chatHandler.CreateUserData)
	chat.Post("/group", chatHandler.CreateGroup)
	chat.Post("/message", chatHandler.CreateMessage)

}
