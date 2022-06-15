package main

import (
	"context"
	"math/rand"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"gitlab.com/c22-cb01/chat-bot-server/internal/logger"
	"gitlab.com/c22-cb01/chat-bot-server/pkg/api"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

// @BasePath  /api/
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	godotenv.Load()
	logger.SetLogger()

	ctx := context.Background()
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIAL_FILE"))

	firebase_app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		zap.L().Fatal("error initializing app: %v\n", zap.Error(err))
	}
	fire_auth, err := firebase_app.Auth(ctx)
	if err != nil {
		zap.L().Fatal("error initializing app: %v\n", zap.Error(err))
	}
	fire_store, err := firebase_app.Firestore(ctx)
	if err != nil {
		zap.L().Fatal("error initializing firestore: %v\n", zap.Error(err))
	}

	defer func() {
		fire_store.Close()
	}()

	s := api.MakeServer(firebase_app, fire_auth, fire_store)
	s.RunServer()
}
