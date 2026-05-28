package entity

type Role struct {
	ID     string  `json:"id" gorm:"type:varchar(20)"`
	Name   string  `json:"role_name" gorm:"uniqueIndex;type:varchar(50)"`
	Admins []Admin `json:"admins" gorm:"foreignKey:role_id"`
}
