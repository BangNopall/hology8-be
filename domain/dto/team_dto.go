package dto

import (
	"time"

	"github.com/BangNopall/hology8-be/domain/entity"
	"github.com/BangNopall/hology8-be/domain/enums"
	"github.com/google/uuid"
)

type TeamRegister struct {
	LeaderID      uuid.UUID `json:"-"`
	CompetitionID int       `json:"competition_id" binding:"numeric,required"`
	Name          string    `json:"team_name" binding:"required,alphnumsympace,min=5,max=50"`
	JoinToken     string    `json:"-"`
}

type TeamUpdate struct {
	LeaderID            uuid.UUID         `json:"leader_id"`
	UniversityID        int               `json:"university_id" binding:"omitempty,numeric"`
	SenderPaymentName   string            `json:"sender_payment_name" binding:"omitempty,max=60,alphnumsympace"`
	BankAccountNumber   string            `json:"bank_account_number" binding:"omitempty,max=60,numeric"`
	BankName            string            `json:"bank_name" binding:"omitempty,max=50,alphnumsympace"`
	ProposalDocLink     string            `json:"-"` // these should be updated via its own endpoint
	VideoLink           string            `json:"video_link" binding:"omitempty,max=255,url"`
	PPTLink             string            `json:"ppt_link" binding:"omitempty,max=255,url"`
	PaymentProofLink    string            `json:"-"` // these should be updated via its own endpoint
	TwibbonProofLink    string            `json:"-"` // these should be updated via its own endpoint
	StatementLetterLink string            `json:"-"` // these should be updated via its own endpoint
	Status              enums.Status      `json:"status"`
	Phase               enums.Phase       `json:"phase"`
	WinnerPlace         enums.WinnerPlace `json:"winner_place"`
}

type TeamResponse struct {
	ID                  string                 `json:"id"`
	LeaderID            uuid.UUID              `json:"-"`
	CompetitionID       int                    `json:"-"`
	UniversityID        int                    `json:"university_id"`
	Name                string                 `json:"team_name"`
	JoinToken           string                 `json:"join_token"`
	SenderPaymentName   string                 `json:"sender_payment_name"`
	BankAccountNumber   string                 `json:"bank_account_number"`
	BankName            string                 `json:"bank_name"`
	PaymentProofLink    string                 `json:"payment_proof_link"`
	TwibbonProofLink    string                 `json:"twibbon_proof_link"`
	ProposalDocLink     string                 `json:"proposal_doc_link"`
	StatementLetterLink string                 `json:"statement_letter_link"`
	VideoLink           string                 `json:"video_link"`
	PPTLink             string                 `json:"ppt_link"`
	Status              enums.Status           `json:"status"`
	Phase               enums.Phase            `json:"phase"`
	WinnerPlace         enums.WinnerPlace      `json:"winner_place"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
	Competition         CompetitionResponse    `json:"competition"`
	University          UniversityResponse     `json:"university"`
	Leader              UserResponse           `json:"leader"`
	Voucher             bool                   `json:"voucher"`
	Members             []UserResponse         `json:"members"`
	Announcements       []AnnouncementResponse `json:"announcements"`
}

type UserTeamsResponse struct {
	Teams []TeamResponse `json:"teams"`
}

type TeamPaginationResponse struct {
	Teams      []TeamResponse     `json:"teams"`
	Pagination PaginationResponse `json:"pagination"`
	Counter    TeamCounter        `json:"counter"`
}

type TeamCounter struct {
	VerifiedTeam     int `json:"verified_team"`
	OnholdTeam       int `json:"onhold_team"`
	EliminatedTeam   int `json:"eliminated_team"`
	FinalTeam        int `json:"final_team"`
	DisqualifiedTeam int `json:"disqualified_team"`
}

type TeamParams struct {
	ID            string
	LeaderID      uuid.UUID
	CompetitionID int
	Name          string
	Status        enums.Status
	Phase         enums.Phase
	WinnerPlace   enums.WinnerPlace
	JoinToken     string
	SortBy        string
}

type TeamNUnivCounter struct {
	TeamCounter int `json:"team_counter"`
	UnivCounter int `json:"univ_counter"`
}

func TeamSliceEntityToResponse(teams []entity.Team) []TeamResponse {
	res := make([]TeamResponse, 0)

	for _, t := range teams {
		res = append(res, TeamEntityToResponse(&t))
	}

	return res
}

func TeamEntityToResponse(team *entity.Team) TeamResponse {

	return TeamResponse{
		ID:                  team.ID,
		LeaderID:            team.LeaderID,
		CompetitionID:       team.CompetitionID,
		UniversityID:        team.UniversityID,
		Name:                team.Name,
		JoinToken:           team.JoinToken,
		SenderPaymentName:   team.SenderPaymentName,
		BankAccountNumber:   team.BankAccountNumber,
		BankName:            team.BankName,
		PaymentProofLink:    team.PaymentProofLink,
		TwibbonProofLink:    team.TwibbonProofLink,
		ProposalDocLink:     team.ProposalDocLink,
		StatementLetterLink: team.StatementLetterLink,
		VideoLink:           team.VideoLink,
		PPTLink:             team.PPTLink,
		Status:              team.Status,
		Phase:               team.Phase,
		WinnerPlace:         team.WinnerPlace,
		CreatedAt:           team.CreatedAt,
		UpdatedAt:           team.UpdatedAt,
		Leader:              *UserEntityToResponseDto(&team.Leader),
		Competition:         CompetitionEntityToResponse(&team.Competition),
		University:          UniversityEntityToDto(&team.University),
		Voucher:             team.VoucherID != nil,
		Members: func() []UserResponse {
			res := make([]UserResponse, 0)

			for _, u := range team.Members {
				res = append(res, *UserEntityToResponseDto(&u.User))
			}

			return res
		}(),
	}
}
