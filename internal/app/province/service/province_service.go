package service

import (
	"context"
	"time"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/contracts"
	"github.com/BangNopall/hology8-be/domain/dto"
)

type provinceService struct {
	provinceRepo contracts.ProvinceRepository
	timeout      time.Duration
}

func NewProvinceService(provinceRepo contracts.ProvinceRepository, timeout time.Duration) contracts.ProvinceService {
	return &provinceService{provinceRepo, timeout}
}

func (s *provinceService) FetchAll(ctx context.Context) ([]dto.ProvinceResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	provinces, err := s.provinceRepo.FetchAll(ctx)

	if err != nil {
		return nil, err
	}

	res := make([]dto.ProvinceResponse, 0)

	for _, prv := range provinces {
		res = append(res, dto.ProvinceEntityToDto(&prv))
	}

	select {
	case <-ctx.Done():
		return nil, domain.ErrTimeout
	default:
		return res, err
	}
}

func (s *provinceService) FetchByID(ctx context.Context, id int) (dto.ProvinceResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	province, err := s.provinceRepo.FetchByID(ctx, id)

	if err != nil {
		return dto.ProvinceResponse{}, err
	}

	res := dto.ProvinceEntityToDto(&province)

	select {
	case <-ctx.Done():
		return dto.ProvinceResponse{}, domain.ErrTimeout
	default:
		return res, err
	}
}
