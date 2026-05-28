package dto

import "time"

type LogResponse struct {
	ID        int       `json:"id"`
	Fullname  string    `json:"fullname"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
}

type LogRequest struct {
	AdminID string `json:"admin_id"`
	Action  string `json:"action"`
}
