package dto

import (
	"github.com/hology8/hology-be/domain/entity"
)

type VoucherResponse struct {
	ID     string `json:"id"`
	TeamID string `json:"team_id"`
}

type VoucherRequest struct {
	ID string `json:"id"`
}

type VoucherRedeem struct {
	ID     string `json:"id"`
	TeamID string `json:"team_id"`
}

func VoucherEntityToResponse(voucher *entity.Voucher) VoucherResponse {
	return VoucherResponse{
		ID:     voucher.ID,
		TeamID: voucher.TeamID,
	}
}
