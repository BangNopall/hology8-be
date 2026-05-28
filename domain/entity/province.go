package entity

type Province struct {
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement;type:integer"`
	Name  string `json:"province_name" gorm:"column:province_name;type:varchar(200);uniqueIndex"`
	Users []User `json:"users" gorm:"foreignKey:province_id"`
}
