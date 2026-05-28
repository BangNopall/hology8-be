package entity

type Voucher struct {
	ID     string `json:"id" gorm:"type:varchar(32);primaryKey"`
	TeamID string `json:"team_id" gorm:"type:varchar(100);default:null"`
	Team   *Team  `json:"team" gorm:"foreignKey:voucher_id"`
}

func NewVoucher(ID string) *Voucher {
	return &Voucher{
		ID: ID,
	}
}
