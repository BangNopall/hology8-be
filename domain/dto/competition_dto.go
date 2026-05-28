package dto

import "github.com/BangNopall/hology8-be/domain/entity"

type CompetitionRequest struct {
	ID     int
	Name   string `json:"name"`
	Desc   string `json:"competition_desc"`
	LinkWA string `json:"competition_link_wa"`
}

type CompetitionResponse struct {
	ID            int                    `json:"id"`
	Name          string                 `json:"name"`
	Desc          string                 `json:"competition_desc"`
	LinkWA        string                 `json:"competition_link_wa"`
	Teams         []TeamResponse         `json:"teams,omitempty"`
	Announcements []AnnouncementResponse `json:"announcements,omitempty"`
}

func CompetitionEntityToResponse(c *entity.Competition) CompetitionResponse {
	return CompetitionResponse{
		ID:     c.ID,
		Name:   c.Name,
		Desc:   c.Desc,
		LinkWA: c.LinkWA,
	}
}
