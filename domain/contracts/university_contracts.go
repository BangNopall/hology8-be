package contracts

import (
	"context"

	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
)

type UniversityRepository interface {
	FetchAll(ctx context.Context, params *dto.UniversityParam) ([]entity.University, error)
	FetchByID(ctx context.Context, id int) (entity.University, error)
}

type UniversityService interface {
	FetchAll(ctx context.Context, params *dto.UniversityParam) ([]dto.UniversityResponse, error)
	FetchByID(ctx context.Context, id int) (dto.UniversityResponse, error)
}
