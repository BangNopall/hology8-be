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

type logRepository struct {
	conn *gorm.DB
}

func NewLogRepository(conn *gorm.DB) contracts.LogRepository {
	return &logRepository{conn}
}

func (logRepo *logRepository) InsertLog(ctx context.Context, l *entity.Log) error {
	res := logRepo.conn.Create(l)

	if res.Error != nil {
		log.Error(log.LogInfo{
			"error": res.Error.Error(),
		}, "[LOG REPOSITORY][InsertLog] failed to create log")

		return res.Error
	}

	return nil
}

func (logRepo *logRepository) FetchOneByID(ctx context.Context, id int) (entity.Log, error) {
	var l entity.Log
	err := logRepo.conn.First(&l, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Log{}, domain.ErrNotFound
		}

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[LOG REPOSITORY][FetchOneByID] failed to fetch log")

		return entity.Log{}, err
	}

	return l, nil
}

func (logRepo *logRepository) FetchAll(ctx context.Context, relations ...string) ([]entity.Log, error) {
	preloadConn := logRepo.conn

	for _, relation := range relations {
		preloadConn = preloadConn.Preload(relation)
	}

	logs := make([]entity.Log, 0)

	err := preloadConn.Order("created_at DESC").Find(&logs).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[LOG REPOSITORY][FetchAll] failed to fetch logs")

		return nil, domain.ErrInternalServer
	}

	return logs, nil
}
