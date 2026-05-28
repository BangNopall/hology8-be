package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/BangNopall/hology8-be/domain/entity"
)

type AnnouncementRequest struct {
	ID            int    `json:"id"`
	TeamID        string `json:"team_id"`
	CompetitionID int    `json:"competition_id"`
	Description   string `json:"description" binding:"required"`
	AdminID       uuid.UUID
}

type AnnouncementResponse struct {
	ID            int       `json:"id"`
	AdminID       uuid.UUID `json:"-"`
	TeamID        *string   `json:"-"`
	CompetitionID *int      `json:"-"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func AnnouncementSliceEntityToResponse(announcements []entity.Announcement) []AnnouncementResponse {
	res := make([]AnnouncementResponse, 0)

	for _, a := range announcements {
		res = append(res, AnnouncementEntityToResponse(&a))
	}

	return res
}

func AnnouncementEntityToResponse(a *entity.Announcement) AnnouncementResponse {
	return AnnouncementResponse{
		ID:            a.ID,
		AdminID:       a.AdminID,
		TeamID:        a.TeamID,
		CompetitionID: a.CompetitionID,
		Description:   a.Description,
		CreatedAt:     a.CreatedAt,
	}
}
