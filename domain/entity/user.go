package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID     `json:"id" gorm:"primaryKey;type:uuid"`
	Email               string        `json:"email" gorm:"uniqueIndex;type:varchar(200)"`
	Password            string        `json:"password" gorm:"type:varchar(255)"`
	Fullname            string        `json:"fullname" gorm:"type:varchar(200)"`
	BirthDate           time.Time     `json:"birth_date" gorm:"type:date"`
	WANumber            string        `json:"wa_number" gorm:"uniqueIndex;type:varchar(30);default:null"`
	LineID              string        `json:"id_line" gorm:"type:varchar(50);default:null"`
	DiscordID           string        `json:"id_discord" gorm:"type:varchar(50);default:null"`
	StudentID           string        `json:"id_student" gorm:"type:varchar(50);default:null"`
	EmailVerifiedToken  string        `json:"email_verified_token" gorm:"type:varchar(100)"`
	ForgotPasswordToken string        `json:"forgot_password_token" gorm:"type:varchar(100)"`
	EmailIsVerified     bool          `json:"email_is_verified" gorm:"type:bool"`
	KtmImageLink        string        `json:"ktm_image_link" gorm:"type:varchar(255)"`
	FollowProofLink     string        `json:"follow_proof_link" gorm:"type:varchar(255)"`
	ShareProofLink      string        `json:"share_proof_link" gorm:"type:varchar(255)"`
	City                string        `json:"city" gorm:"type:varchar(100)"`
	HomeAddress         string        `json:"home_address" gorm:"type:varchar(100)"`
	ProvinceID          int           `json:"-" gorm:"type:integer;default:null"`
	UniversityID        int           `json:"-" gorm:"type:integer;default:null"`
	ExpiredToken        time.Time     `json:"-" gorm:"type:timestamp"`
	ExpiredTokenForgot  time.Time     `json:"-" gorm:"type:timestamp"`
	University          University    `json:"university" gorm:"foreignKey:university_id"`
	Province            Province      `json:"province" gorm:"foreignKey:province_id"`
	Teams               []DetailTeams `json:"detail_teams" gorm:"foreignKey:user_id"`
}
