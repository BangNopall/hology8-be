package repository

import (
	"context"
	"errors"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/contracts"
	"github.com/BangNopall/hology8-be/domain/entity"
	"github.com/BangNopall/hology8-be/pkg/helpers"
	"github.com/BangNopall/hology8-be/pkg/log"
	"gorm.io/gorm"
)

type voucherRepository struct {
	conn *gorm.DB
}

func NewVoucherRepository(conn *gorm.DB) contracts.VoucherRepository {
	return &voucherRepository{conn: conn}
}

func (r *voucherRepository) FetchAll(ctx context.Context) ([]entity.Voucher, error) {
	vouchers := make([]entity.Voucher, 0)

	err := r.conn.Find(&vouchers).Error

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[VOUCHER REPOSITORY][FetchAll] failed to fetch vouchers data")

		return nil, err
	}

	return vouchers, nil
}

func (r *voucherRepository) FetchByID(ctx context.Context, id string) (entity.Voucher, error) {
	var voucher entity.Voucher

	err := r.conn.Where("\"vouchers\".\"id\" = ?", id).First(&voucher).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Voucher{}, domain.ErrNotFound
		}

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[VOUCHER REPOSITORY][FetchByID] failed to fetch voucher data")

		return entity.Voucher{}, err
	}

	return voucher, nil
}

func (r *voucherRepository) InsertVoucher(ctx context.Context, voucher *entity.Voucher) error {
	err := r.conn.Create(voucher).Error

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[VOUCHER REPOSITORY][InsertVoucher] failed to insert voucher to database")

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateEntry
		}

		return err
	}

	return nil
}

func (r *voucherRepository) UpdateVoucher(ctx context.Context, voucher *entity.Voucher, tx *gorm.DB) error {
	conn := r.conn

	if tx != nil {
		conn = tx
	}

	res := conn.Updates(voucher)

	if err := helpers.CheckRowsAffected(res.RowsAffected); err != nil {
		if err != domain.ErrNotFound {
			log.Error(log.LogInfo{
				"error": err.Error(),
			}, "[VOUCHER REPOSITORY][UpdateVoucher] failed to update voucher data")
		}

		return err
	}

	if err := res.Error; err != nil {
		if err != domain.ErrNotFound {
			log.Error(log.LogInfo{
				"error": err.Error(),
			}, "[VOUCHER REPOSITORY][UpdateVoucher] failed to update voucher data")
		}

		return err
	}

	return nil
}
