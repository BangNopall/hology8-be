package entity

type University struct {
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement;type:integer"`
	Name  string `json:"university_name" gorm:"type:varchar(200);uniqueIndex;column:university_name"`
	Users []User `json:"users" gorm:"foreignKey:university_id"`
	Teams []Team `json:"teams" gorm:"foreignKey:university_id"`
}
