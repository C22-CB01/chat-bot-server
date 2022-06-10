package api

import (
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
	"gitlab.com/c22-cb01/chat-bot-server/docs"
	"go.uber.org/zap"

	fiberSwagger "github.com/gofiber/swagger"
)

type Server struct {
	Router       *fiber.App
	FirebaseApp  *firebase.App
	FirebaseAuth *auth.Client
	Firestore    *firestore.Client
}

func MakeServer(firebase_app *firebase.App, fire_auth *auth.Client, fire_store *firestore.Client) Server {
	r := fiber.New()
	server := Server{
		Router:       r,
		FirebaseApp:  firebase_app,
		FirebaseAuth: fire_auth,
		Firestore:    fire_store,
	}
	return server
}

func (s *Server) RunServer() {
	s.SetupSwagger()
	s.SetupRouter()

	port := os.Getenv("PORT")
	err := s.Router.Listen(":" + port)
	if err != nil {
		zap.L().Fatal("Failed to listen port "+port, zap.Error(err))
	}
}

func (s *Server) SetupSwagger() {

	docs.SwaggerInfo.Title = "Swagger API"
	docs.SwaggerInfo.Description = "Chat-bot server"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = os.Getenv("SWAGGER_HOST")
	docs.SwaggerInfo.Schemes = []string{"http"}

	s.Router.Get("/swagger/*", fiberSwagger.HandlerDefault)
}
