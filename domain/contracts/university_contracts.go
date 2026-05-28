package contracts

import (
	"context"

	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/domain/entity"
)

type UniversityRepository interface {
	FetchAll(ctx context.Context, params *dto.UniversityParam) ([]entity.University, error)
	FetchByID(ctx context.Context, id int) (entity.University, error)
}

type UniversityService interface {
	FetchAll(ctx context.Context, params *dto.UniversityParam) ([]dto.UniversityResponse, error)
	FetchByID(ctx context.Context, id int) (dto.UniversityResponse, error)
}
