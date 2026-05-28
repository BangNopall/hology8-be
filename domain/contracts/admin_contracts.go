package contracts

import (
	"context"

	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/domain/entity"
)

type AdminRepository interface {
	FindAdmin(admin *entity.Admin, adminParam *dto.AdminParam) error
}

type AdminService interface {
	Login(ctx context.Context, adminLogin dto.AdminLogin) (dto.AdminLoginResponse, error)
	SendEmail(ctx context.Context, to string, emailMessage dto.EmailMessage) error
}
