package service

import (
	"context"
	"time"

	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/dto"
)

type universityService struct {
	uniRepo contracts.UniversityRepository
	timeout time.Duration
}

func NewUniversityService(uniRepo contracts.UniversityRepository, timeout time.Duration) contracts.UniversityService {
	return &universityService{uniRepo, timeout}
}

func (s *universityService) FetchAll(ctx context.Context, params *dto.UniversityParam) ([]dto.UniversityResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	universities, err := s.uniRepo.FetchAll(ctx, params)

	if err != nil {
		return nil, err
	}

	res := make([]dto.UniversityResponse, 0)

	for _, uni := range universities {
		res = append(res, dto.UniversityEntityToDto(&uni))
	}

	select {
	case <-ctx.Done():
		return nil, domain.ErrTimeout
	default:
		return res, err
	}
}

func (s *universityService) FetchByID(ctx context.Context, id int) (dto.UniversityResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	university, err := s.uniRepo.FetchByID(ctx, id)

	if err != nil {
		return dto.UniversityResponse{}, err
	}

	res := dto.UniversityEntityToDto(&university)

	select {
	case <-ctx.Done():
		return dto.UniversityResponse{}, domain.ErrTimeout
	default:
		return res, err
	}
}
