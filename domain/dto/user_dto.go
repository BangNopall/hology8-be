package dto

import (
	"time"

	"github.com/BangNopall/hology8-be/domain/entity"
	"github.com/google/uuid"
)

type UserUpdate struct {
	Password            string    `json:"password" binding:"omitempty,validpassword"`
	Fullname            string    `json:"fullname" binding:"omitempty,alphnumsympace,max=100"`
	BirthDate           string    `json:"birth_date" binding:"omitempty,validdate"`
	WANumber            string    `json:"wa_number" binding:"omitempty,plusnumeric"`
	LineID              string    `json:"id_line" binding:"omitempty,max=70"`
	DiscordID           string    `json:"id_discord" binding:"omitempty,max=70"`
	StudentID           string    `json:"id_student" binding:"omitempty,max=70"`
	City                string    `json:"city" binding:"omitempty,max=70"`
	UniversityID        int       `json:"university_id" binding:"omitempty,numeric"`
	ProvinceID          int       `json:"province_id" binding:"omitempty,numeric"`
	KtmImageLink        string    `json:"-"`
	FollowProofLink     string    `json:"-"`
	ShareProofLink      string    `json:"-"`
	EmailIsVerified     bool      `json:"-"`
	EmailVerifiedToken  string    `json:"-"`
	ForgotPasswordToken string    `json:"-"`
	ExpiredToken        time.Time `json:"-"`
	ExpiredTokenForgot  time.Time `json:"-"`
}

type UserRegister struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6,max=40,validpassword"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6,max=40,validpassword"`
}

type UserOauth struct {
	Email      string `json:"email" gorm:"primary key"`
	FamilyName string `json:"family_name"`
	GivenName  string `json:"given_name"`
	Id         string `json:"id"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
}

type UserResponse struct {
	ID                  uuid.UUID          `json:"id"`
	Email               string             `json:"email"`
	Password            string             `json:"password"`
	Fullname            string             `json:"fullname"`
	BirthDate           time.Time          `json:"birth_date"`
	WANumber            string             `json:"wa_number"`
	LineID              string             `json:"id_line"`
	DiscordID           string             `json:"id_discord"`
	StudentID           string             `json:"id_student"`
	EmailVerifiedToken  string             `json:"email_verified_token"`
	ForgotPasswordToken string             `json:"forgot_password_token"`
	EmailIsVerified     bool               `json:"email_is_verified"`
	KtmImageLink        string             `json:"ktm_image_link"`
	FollowProofLink     string             `json:"follow_proof_link"`
	ShareProofLink      string             `json:"share_proof_link"`
	City                string             `json:"city"`
	HomeAddress         string             `json:"home_address"`
	ProvinceID          int                `json:"province_id"`
	UniversityID        int                `json:"university_id"`
	ExpiredToken        time.Time          `json:"-"`
	ExpiredTokenForgot  time.Time          `json:"-"`
	Province            ProvinceResponse   `json:"province"`
	University          UniversityResponse `json:"university"`
}

type UserPaginationResponse struct {
	Users      []UserResponse     `json:"users"`
	Pagination PaginationResponse `json:"pagination"`
}

type UserParam struct {
	ID                  uuid.UUID `json:"id"`
	Email               string    `json:"email"`
	ForgotPasswordToken string    `json:"forgot_password_token"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserForgotPassword struct {
	Email string `json:"email" binding:"required,email"`
}

type UserResetPassword struct {
	Password        string `json:"password" binding:"required,min=6,max=20"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6,max=20"`
}

type UserOuathLink struct {
	RedirectLink string `json:"redirect_link"`
}

func UserEntityToResponseDto(user *entity.User) *UserResponse {
	return &UserResponse{
		ID:                  user.ID,
		Email:               user.Email,
		Password:            user.Password,
		Fullname:            user.Fullname,
		BirthDate:           user.BirthDate,
		WANumber:            user.WANumber,
		LineID:              user.LineID,
		DiscordID:           user.DiscordID,
		StudentID:           user.StudentID,
		EmailVerifiedToken:  user.EmailVerifiedToken,
		ForgotPasswordToken: user.ForgotPasswordToken,
		EmailIsVerified:     user.EmailIsVerified,
		KtmImageLink:        user.KtmImageLink,
		FollowProofLink:     user.FollowProofLink,
		ShareProofLink:      user.ShareProofLink,
		City:                user.City,
		HomeAddress:         user.HomeAddress,
		ProvinceID:          user.ProvinceID,
		UniversityID:        user.UniversityID,
		ExpiredToken:        user.ExpiredToken,
		ExpiredTokenForgot:  user.ExpiredTokenForgot,
	}
}

func ConvertUserEntityToResponseDto(user *entity.User) *UserResponse {
	userResp := UserEntityToResponseDto(user)

	userResp.Province = ProvinceEntityToDto(&user.Province)

	userResp.University = UniversityEntityToDto(&user.University)

	return userResp
}
