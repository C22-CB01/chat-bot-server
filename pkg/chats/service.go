package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

type Service interface {
	CreateUserData(uid string) (message string, status int, err error)
	CreateGroup(uid string) (message string, status int, err error)
	CreateMessageUser(uid string, text string) (message string, status int, err error, groupId string)
	CreateMessageBot(groupUID string, text string) (message string, status int, err error)
	ProcessedML(text string, ml_server string) (response Text_message, status int, err error)
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
	SendBy      string    `firestore:"sendBy,omitempty"`
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

func (s *service) CreateMessageUser(uid string, text string) (message string, status int, err error, groupUID string) {
	status = http.StatusOK
	u, err := s.FirebaseAuth.GetUser(context.Background(), uid)
	if err != nil {
		message = "Cannot obtain user"
		status = http.StatusInternalServerError
	}

	// get groupId
	iter := s.Firestore.Collection("groups").Where("createdBy", "==", u.UID).Limit(1).Documents(context.Background())
	doc, err := iter.Next()
	if err != nil {
		s.CreateGroup(uid)
		return
	}
	groupUID = doc.Data()["uid"].(string)

	text_message := messages{
		SendBy:  uid,
		Message: text,
	}

	_, _, err = s.Firestore.Collection("messages").Doc(groupUID).Collection("group_messages").Add(context.Background(), text_message)

	if err != nil {
		status = http.StatusInternalServerError
		message = "Message failed to be sent"

		return
	}
	message = "Message has been sent"
	return
}

func (s *service) CreateMessageBot(groupUID string, text string) (message string, status int, err error) {
	status = http.StatusOK

	text_message := messages{
		SendBy:  "BOT",
		Message: text,
	}

	_, _, err = s.Firestore.Collection("messages").Doc(groupUID).Collection("group_messages").Add(context.Background(), text_message)

	if err != nil {
		status = http.StatusInternalServerError
		message = "Message failed to be sent"

		return
	}
	message = "Message has been sent"
	return
}

func (s *service) ProcessedML(text string, ml_server string) (ml_resp Text_message, status int, err error) {
	status = http.StatusOK
	post_body, _ := json.Marshal(map[string]string{
		"text": text,
	})
	response_body := bytes.NewBuffer(post_body)
	resp, err := http.Post(ml_server, "application/json", response_body)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&ml_resp)
	return
}
