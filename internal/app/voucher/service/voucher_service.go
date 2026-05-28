package service

import (
	"context"
	"time"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/contracts"
	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/domain/entity"
	"gorm.io/gorm"
)

type voucherService struct {
	voucherRepo contracts.VoucherRepository
	teamRepo    contracts.TeamRepository
	timeout     time.Duration
	db          *gorm.DB
}

func NewVoucherService(
	voucherRepo contracts.VoucherRepository,
	teamRepo contracts.TeamRepository,
	timeout time.Duration,
	db *gorm.DB,
) contracts.VoucherService {
	return &voucherService{voucherRepo, teamRepo, timeout, db}
}

func (s *voucherService) FetchAll(ctx context.Context) ([]dto.VoucherResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	vouchers, err := s.voucherRepo.FetchAll(ctx)

	if err != nil {
		return nil, err
	}

	res := make([]dto.VoucherResponse, 0)

	for _, vch := range vouchers {
		res = append(res, dto.VoucherEntityToResponse(&vch))
	}

	select {
	case <-ctx.Done():
		return nil, domain.ErrTimeout
	default:
		return res, err
	}
}

func (s *voucherService) FetchByID(ctx context.Context, id string) (dto.VoucherResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	voucher, err := s.voucherRepo.FetchByID(ctx, id)

	if err != nil {
		return dto.VoucherResponse{}, err
	}

	res := dto.VoucherEntityToResponse(&voucher)

	select {
	case <-ctx.Done():
		return dto.VoucherResponse{}, domain.ErrTimeout
	default:
		return res, err
	}
}

func (s *voucherService) InsertVoucher(ctx context.Context, voucher *dto.VoucherRequest) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	payload := entity.NewVoucher(voucher.ID)

	err := s.voucherRepo.InsertVoucher(ctx, payload)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (s *voucherService) RedeemVoucher(ctx context.Context, voucherRedeem *dto.VoucherRedeem) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	tx := s.db.Begin()

	team, err := s.teamRepo.FetchOneByID(ctx, voucherRedeem.TeamID)
	if err != nil {
		tx.Rollback()
		return domain.ErrTeamNotFound
	}

	voucher, err := s.voucherRepo.FetchByID(ctx, voucherRedeem.ID)
	if err != nil {
		tx.Rollback()
		return domain.ErrVoucherNotFound
	}

	if team.VoucherID != nil || voucher.TeamID != "" {
		tx.Rollback()
		return domain.ErrVoucherAlreadyRedeemed
	}

	payload := entity.NewVoucher(voucherRedeem.ID)
	payload.TeamID = team.ID
	payload.Team = &team

	err = s.voucherRepo.UpdateVoucher(ctx, payload, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = s.teamRepo.LinkVoucher(ctx, &team, payload, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}
