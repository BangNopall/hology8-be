package gomail

import (
	"strconv"

	"gopkg.in/gomail.v2"

	"github.com/hology8/hology-be/internal/infra/env"
	"github.com/hology8/hology-be/pkg/log"
)

type GoMailInterface interface {
	SendEmail(subject string, htmlBody string, email string) error
	SendEmails(subject string, htmlBody string, emails []string) error
}

type GoMailStruct struct {
	host     string
	port     int
	username string
	password string
}

var Gomail = getGoMail()

func getGoMail() GoMailInterface {
	portInt, _ := strconv.Atoi(env.AppEnv.GoMailPort)
	return &GoMailStruct{
		host:     env.AppEnv.GoMailHost,
		port:     portInt,
		username: env.AppEnv.GoMailUsername,
		password: env.AppEnv.GoMailPassword,
	}
}

func (g *GoMailStruct) SendEmail(subject string, htmlBody string, email string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", g.username)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(g.host, g.port, g.username, g.password)

	if err := d.DialAndSend(m); err != nil {
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[GOMAIL][SendEmail]failed to send email")

		return err
	}

	return nil
}

func (g *GoMailStruct) SendEmails(subject string, htmlBody string, emails []string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", g.username)
	m.SetHeader("Subject", subject)
	m.SetHeader("To", "hology@ub.ac.id")
	m.SetHeader("Bcc", emails...)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(g.host, g.port, g.username, g.password)

	if err := d.DialAndSend(m); err != nil {
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[GOMAIL][SendEmail]failed to send email")

		return err
	}

	return nil
}
