package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/domain/entity"
	"github.com/BangNopall/hology8-be/pkg/helpers/http/response"
)

func (m *Middleware) VerifiedUser(ctx *gin.Context) {
	id := ctx.GetString("id")

	userId, err := uuid.Parse(id)

	if err != nil {
		response.SendErrResp(
			ctx,
			http.StatusUnauthorized,
			response.Fail,
			"failed to verify user",
			err,
		)
		ctx.Abort()
		return
	}

	var user entity.User

	err = m.userRepo.FindUser(&user, &dto.UserParam{ID: userId})

	if err != nil {
		response.SendErrResp(
			ctx,
			http.StatusUnauthorized,
			response.Fail,
			"failed to verify user",
			err,
		)
		ctx.Abort()
		return
	}

	if !user.EmailIsVerified {
		response.SendErrResp(
			ctx,
			http.StatusUnauthorized,
			response.Fail,
			"user is not verified",
			errors.New("user email is not verified yet"),
		)

		ctx.Abort()
		return
	}

	ctx.Next()
}
