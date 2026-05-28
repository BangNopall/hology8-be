package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/internal/infra/env"
	"github.com/hology8/hology-be/pkg/helpers/http/response"
	"github.com/hology8/hology-be/pkg/log"
)

func (m *Middleware) Authentication(ctx *gin.Context) {
	bearer := ctx.GetHeader("Authorization")
	if bearer == "" {
		log.Warn(log.LogInfo{
			"error": errors.New("failed to get bearer token"),
		}, "[MIDDLEWARE][Authentication] failed to get bearer token")

		response.SendErrResp(
			ctx,
			http.StatusUnauthorized,
			response.Error,
			"failed to authenticate user",
			errors.New("failed to get bearer token"),
		)

		ctx.Abort()
		return
	}

	splitted := strings.Split(bearer, " ")

	if len(splitted) < 2 {
		response.SendErrResp(
			ctx,
			400,
			response.Fail,
			"failed to authenticate user",
			fmt.Errorf("invalid token"),
		)
		ctx.Abort()
		return
	}

	tokenString := splitted[1]

	id, user, role, err := m.jwt.ValidateToken(tokenString)
	if err != nil {
		log.Warn(log.LogInfo{
			"error": err,
		}, "[MIDDLEWARE][Authentication] failed to validate token")

		response.SendErrResp(
			ctx,
			http.StatusUnauthorized,
			response.Error,
			"failed to authenticate user",
			err,
		)

		ctx.Abort()
		return
	}

	if user != env.AppEnv.JwtUserRole && user != env.AppEnv.JwtAdminRole {
		response.SendErrResp(
			ctx,
			http.StatusUnauthorized,
			response.Fail,
			"failed to authenticate user",
			err,
		)
		ctx.Abort()
		return
	}

	val, err := m.redis.Get(ctx.Request.Context(), tokenString)

	if err != nil {
		response.SendErrResp(
			ctx,
			http.StatusInternalServerError,
			response.Error,
			"failed to authenticate user",
			err,
		)

		ctx.Abort()
		return
	}

	if val != "" {
		response.SendErrResp(
			ctx,
			http.StatusUnauthorized,
			response.Fail,
			"failed to authenticate user",
			nil,
		)

		ctx.Abort()
		return
	}

	ctx.Set("id", id.String())
	ctx.Set("user", user)
	ctx.Set("role", role)
	ctx.Next()
}

func (m *Middleware) AuthorizationAdmin(ctx *gin.Context) {
	user := ctx.GetString("user")
	id := ctx.GetString("id")

	if user != env.AppEnv.JwtAdminRole {
		response.SendErrResp(
			ctx,
			http.StatusUnauthorized,
			response.Fail,
			"unauthorized",
			errors.New("unauthorized"),
		)
		ctx.Abort()
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		log.Warn(log.LogInfo{
			"error": err,
		}, "[MIDDLEWARE][AuthorizationAdmin] failed to parse uuid")

		response.SendErrResp(
			ctx,
			http.StatusInternalServerError,
			response.Error,
			"failed to parse uuid",
			err,
		)

		ctx.Abort()
		return
	}

	err = m.adminRepo.FindAdmin(&entity.Admin{}, &dto.AdminParam{ID: uuid})
	if err != nil {
		if err == domain.ErrNotFound {
			response.SendErrResp(
				ctx,
				http.StatusUnauthorized,
				response.Fail,
				"unauthorized",
				err,
			)
			ctx.Abort()
			return
		}

		response.SendErrResp(
			ctx,
			http.StatusInternalServerError,
			response.Error,
			"failed to authorize admin",
			err,
		)

		ctx.Abort()
		return
	}

	ctx.Next()

}

func (m *Middleware) AuthorizationTeam(ctx *gin.Context) {
	teamId := ctx.Param("id")
	id := ctx.GetString("id")

	userId, err := uuid.Parse(id)

	if err != nil {
		log.Error(log.LogInfo{
			"error": err,
		}, "[MIDDLEWARE][AuthorizationTeam] failed to parse uuid")

		response.SendErrResp(
			ctx,
			http.StatusInternalServerError,
			response.Error,
			"failed to parse uuid",
			err,
		)

		ctx.Abort()
		return
	}

	var (
		errChans = make(chan error, 2)
		wg       sync.WaitGroup
	)

	wg.Add(2)

	go func(ctx context.Context, teamId string, userId uuid.UUID) {
		defer wg.Done()
		_, err = m.teamRepo.FetchTeamMember(ctx, teamId, userId)

		if err != nil {
			errChans <- err
		}
	}(ctx.Request.Context(), teamId, userId)

	go func(ctx context.Context, userId uuid.UUID) {
		defer wg.Done()
		_, err := m.teamRepo.FetchOneByParams(ctx, &dto.TeamParams{LeaderID: userId})

		if err != nil {
			errChans <- err
		}
	}(ctx.Request.Context(), userId)

	go func() {
		wg.Wait()
		close(errChans)
	}()

	counter := 0

	for err := range errChans {
		if err != nil {
			if err == domain.ErrNotFound {
				counter++
				continue
			}

			log.Error(log.LogInfo{
				"error": err,
			}, "[MIDDLEWARE][AuthorizationTeam] failed to authorized team")

			response.SendErrResp(
				ctx,
				http.StatusInternalServerError,
				response.Error,
				"failed to authorize team",
				err,
			)

			ctx.Abort()
			return
		}
	}

	if counter > 1 {
		response.SendErrResp(
			ctx,
			http.StatusUnauthorized,
			response.Fail,
			"unauthorized",
			err,
		)
		ctx.Abort()
		return
	}

	ctx.Next()
}
