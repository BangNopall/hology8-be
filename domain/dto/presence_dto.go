package dto

import (
	"time"

	"github.com/google/uuid"
)

type PresenceCreateRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required,uuid"`
}

// Response for POST /presences
type PresenceCreateResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	Fullname  string    `json:"fullname"`
	TeamName  string    `json:"team_name"`
	CreatedAt time.Time `json:"created_at"`
}

// Response for GET /presences/:user_id
type PresenceCheckResponse struct {
	UserID    uuid.UUID  `json:"user_id"`
	Exists    bool       `json:"exists"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type PresenceListItem struct {
	UserID    uuid.UUID `json:"user_id"`
	Fullname  string    `json:"fullname"`
	TeamName  string    `json:"team_name"`
	CreatedAt time.Time `json:"created_at"`
}

type PresenceListResponse struct {
	Presences  []PresenceListItem `json:"presences"`
	Pagination PaginationResponse `json:"pagination"`
}
