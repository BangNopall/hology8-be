package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/pkg/log"
)

type provinceRepository struct {
	conn *gorm.DB
}

func NewProvinceRepository(conn *gorm.DB) contracts.ProvinceRepository {
	return &provinceRepository{conn}
}

func (r *provinceRepository) FetchAll(ctx context.Context) ([]entity.Province, error) {
	provinces := make([]entity.Province, 0)

	err := r.conn.Find(&provinces).Error

	if err != nil {

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[PROVINCE REPOSITORY][FetchAll] failed to fetch provinces data")

		return nil, err
	}

	return provinces, nil
}

func (r *provinceRepository) FetchByID(ctx context.Context, id int) (entity.Province, error) {
	var province entity.Province

	err := r.conn.First(&province, id).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Province{}, domain.ErrNotFound
		}

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[PROVINCE REPOSITORY][FetchByID] failed to fetch province data")

		return entity.Province{}, err
	}

	return province, nil
}
