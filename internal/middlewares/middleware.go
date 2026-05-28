package middlewares

import (
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/pkg/jwt"
	"github.com/hology8/hology-be/pkg/redis"
)

type Middleware struct {
	jwt       jwt.JwtInterface
	adminRepo contracts.AdminRepository
	teamRepo  contracts.TeamRepository
	userRepo  contracts.UserRepository
	redis     redis.RedisInterface
}

func NewMiddleware(
	jwt jwt.JwtInterface,
	adminRepo contracts.AdminRepository,
	teamRepo contracts.TeamRepository,
	userRepo contracts.UserRepository,
	redis redis.RedisInterface,
) *Middleware {
	return &Middleware{
		jwt:       jwt,
		adminRepo: adminRepo,
		teamRepo:  teamRepo,
		userRepo:  userRepo,
		redis:     redis,
	}
}
