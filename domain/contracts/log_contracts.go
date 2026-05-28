package contracts

import (
	"context"

	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/domain/entity"
)

type LogRepository interface {
	InsertLog(ctx context.Context, log *entity.Log) error
	FetchOneByID(ctx context.Context, id int) (entity.Log, error)
	FetchAll(ctx context.Context, relations ...string) ([]entity.Log, error)
}

type LogService interface {
	InsertLog(ctx context.Context, log *dto.LogRequest) error
	FetchAllLogs(ctx context.Context) ([]dto.LogResponse, error)
}
