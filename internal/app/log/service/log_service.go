package service

import (
	"context"
	"time"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/contracts"
	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/domain/entity"
	"github.com/google/uuid"
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
