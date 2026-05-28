package repository

import (
	"context"
	"errors"

	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/pkg/log"
	"gorm.io/gorm"
)

type partnerRepository struct {
	conn *gorm.DB
}

func NewPartnerRepository(conn *gorm.DB) contracts.PartnerRepository {
	return &partnerRepository{conn}
}

func (p *partnerRepository) FetchAll(ctx context.Context, params *dto.PartnerParams) ([]entity.Partner, error) {
	partners := make([]entity.Partner, 0)

	err := p.conn.Find(&partners, params).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[PARTNER REPOSITORY][FetchAll] failed to fetch partners data")

		return nil, err
	}

	return partners, nil
}

func (p *partnerRepository) FetchAllTypes(ctx context.Context) ([]entity.PartnerType, error) {
	partnerTypes := make([]entity.PartnerType, 0)

	err := p.conn.Find(&partnerTypes).Error

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[PARTNER REPOSITORY][FetchAllTypes] failed to fetch partner types data")

		return nil, err
	}

	return partnerTypes, nil
}

func (p *partnerRepository) FetchOneByID(ctx context.Context, id int) (entity.Partner, error) {
	var partner entity.Partner

	err := p.conn.First(&partner, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Partner{}, domain.ErrNotFound
		}

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[PARTNER REPOSITORY][FetchOneByID] failed to fetch partner data")

		return entity.Partner{}, err
	}

	return partner, nil
}

func (p *partnerRepository) InsertPartner(ctx context.Context, partner *entity.Partner) error {
	err := p.conn.Create(partner).Error

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[PARTNER REPOSITORY][InsertPartner] failed to insert partner data")

		return err
	}

	return nil
}

func (p *partnerRepository) UpdatePartner(ctx context.Context, id int, partner *entity.Partner) error {
	err := p.conn.Model(&entity.Partner{}).Where("id = ?", id).Updates(partner).Error

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[PARTNER REPOSITORY][UpdatePartner] failed to update partner data")

		return err
	}

	return nil
}

func (p *partnerRepository) DeletePartner(ctx context.Context, id int) error {
	err := p.conn.Delete(&entity.Partner{}, id).Error

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[PARTNER REPOSITORY][DeletePartner] failed to delete partner data")

		return err
	}

	return nil
}
