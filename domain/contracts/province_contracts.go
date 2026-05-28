package contracts

import (
	"context"

	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
)

type ProvinceRepository interface {
	FetchAll(ctx context.Context) ([]entity.Province, error)
	FetchByID(ctx context.Context, id int) (entity.Province, error)
}

type ProvinceService interface {
	FetchAll(ctx context.Context) ([]dto.ProvinceResponse, error)
	FetchByID(ctx context.Context, id int) (dto.ProvinceResponse, error)
}
