package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/pkg/log"
)

type universityRepository struct {
	conn *gorm.DB
}

func NewUniversityRepository(conn *gorm.DB) contracts.UniversityRepository {
	return &universityRepository{conn}
}

func (r *universityRepository) FetchAll(ctx context.Context, params *dto.UniversityParam) ([]entity.University, error) {
	universities := make([]entity.University, 0)

	conn := r.conn

	if params != nil {
		conn = r.conn.Where("university_name LIKE ?", "%"+params.Name+"%")
	}

	err := conn.Find(&universities).Error

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[UNIVERSITY REPOSITORY][FetchAll] failed to fetch universities data")

		return nil, err
	}

	return universities, nil
}

func (r *universityRepository) FetchByID(ctx context.Context, id int) (entity.University, error) {
	var university entity.University

	err := r.conn.First(&university, id).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.University{}, domain.ErrNotFound
		}

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[UNIVERSITY REPOSITORY][FetchByID] failed to fetch university data")

		return entity.University{}, err
	}

	return university, nil
}
