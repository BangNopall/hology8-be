package domain

import (
	"errors"
	"net/http"
)

var (
	ErrAnnouncementNotFound     = errors.New("announcement not found")
	ErrCompetitionNotFound      = errors.New("competition not found")
	ErrTeamNotFound             = errors.New("team not found")
	ErrUserNotTeamLeader        = errors.New("user not a team leader")
	ErrVoucherNotFound          = errors.New("invalid voucher")
	ErrVoucherAlreadyRedeemed   = errors.New("voucher already redeemed")
	ErrInvalidCompeTeamID       = errors.New("invalid compe or team id")
	ErrIllegalEntry             = errors.New("illegal entry")
	ErrUserAlreadyRegistered    = errors.New("user already registered")
	ErrInternalServer           = errors.New("internal server error")
	ErrEmailRegistered          = errors.New("email already registered")
	ErrCheckEmail               = errors.New("please check your email to verify")
	ErrInvalidToken             = errors.New("invalid token")
	ErrWrongEmailOrPassword     = errors.New("invalid username/password")
	ErrFileTooBig               = errors.New("file size too big")
	ErrNotFound                 = errors.New("item not found")
	ErrTimeout                  = errors.New("operation timed out")
	ErrDuplicateEntry           = errors.New("data already exists")
	ErrUserNotFound             = errors.New("user not found")
	ErrUniversityNotFound       = errors.New("university not found")
	ErrAlreadyAttend            = errors.New("user already attended")
	ErrForbiddenUpdate          = errors.New("forbidden to update status or phase in this service")
	ErrInvalidEnumInput         = errors.New("invalid enum input")
	ErrMissingAttribute         = errors.New("missing required attribute")
	ErrConfirmPasswordNotMatch  = errors.New("password and confirm password doesn't match")
	ErrTeamFull                 = errors.New("team is full")
	ErrInvalidProofType         = errors.New("invalid proof type")
	ErrBadRequest               = errors.New("bad data request")
	ErrUserPresenceAlreadyExist = errors.New("user presence already exist")
)

func GetCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case ErrInternalServer:
		return http.StatusInternalServerError
	case ErrNotFound,
		ErrVoucherNotFound,
		ErrUserNotFound:
		return http.StatusNotFound
	case ErrInvalidEnumInput:
		return http.StatusBadRequest
	case ErrUserNotTeamLeader,
		ErrAnnouncementNotFound,
		ErrCompetitionNotFound,
		ErrTeamNotFound,
		ErrIllegalEntry,
		ErrMissingAttribute,
		ErrTeamFull,
		ErrInvalidCompeTeamID,
		ErrInvalidProofType,
		ErrUniversityNotFound,
		ErrBadRequest:
		return http.StatusBadRequest
	case ErrTimeout:
		return http.StatusRequestTimeout
	case ErrDuplicateEntry, ErrUserAlreadyRegistered, ErrEmailRegistered, ErrVoucherAlreadyRedeemed, ErrUserPresenceAlreadyExist:
		return http.StatusConflict
	case ErrAlreadyAttend:
		return http.StatusUnprocessableEntity
	case ErrCheckEmail:
		return http.StatusBadRequest
	case ErrInvalidToken:
		return http.StatusUnauthorized
	case ErrWrongEmailOrPassword:
		return http.StatusUnauthorized
	case ErrFileTooBig:
		return http.StatusRequestEntityTooLarge
	case ErrForbiddenUpdate:
		return http.StatusForbidden
	case ErrConfirmPasswordNotMatch:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
