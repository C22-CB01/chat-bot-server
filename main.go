package main

import (
	"context"
	"math/rand"
	"time"

	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"gitlab.com/c22-cb01/chat-bot-server/internal/logger"
	"gitlab.com/c22-cb01/chat-bot-server/pkg/api"
	"go.uber.org/zap"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	godotenv.Load()
	logger.SetLogger()

	ctx := context.Background()
	firebase_app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		zap.L().Fatal("error initializing app: %v\n", zap.Error(err))
	}
	firestore, err := firebase_app.Firestore(ctx)
	if err != nil {
		zap.L().Fatal("error initializing firestore: %v\n", zap.Error(err))
	}

	defer func() {
		firestore.Close()
	}()
	s := api.MakeServer(firebase_app, firestore)
	s.RunServer()
}
