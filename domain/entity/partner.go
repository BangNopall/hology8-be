package entity

type Partner struct {
	ID            int         `json:"id" gorm:"primaryKey;autoIncrement;type:integer"`
	Name          string      `json:"partner_name" gorm:"type:varchar(200);uniqueIndex;column:partner_name"`
	ImageLink     string      `json:"image_link" gorm:"type:varchar(255)"`
	PartnerTypeID int         `json:"-" gorm:"type:integer"`
	PartnerType   PartnerType `json:"partner_type" gorm:"foreignKey:partner_type_id"`
}

type PartnerType struct {
	ID       int       `json:"id" gorm:"primaryKey;autoIncrement;type:integer"`
	Name     string    `json:"partner_type_name" gorm:"type:varchar(200);uniqueIndex;column:partner_type_name"`
	Partners []Partner `json:"partners" gorm:"foreignKey:partner_type_id"`
}

func NewPartner(
	name string,
	imageLink string,
	partnerTypeId int,
) *Partner {
	return &Partner{
		Name:          name,
		ImageLink:     imageLink,
		PartnerTypeID: partnerTypeId,
	}
}
