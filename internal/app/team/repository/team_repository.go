package repository

import (
	"context"
	"errors"
	"math"
	"strings"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/pkg/helpers"
	"github.com/hology8/hology-be/pkg/log"
)

type teamRepository struct {
	conn *gorm.DB
	tx   *gorm.DB
}

func NewTeamRepository(conn *gorm.DB) contracts.TeamRepository {
	return &teamRepository{conn: conn}
}

func (r *teamRepository) Begin() error {
	r.tx = r.conn.Begin()

	if err := r.tx.Error; err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[TEAM REPOSITORY][Begin] failed to start db transaction")

		return err
	}

	return nil
}

func (r *teamRepository) Commit() error {
	if r.tx == nil {
		log.Error(log.LogInfo{
			"error": "transaction attr is null",
		}, "[TEAM REPOSITORY][Commit] failed to commit to transaction, tx is null")

		return errors.New("transaction is null")
	}

	err := r.tx.Commit().Error

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[TEAM REPOSITORY][Commit] failed to commit to transaction")

		return err
	}

	r.tx = nil

	return nil
}

func (r *teamRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}

	err := r.tx.Rollback().Error

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[TEAM REPOSITORY][Commit] failed to rollback the transaction")

		return err
	}

	r.tx = nil

	return nil
}

func (teamRepo *teamRepository) FetchOneByID(ctx context.Context, id string) (entity.Team, error) {
	var team entity.Team

	err := teamRepo.conn.Where("\"teams\".\"id\" = ?", id).First(&team).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Team{}, domain.ErrNotFound
		}

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[TEAM REPOSITORY][FetchOneByID] failed to fetch team")

		return entity.Team{}, err
	}

	return team, nil
}

func (r *teamRepository) FetchAllTeams(ctx context.Context, params *dto.TeamParams, relations ...string) ([]entity.Team, error) {
	preloadConn := r.conn

	for _, relation := range relations {
		preloadConn = preloadConn.Preload(relation)
	}

	teams := make([]entity.Team, 0)

	err := preloadConn.Find(&teams, params).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[TEAM REPOSITORY][FetchOneByID] failed to fetch team")

		return nil, err
	}

	return teams, nil
}

func (r *teamRepository) FetchAllByConditionAndRelation(
	ctx context.Context,
	condition string,
	args []interface{},
	order string,
	pageParam *dto.PaginationRequest,
	preload ...string,
) ([]entity.Team, dto.PaginationResponse, error) {
	var teams []entity.Team
	var pageResp dto.PaginationResponse

	preloadConn := r.conn

	for _, relation := range preload {
		if len(relation) < 1 {
			continue
		}
		preloadConn = preloadConn.Preload(relation)
	}

	if pageParam != nil {
		preloadConn = preloadConn.Offset(pageParam.Offset).Limit(pageParam.Limit)

		var count int64

		r.conn.Model(entity.Team{}).Count(&count)

		pageResp.TotalPages = int(math.Ceil(float64(count) / float64(pageParam.Limit)))
		pageResp.Page = pageParam.Page
	}

	var res *gorm.DB
	if len(args) < 1 {
		res = preloadConn.Order(order).Find(&teams)
	} else {
		res = preloadConn.Where(condition, args...).Order(order).Find(&teams)
	}

	if res.Error != nil {
		log.Error(log.LogInfo{
			"error": res.Error.Error(),
		}, "[TEAM REPOSITORY][FetchAllByConditionAndRelation] failed to fetch teams with relation and condition")

		return nil, pageResp, domain.ErrInternalServer
	}

	if res.RowsAffected < 1 {
		log.Error(log.LogInfo{
			"error": gorm.ErrRecordNotFound.Error(),
		}, "[TEAM REPOSITORY][FetchAllByConditionAndRelation] failed to fetch teams with relation and condition")

		return nil, pageResp, domain.ErrNotFound
	}

	return teams, pageResp, nil
}

func (r *teamRepository) FetchOneByParams(ctx context.Context, params *dto.TeamParams, relations ...string) (entity.Team, error) {
	preloadConn := r.conn

	for _, relation := range relations {
		preloadConn = preloadConn.Preload(relation)
	}

	team := entity.Team{}

	err := preloadConn.Find(&team, params).Error

	if err != nil {

		if err == gorm.ErrRecordNotFound {
			return entity.Team{}, domain.ErrNotFound
		}

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[TEAM REPOSITORY][FetchAllTeams] failed to fetch all teams")

		return entity.Team{}, err
	}

	// for some reason, gorm wont return err record not found
	if team.ID == "" {
		return entity.Team{}, domain.ErrNotFound
	}

	return team, nil
}

func (r *teamRepository) InsertTeam(ctx context.Context, team *entity.Team) error {
	err := r.conn.Create(team).Error

	if err != nil {

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateEntry
		}

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[TEAM REPOSITORY][InsertTeam] failed to insert team")

		return err
	}

	return nil
}

func (r *teamRepository) FetchMemberTeams(ctx context.Context, userId uuid.UUID) ([]entity.DetailTeams, error) {
	var member []entity.DetailTeams

	err := r.conn.
		Preload("Team").
		Where("user_id = ?", userId).
		Preload("Team").
		Preload("Team.Competition").
		Find(&member).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.DetailTeams{}, domain.ErrNotFound
		}

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[TEAM REPOSITORY][FetchTeamMember] failed to find team member")

		return []entity.DetailTeams{}, err
	}

	return member, nil
}

func (r *teamRepository) FetchTeamMember(ctx context.Context, teamId string, userId uuid.UUID) (entity.DetailTeams, error) {
	var member entity.DetailTeams

	err := r.conn.
		Preload("User").
		Preload("Team").
		Where("team_id = ? AND user_id = ?", teamId, userId).
		First(&member).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.DetailTeams{}, domain.ErrNotFound
		}

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[TEAM REPOSITORY][FetchTeamMember] failed to find team member")

		return entity.DetailTeams{}, err
	}

	return member, nil
}

func (r *teamRepository) InsertTeamMember(ctx context.Context, teamMember *entity.DetailTeams) error {
	err := r.conn.Create(teamMember).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateEntry
		}

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[TEAM REPOSITORY][InsertTeamMember] failed to insert member to a team")

		return err
	}

	return nil
}

func (r *teamRepository) UpdateTeam(ctx context.Context, id string, team *dto.TeamUpdate) error {
	conn := r.conn

	if r.tx != nil {
		conn = r.tx
	}

	res := conn.Model(entity.Team{}).Where("id = ?", id).Updates(team)

	if err := res.Error; err != nil {

		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return domain.ErrUniversityNotFound
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateEntry
		}

		if strings.Contains(err.Error(), "invalid input value for enum") {
			return domain.ErrInvalidEnumInput
		}

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[TEAM REPOSITORY][UpdateTeam] failed to update team data")

		return domain.ErrInternalServer
	}

	if err := helpers.CheckRowsAffected(res.RowsAffected); err != nil {
		if err != domain.ErrNotFound {
			log.Error(log.LogInfo{
				"error": err.Error(),
			}, "[TEAM REPOSITORY][UpdateTeam] failed to update team data")
		}

		return err
	}

	return nil
}

func (r *teamRepository) DeleteTeamMember(ctx context.Context, member *entity.DetailTeams) error {
	conn := r.conn

	if r.tx != nil {
		conn = r.tx
	}

	res := conn.Delete(member, "team_id = ? AND user_id = ?", member.TeamID, member.UserID)

	if err := helpers.CheckRowsAffected(res.RowsAffected); err != nil {
		return err
	}

	if err := res.Error; err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[TEAM REPOSITORY][DeleteTeamMember] failed to delete team member")

		return err
	}

	return nil
}

func (r *teamRepository) LinkVoucher(ctx context.Context, team *entity.Team, voucher *entity.Voucher, tx *gorm.DB) error {
	conn := r.conn

	if r.tx != nil {
		conn = r.tx
	}

	if tx != nil {
		conn = tx
	}

	team.VoucherID = &voucher.ID
	team.Voucher = voucher

	res := conn.Updates(team)

	if err := helpers.CheckRowsAffected(res.RowsAffected); err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[TEAM REPOSITORY][LinkVoucher] failed to link voucher")

		return err
	}

	if err := res.Error; err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[TEAM REPOSITORY][LinkVoucher] failed to link voucher")

		return err
	}

	return nil
}

func (r *teamRepository) Count(condition string, params []interface{}, groupBy string) (int64, error) {
	var count int64

	query := r.conn.
		Model(&entity.Team{}).
		Where(condition, params...)

	if groupBy != "" {
		query = query.Group(groupBy)
	}

	res := query.Debug().Count(&count)

	if res.Error != nil {
		log.Error(log.LogInfo{
			"error": res.Error.Error(),
		}, "[TEAM REPOSITORY][Count] failed to count team")

		return 0, res.Error
	}

	return count, nil
}
