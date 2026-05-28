package dto

import "github.com/google/uuid"

type AdminLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminParam struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

type EmailMessage struct {
	Subject string   `json:"subject"`
	Content string   `json:"content"`
	Name    []string `json:"name"`
}

type AdminLoginResponse struct {
	Token string `json:"token"`
}
