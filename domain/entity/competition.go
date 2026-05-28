package entity

type Competition struct {
	ID            int            `json:"id" gorm:"autoIncrement;primaryKey;type:integer"`
	Name          string         `json:"name" gorm:"column:competition_name;type:varchar(70);uniqueIndex"`
	Desc          string         `json:"competition_desc" gorm:"column:competition_description;type:text"`
	LinkWA        string         `json:"competition_link_wa" gorm:"column:competition_link_wa"`
	Teams         []Team         `json:"teams" gorm:"foreignKey:competition_id"`
	Announcements []Announcement `json:"announcements" gorm:"foreignKey:competition_id"`
}

func NewCompetition(name string, desc string) *Competition {
	return &Competition{
		Name: name,
		Desc: desc,
	}
}
