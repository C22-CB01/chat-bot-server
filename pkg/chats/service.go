package chat

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

type Service interface {
	HelloWorld() (message string, status int, err error)
	CreateUserData(uid string) (message string, status int, err error)
	CreateGroup(uid string) (message string, status int, err error)
}

type service struct {
	FirebaseApp  *firebase.App
	FirebaseAuth *auth.Client
	Firestore    *firestore.Client
}

type users struct {
	UID        string `firestore:"uid,omitempty"`
	Email      string `firestore:"email,omitempty"`
	PhotoURL   string `firestore:"photoURL,omitempty"`
	DisplaName string `firestore:"displaName,omitempty"`
}

type groups struct {
	UID         string    `firestore:"uid,omitempty"`
	DateCreated time.Time `firestore:"dateCreated,serverTimestamp"`
	CreatedBy   string    `firestore:"createdBy,omitempty"`
}

type messages struct {
	GroupID     string    `firestore:"groupID,omitempty"`
	DateCreated time.Time `firestore:"dateCreated,serverTimestamp"`
	Message     string    `firestore:"message,omitempty"`
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
		user := users{
			UID:        u.UID,
			Email:      u.Email,
			PhotoURL:   u.PhotoURL,
			DisplaName: u.DisplayName,
		}
		_, err = s.Firestore.Collection("users").Doc(uid).Set(context.Background(), user)
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

func (s *service) CreateGroup(uid string) (message string, status int, err error) {
	status = http.StatusOK
	u, err := s.FirebaseAuth.GetUser(context.Background(), uid)
	if err != nil {
		message = "Cannot obtain user"
		status = http.StatusInternalServerError
	}

	iter := s.Firestore.Collection("groups").Where("createdBy", "==", u.UID).Limit(1).Documents(context.Background())

	_, err = iter.Next()

	if err != nil {
		ref := s.Firestore.Collection("groups").NewDoc()
		group := groups{
			UID:       ref.ID,
			CreatedBy: u.UID,
		}
		ref.Set(context.Background(), group)

		message = "Group successfully created"

		return
	}
	message = "Group has already been created"
	return
}
