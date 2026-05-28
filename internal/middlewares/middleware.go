package middlewares

import (
	"github.com/BangNopall/hology8-be/domain/contracts"
	"github.com/BangNopall/hology8-be/pkg/jwt"
	"github.com/BangNopall/hology8-be/pkg/redis"
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
