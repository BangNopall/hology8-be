package contracts

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
	"gorm.io/gorm"
)

type TeamRepository interface {
	Begin() error
	Commit() error
	Rollback() error
	FetchAllTeams(ctx context.Context, params *dto.TeamParams, relations ...string) ([]entity.Team, error)
	FetchAllByConditionAndRelation(
		ctx context.Context,
		condition string,
		args []interface{},
		order string,
		pageParam *dto.PaginationRequest,
		preload ...string,
	) ([]entity.Team, dto.PaginationResponse, error)
	FetchMemberTeams(ctx context.Context, userId uuid.UUID) ([]entity.DetailTeams, error)
	FetchOneByID(ctx context.Context, id string) (entity.Team, error)
	FetchOneByParams(ctx context.Context, params *dto.TeamParams, relations ...string) (entity.Team, error)
	FetchTeamMember(ctx context.Context, teamId string, userId uuid.UUID) (entity.DetailTeams, error)
	InsertTeam(ctx context.Context, team *entity.Team) error
	InsertTeamMember(ctx context.Context, teamMember *entity.DetailTeams) error
	UpdateTeam(ctx context.Context, id string, team *dto.TeamUpdate) error
	LinkVoucher(ctx context.Context, team *entity.Team, voucher *entity.Voucher, tx *gorm.DB) error
	Count(condition string, params []interface{}, groupBy string) (int64, error)
	DeleteTeamMember(ctx context.Context, member *entity.DetailTeams) error
}

type TeamService interface {
	FetchTeamData(ctx context.Context, teamId string) (dto.TeamResponse, error)
	FetchStatisticData() (dto.TeamCounter, error)
	FetchUserTeams(ctx context.Context, userId uuid.UUID) (dto.UserTeamsResponse, error)
	FetchAll(ctx context.Context, params *dto.TeamParams, pageParam *dto.PaginationRequest) (dto.TeamPaginationResponse, error)
	CreateTeam(ctx context.Context, leaderID uuid.UUID, team dto.TeamRegister) error
	UploadPaymentProof(ctx context.Context, id string, userId string, paymentFile *multipart.FileHeader) error
	UploadTwibbonProof(ctx context.Context, id string, userId string, twibbonFile *multipart.FileHeader) error
	UploadProposalDoc(ctx context.Context, id string, userId string, twibbonFile *multipart.FileHeader) error
	UploadStatementLetter(ctx context.Context, id string, userId string, twibbonFile *multipart.FileHeader) error
	UpdateTeamData(ctx context.Context, id string, userId string, team *dto.TeamUpdate) error
	UpdateTeamStatus(ctx context.Context, id string, team *dto.TeamUpdate) error
	UpdateLeader(ctx context.Context, teamId string, leaderId uuid.UUID) error
	RemoveMember(ctx context.Context, member *entity.DetailTeams) error
	JoinTeam(ctx context.Context, joinToken string, userId uuid.UUID) error
	CountTeamNUniv(ctx context.Context) (dto.TeamNUnivCounter, error)
}
