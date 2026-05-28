package entity

import (
	"time"

	"github.com/google/uuid"
)

type Presence struct {
	UserID    uuid.UUID `json:"user_id" gorm:"primaryKey;type:uuid"`
	User      User      `json:"user" gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID;references:ID"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;index:idx_presences_created_at"`
}

func NewPresence(userID uuid.UUID) *Presence {
	return &Presence{
		UserID:    userID,
		CreatedAt: time.Now(),
	}
}
