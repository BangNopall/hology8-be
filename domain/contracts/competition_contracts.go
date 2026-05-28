package contracts

import (
	"context"

	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/domain/entity"
)

type CompetitionRepository interface {
	FetchAll(ctx context.Context) ([]entity.Competition, error)
	FetchAllByConditionAndRelation(ctx context.Context, condition string, args []interface{}, preload ...string) ([]entity.Competition, error)
	FetchOneByID(ctx context.Context, id int) (entity.Competition, error)
	FetchOneWithRelations(ctx context.Context, id int, relations ...string) (entity.Competition, error)
	InsertCompe(ctx context.Context, compe *entity.Competition) error
	UpdateCompe(ctx context.Context, compe *entity.Competition) error
	DeleteCompe(ctx context.Context, id int) error
}

type CompetitionService interface {
	FetchAll(ctx context.Context, relations string) ([]dto.CompetitionResponse, error)
	FetchOne(ctx context.Context, id int, relations string) (dto.CompetitionResponse, error)
	InsertCompe(ctx context.Context, compe *dto.CompetitionRequest) error
	UpdateCompe(ctx context.Context, compe *dto.CompetitionRequest) error
	DeleteCompe(ctx context.Context, id int) error
}
