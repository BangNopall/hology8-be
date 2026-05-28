package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/pkg/log"
)

type presenceRepository struct {
	conn *gorm.DB
}

func NewPresenceRepository(conn *gorm.DB) contracts.PresenceRepository {
	return &presenceRepository{conn}
}

func (presenceRepo *presenceRepository) InsertPresence(ctx context.Context, presence *entity.Presence) error {
	res := presenceRepo.conn.WithContext(ctx).Create(presence)
	if res.Error != nil {
		// if user not found
		if errors.Is(res.Error, gorm.ErrForeignKeyViolated) {
			return domain.ErrUserNotFound
		}

		// if presence already exist
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			return domain.ErrUserPresenceAlreadyExist
		}

		log.Error(log.LogInfo{
			"error": res.Error.Error(),
		}, "[PRESENCE REPOSITORY][InsertPresence] failed to create presence")

		return res.Error
	}

	return nil
}

func (presenceRepo *presenceRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.Presence, error) {
	var presence entity.Presence
	res := presenceRepo.conn.WithContext(ctx).
		Preload("User").
		Preload("User.Teams.Team").
		First(&presence, "user_id = ?", userID)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error(log.LogInfo{
			"error": res.Error.Error(),
		}, "[PRESENCE REPOSITORY][FindByUserID] failed to find presence")
		return nil, res.Error
	}

	return &presence, nil
}

func (presenceRepo *presenceRepository) FetchAll(ctx context.Context, offset, limit int) ([]entity.Presence, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	var total int64
	if err := presenceRepo.conn.WithContext(ctx).
		Model(&entity.Presence{}).
		Count(&total).Error; err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[PRESENCE REPOSITORY][FetchAll] failed to count")
		return nil, 0, err
	}

	var presences []entity.Presence
	if err := presenceRepo.conn.WithContext(ctx).
		Preload("User").
		Preload("User.Teams.Team").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&presences).Error; err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[PRESENCE REPOSITORY][FetchAll] failed to fetch list")
		return nil, 0, err
	}

	return presences, total, nil
}
