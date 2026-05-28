package contracts

import (
	"context"
	"mime/multipart"

	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/domain/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(user *entity.User) error
	FetchAllByConditionAndRelation(
		condition string,
		args []interface{},
		joins []string,
		pageParam *dto.PaginationRequest,
		preload ...string,
	) ([]entity.User, dto.PaginationResponse, error)
	FindUser(user *entity.User, userParam *dto.UserParam, relations ...string) error
	UpdateUser(updateUser *dto.UserUpdate, userId uuid.UUID) error
	DeleteUnverifiedUser() error
}

type UserService interface {
	LoginRegisterOauth(ctx context.Context, user dto.UserOauth) (dto.UserLoginResponse, error)
	Register(ctx context.Context, user dto.UserRegister, referer string) error
	VerifyEmail(ctx context.Context, email string, emailVerPass string) error
	LoginWithEmail(ctx context.Context, user dto.UserLogin) (dto.UserLoginResponse, error)
	ResetPassword(ctx context.Context, user dto.UserResetPassword, forgotPasswordToken string) error
	ForgotPassword(ctx context.Context, user dto.UserForgotPassword, referer string) error
	UpdateUser(ctx context.Context, userID uuid.UUID, userUpdate dto.UserUpdate) error
	UploadKtmImage(ctx context.Context, userId uuid.UUID, ktmImage *multipart.FileHeader) error
	UploadProofImage(ctx context.Context, userId uuid.UUID, proofName string, file *multipart.FileHeader) error
	FetchByParam(ctx context.Context, userParam *dto.UserParam) (dto.UserResponse, error)
	FetchAll(ctx context.Context, userParam *dto.UserParam, pageParam *dto.PaginationRequest) ([]dto.UserResponse, dto.PaginationResponse, error)
	Logout(ctx context.Context, jwtToken string) error
	DeleteUnverifiedUsers()
}
