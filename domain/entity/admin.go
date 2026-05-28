package entity

import "github.com/google/uuid"

type Admin struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Fullname      string         `json:"fullname" gorm:"type:varchar(200)"`
	Username      string         `json:"username" gorm:"type:varchar(200)"`
	Password      string         `json:"password" gorm:"type:varchar(250)"`
	RoleID        string         `json:"-" gorm:"type:varchar(20)"`
	Role          Role           `json:"role" gorm:"foreignKey:role_id"`
	Announcements []Announcement `json:"announcements" gorm:"foreignKey:admin_id"`
	Logs          []Log          `json:"logs" gorm:"foreignKey:admin_id"`
}
