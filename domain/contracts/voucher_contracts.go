package contracts

import (
	"context"

	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/domain/entity"
	"gorm.io/gorm"
)

type VoucherRepository interface {
	FetchAll(ctx context.Context) ([]entity.Voucher, error)
	FetchByID(ctx context.Context, id string) (entity.Voucher, error)
	InsertVoucher(ctx context.Context, voucher *entity.Voucher) error
	UpdateVoucher(ctx context.Context, voucher *entity.Voucher, tx *gorm.DB) error
}

type VoucherService interface {
	FetchAll(ctx context.Context) ([]dto.VoucherResponse, error)
	FetchByID(ctx context.Context, id string) (dto.VoucherResponse, error)
	InsertVoucher(ctx context.Context, voucher *dto.VoucherRequest) error
	RedeemVoucher(ctx context.Context, voucher *dto.VoucherRedeem) error
}
