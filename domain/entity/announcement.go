package entity

import (
	"time"

	"github.com/google/uuid"
)

type Announcement struct {
	ID            int          `json:"id" gorm:"type:integer;autoIncrement;primaryKey"`
	AdminID       uuid.UUID    `json:"-" gorm:"type:uuid"`
	TeamID        *string      `json:"-" gorm:"type:varchar(100);default:null"`
	CompetitionID *int         `json:"-" gorm:"type:integer;default:null"`
	Description   string       `json:"description" gorm:"type:text"`
	CreatedAt     time.Time    `json:"created_at" gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt     time.Time    `json:"updated_at" gorm:"type:timestamp;default:current_timestamp"`
	Admin         Admin        `json:"admin" gorm:"foreignKey:admin_id"`
	Team          *Team        `json:"team" gorm:"foreignKey:team_id"`
	Competition   *Competition `json:"competition" gorm:"foreignKey:competition_id"`
}

func NewAnnouncement(teamID string, competitionID int, description string, adminID uuid.UUID) *Announcement {
	a := &Announcement{
		TeamID:        &teamID,
		CompetitionID: &competitionID,
		Description:   description,
		AdminID:       adminID,
	}

	if competitionID == 0 {
		a.CompetitionID = nil 
	}

	if teamID == "" {
		a.TeamID = nil 
	}

	return a
}
