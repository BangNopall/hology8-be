package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/hology8/hology-be/domain/enums"
)

type Team struct {
	ID                  string            `json:"id" gorm:"type:varchar(100);primaryKey"`
	LeaderID            uuid.UUID         `json:"-" gorm:"type:uuid;uniqueIndex"`
	CompetitionID       int               `json:"-" gorm:"type:integer"`
	UniversityID        int               `json:"university_id" gorm:"type:integer;default:null"`
	Name                string            `json:"team_name" gorm:"column:team_name;type:varchar(170);uniqueIndex"`
	JoinToken           string            `json:"join_token" gorm:"type:varchar(20);uniqueIndex"`
	SenderPaymentName   string            `json:"sender_payment_name" gorm:"type:varchar(80)"`
	BankAccountNumber   string            `json:"bank_account_number" gorm:"type:varchar(80)"`
	BankName            string            `json:"bank_name" gorm:"type:varchar(80)"`
	PaymentProofLink    string            `json:"payment_proof_link" gorm:"type:varchar(255)"`
	TwibbonProofLink    string            `json:"twibbon_proof_link" gorm:"type:varchar(255)"`
	ProposalDocLink     string            `json:"proposal_doc_link" gorm:"type:varchar(255)"`
	StatementLetterLink string            `json:"statement_letter_link" gorm:"type:varchar(255)"`
	VideoLink           string            `json:"video_link" gorm:"type:varchar(255)"`
	PPTLink             string            `json:"ppt_link" gorm:"type:varchar(255)"`
	Status              enums.Status      `json:"status" gorm:"type:status"`
	Phase               enums.Phase       `json:"phase" gorm:"type:phase"`
	WinnerPlace         enums.WinnerPlace `json:"winner_place" gorm:"type:winner_place"`
	CreatedAt           time.Time         `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt           time.Time         `json:"updated_at" gorm:"type:timestamp"`
	Leader              User              `json:"leader" gorm:"foreignKey:leader_id"`
	Competition         Competition       `json:"competition" gorm:"foreignKey:competition_id"`
	University          University        `json:"university" gorm:"foreignKey:university_id"`
	Voucher             *Voucher          `json:"voucher" gorm:"foreignKey:voucher_id"`
	VoucherID           *string           `json:"voucher_id" gorm:"type:varchar(36);default:null"`
	Announcements       []Announcement    `json:"announcements" gorm:"foreignKey:team_id"`
	Members             []DetailTeams     `json:"members" gorm:"foreignKey:team_id"`
}

type DetailTeams struct {
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;index:idx_team_members,unique"`
	TeamID    string    `json:"team_id" gorm:"type:varchar(50);index:idx_team_members,unique"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;default:current_timestamp"`
	User      User      `json:"user" gorm:"foreignKey:user_id"`
	Team      Team      `json:"team" gorm:"foreignKey:team_id"`
}

func NewTeam(
	teamId string,
	leaderId uuid.UUID,
	compeId int,
	name string,
	joinToken string,
	status enums.Status,
	phase enums.Phase,
	winnerPlace enums.WinnerPlace,
) *Team {
	return &Team{
		ID:            teamId,
		LeaderID:      leaderId,
		CompetitionID: compeId,
		Name:          name,
		JoinToken:     joinToken,
		Status:        status,
		Phase:         phase,
		WinnerPlace:   winnerPlace,
	}
}
