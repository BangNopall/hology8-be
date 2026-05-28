package repository

import (
	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/pkg/log"
	"gorm.io/gorm"
)

type adminRepository struct {
	conn *gorm.DB
}

func NewAdminRepository(conn *gorm.DB) contracts.AdminRepository {
	return &adminRepository{conn}
}

func (r *adminRepository) FindAdmin(admin *entity.Admin, adminParam *dto.AdminParam) error {
	err := r.conn.Preload("Role").First(admin, adminParam).Error
	if err != nil {

		if err == gorm.ErrRecordNotFound {
			return domain.ErrNotFound
		}

		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[ADMIN REPOSITORY][FindByUsername] failed to find admin")

		return err
	}

	return nil
}
