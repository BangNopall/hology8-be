package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
)

type logService struct {
	logRepo contracts.LogRepository
	timeout time.Duration
}

func NewLogService(logRepo contracts.LogRepository, timeout time.Duration) contracts.LogService {
	return &logService{logRepo, timeout}
}

func (logSvc *logService) InsertLog(ctx context.Context, log *dto.LogRequest) error {
	ctx, cancel := context.WithTimeout(ctx, logSvc.timeout)
	defer cancel()

	logE := &entity.Log{
		AdminID: uuid.MustParse(log.AdminID),
		Action:  log.Action,
	}

	err := logSvc.logRepo.InsertLog(ctx, logE)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (logSvc *logService) FetchAllLogs(ctx context.Context) ([]dto.LogResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, logSvc.timeout)
	defer cancel()

	logs, err := logSvc.logRepo.FetchAll(ctx, "Admin")

	var logResponses []dto.LogResponse

	for _, log := range logs {
		logResponses = append(logResponses, dto.LogResponse{
			ID:        log.ID,
			Fullname:  log.Admin.Fullname,
			Action:    log.Action,
			CreatedAt: log.CreatedAt,
		})
	}

	select {
	case <-ctx.Done():
		return nil, domain.ErrTimeout
	default:
		return logResponses, err
	}
}
