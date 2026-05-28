package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/pkg/helpers"
	"github.com/hology8/hology-be/pkg/log"
)

type competitionRepository struct {
	conn *gorm.DB
}

func NewCompetitionRepository(conn *gorm.DB) contracts.CompetitionRepository {
	return &competitionRepository{conn}
}

func (r *competitionRepository) FetchAll(ctx context.Context) ([]entity.Competition, error) {
	compes := make([]entity.Competition, 0)

	err := r.conn.Find(&compes).Error

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[COMPETITION REPOSITORY][FetchAll] failed to fetch all competitions")

		return nil, err
	}

	return compes, nil
}

func (r *competitionRepository) FetchAllByConditionAndRelation(ctx context.Context, condition string, args []interface{}, preload ...string) ([]entity.Competition, error) {
	var competitions []entity.Competition

	preloadConn := r.conn

	for _, relation := range preload {
		if len(relation) < 1 {
			continue
		}
		preloadConn = preloadConn.Preload(relation)
	}

	preloadConn = preloadConn.Order("id ASC")

	var res *gorm.DB
	if len(args) < 1 {
		res = preloadConn.Find(&competitions)
	} else {
		res = preloadConn.Where(condition, args).Find(&competitions)
	}

	if res.Error != nil {
		log.Error(log.LogInfo{
			"error": res.Error.Error(),
		}, "[COMPETITION REPOSITORY][FetchAllByConditionAndRelation] failed to fetch competition with relation and condition")

		return nil, res.Error
	}

	if res.RowsAffected < 1 {
		log.Error(log.LogInfo{
			"error": gorm.ErrRecordNotFound.Error(),
		}, "[COMPETITION REPOSITORY][FetchAllByConditionAndRelation] failed to fetch competition with relation and condition")

		return nil, domain.ErrNotFound
	}

	return competitions, nil
}

func (competitionRepo *competitionRepository) FetchOneByID(ctx context.Context, id int) (entity.Competition, error) {
	var competition entity.Competition

	err := competitionRepo.conn.First(&competition, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Competition{}, domain.ErrNotFound
		}

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[COMPETITION REPOSITORY][FetchOneByID] failed to fetch competition")

		return entity.Competition{}, err
	}

	return competition, nil
}

func (r *competitionRepository) FetchOneWithRelations(ctx context.Context, id int, relations ...string) (entity.Competition, error) {
	preloadConn := r.conn

	for _, relation := range relations {
		if len(relation) < 1 {
			continue
		}
		preloadConn = preloadConn.Preload(relation)
	}

	compe := entity.Competition{}

	err := preloadConn.First(&compe, id).Error

	if err != nil {

		if err == gorm.ErrRecordNotFound {
			return entity.Competition{}, domain.ErrNotFound
		}

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[COMPETITION REPOSITORY][FetchOneWithRelations] failed to fetch competition with relations")

		return entity.Competition{}, err
	}

	return compe, nil
}

func (r *competitionRepository) InsertCompe(ctx context.Context, compe *entity.Competition) error {
	err := r.conn.Create(compe).Error

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[COMPETITION REPOSITORY][InsertCompe] failed to insert competition to database")

		return err
	}

	return nil
}

func (r *competitionRepository) UpdateCompe(ctx context.Context, compe *entity.Competition) error {
	res := r.conn.Updates(compe)

	if err := helpers.CheckRowsAffected(res.RowsAffected); err != nil {

		if err != domain.ErrNotFound {
			log.Error(log.LogInfo{
				"error": err.Error(),
			}, "[COMPETITION REPOSITORY][UpdateCompe] failed to update competition data")
		}

		return err
	}

	if err := res.Error; err != nil {
		if err != domain.ErrNotFound {
			log.Error(log.LogInfo{
				"error": err.Error(),
			}, "[COMPETITION REPOSITORY][UpdateCompe] failed to update competition data")
		}

		return err
	}

	return nil
}

func (r *competitionRepository) DeleteCompe(ctx context.Context, id int) error {
	res := r.conn.Delete(&entity.Competition{ID: id})

	if err := helpers.CheckRowsAffected(res.RowsAffected); err != nil {

		if err != domain.ErrNotFound {
			log.Error(log.LogInfo{
				"error": err.Error(),
			}, "[COMPETITION REPOSITORY][DeleteCompe] failed to delete competition")
		}

		return err
	}

	if err := res.Error; err != nil {
		if err != domain.ErrNotFound {
			log.Error(log.LogInfo{
				"error": err.Error(),
			}, "[COMPETITION REPOSITORY][DeleteCompe] failed to delete competition")
		}

		return err
	}

	return nil
}
