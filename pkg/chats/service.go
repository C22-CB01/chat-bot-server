package chat

import "net/http"

type Service interface {
	HelloWorld() (message string, status int, err error)
}

type service struct {
}

func NewService() *service {
	svc := &service{}

	return svc
}

func (s *service) HelloWorld() (message string, status int, err error) {
	message = "Hello World"
	status = http.StatusOK
	err = nil
	return
}
