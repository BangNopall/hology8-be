package repository

import (
	"math"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/pkg/log"
)

type userRepository struct {
	conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) contracts.UserRepository {
	return &userRepository{conn}
}

func (r *userRepository) FindUser(user *entity.User, userParam *dto.UserParam, relations ...string) error {
	preloadConn := r.conn

	for _, relation := range relations {
		preloadConn = preloadConn.Preload(relation)
	}

	err := preloadConn.First(user, userParam).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.ErrNotFound
		}

		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[USER REPOSITORY][FindUser] failed to find user")
		return domain.ErrInternalServer
	}

	return nil
}

func (r *userRepository) FetchAllByConditionAndRelation(
	condition string,
	args []interface{},
	joins []string,
	pageParam *dto.PaginationRequest,
	preload ...string,
) ([]entity.User, dto.PaginationResponse, error) {
	var users []entity.User
	var pageResp dto.PaginationResponse

	preloadConn := r.conn

	for _, join := range joins {
		if len(join) < 1 {
			continue
		}

		preloadConn = preloadConn.Joins(join)
	}

	for _, relation := range preload {
		if len(relation) < 1 {
			continue
		}
		preloadConn = preloadConn.Preload(relation)
	}

	if pageParam != nil {
		preloadConn = preloadConn.Offset(pageParam.Offset).Limit(pageParam.Limit)

		var count int64

		r.conn.Model(entity.User{}).Count(&count)

		pageResp.TotalPages = int(math.Ceil(float64(count) / float64(pageParam.Limit)))
		pageResp.Page = pageParam.Page
	}

	var res *gorm.DB
	if len(args) < 1 {
		res = preloadConn.Find(&users)
	} else {
		res = preloadConn.Where(condition, args).Find(&users)
	}

	if res.Error != nil {
		log.Error(log.LogInfo{
			"error": res.Error.Error(),
		}, "[USER REPOSITORY][FetchAllByConditionAndRelation] failed to fetch users with relation and condition")

		return nil, pageResp, res.Error
	}

	if res.RowsAffected < 1 {
		log.Error(log.LogInfo{
			"error": gorm.ErrRecordNotFound.Error(),
		}, "[USER REPOSITORY][FetchAllByConditionAndRelation] failed to fetch users with relation and condition")

		return nil, pageResp, domain.ErrNotFound
	}

	return users, pageResp, nil
}

func (r *userRepository) CreateUser(user *entity.User) error {
	err := r.conn.Create(user).Error

	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			return domain.ErrDuplicateEntry
		}

		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[USER REPOSITORY][Registration] failed to register user")

		return domain.ErrInternalServer
	}

	return nil
}

func (r *userRepository) UpdateUser(updateUser *dto.UserUpdate, userId uuid.UUID) error {
	err := r.conn.Model(&entity.User{}).Where("id = ?", userId).Updates(updateUser).Error
	if err != nil {

		if err == gorm.ErrDuplicatedKey {
			return domain.ErrDuplicateEntry
		}

		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[USER REPOSITORY][UpdateUser] failed to update user")

		return domain.ErrInternalServer
	}

	return nil
}

func (r *userRepository) DeleteUnverifiedUser() error {
	err := r.conn.Where("email_is_verified = ?", false).Delete(&entity.User{}).Error

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[USER REPOSITORY][DeleteUnverifiedUser] failed to delete unverified user")

		return domain.ErrInternalServer
	}

	return nil
}
