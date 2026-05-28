package entity

import (
	"time"

	"github.com/google/uuid"
)

type Log struct {
	ID        int       `json:"id" gorm:"type:integer;autoIncrement;primaryKey"`
	AdminID   uuid.UUID `json:"-" gorm:"type:uuid"`
	Action    string    `json:"action" gorm:"type:varchar(255)"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;default:current_timestamp"`
	Admin     Admin     `json:"admin" gorm:"foreignKey:admin_id"`
}

func NewLog(adminID uuid.UUID, action string) *Log {
	return &Log{
		AdminID: adminID,
		Action:  action,
	}
}
