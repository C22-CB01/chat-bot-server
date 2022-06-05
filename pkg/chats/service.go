package chat

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

type Service interface {
	HelloWorld() (message string, status int, err error)
	CreateUserData(uid string) (message string, status int, err error)
}

type service struct {
	FirebaseApp  *firebase.App
	FirebaseAuth *auth.Client
	Firestore    *firestore.Client
}

func NewService(firebase_app *firebase.App, fire_auth *auth.Client, firestore *firestore.Client) *service {
	svc := &service{
		FirebaseApp:  firebase_app,
		FirebaseAuth: fire_auth,
		Firestore:    firestore,
	}

	return svc
}

func (s *service) HelloWorld() (message string, status int, err error) {
	message = "Hello World"
	status = http.StatusOK
	err = nil
	return
}

func (s *service) CreateUserData(uid string) (message string, status int, err error) {
	status = http.StatusOK
	u, err := s.FirebaseAuth.GetUser(context.Background(), uid)
	if err != nil {
		message = "Cannot obtain user"
		status = http.StatusInternalServerError
	}
	_, err = s.Firestore.Collection("users").Doc(uid).Get(context.Background())
	if err != nil {
		_, err = s.Firestore.Collection("users").Doc(uid).Set(context.Background(), map[string]interface{}{
			"uid":         u.UID,
			"email":       u.Email,
			"photoURL":    u.PhotoURL,
			"displayName": u.DisplayName,
		})
		if err != nil {
			message = "Something went wrong when creating user in the collection"
			status = http.StatusInternalServerError
		}
		message = "User successfully created"
        return 
	}
	message = "User has already been created"

	return
}
