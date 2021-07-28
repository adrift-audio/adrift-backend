package utilities

import (
	"errors"
	"log"
	"net/smtp"
	"os"
)

type auth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &auth{username, password}
}

func (a *auth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *auth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

func SendEmail(destination, subject, message string) error {
	from := os.Getenv("MAILER_FROM")
	password := os.Getenv("MAILER_PASSWORD")
	username := os.Getenv("MAILER_USERNAME")

	if from == "" || password == "" || username == "" {
		log.Fatal("Please provide credentials for the mailing service!")
	}

	auth := LoginAuth(username, password)

	to := []string{destination}
	byteMessage := []byte(
		"To: " + destination + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			message,
	)

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		from,
		to,
		byteMessage,
	)
	if err != nil {
		return err
	}

	return nil
}
