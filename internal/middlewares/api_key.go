package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/BangNopall/hology8-be/internal/infra/env"
	"github.com/BangNopall/hology8-be/pkg/helpers/http/response"
)

func ApiKey() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if strings.Contains(ctx.FullPath(), "/api/v1/users/oauth/callback") {
			ctx.Next()
			return
		}

		headerReq := ctx.GetHeader("x-api-key")

		splitted := strings.Split(headerReq, " ")

		if len(splitted) < 2 {
			response.SendErrResp(
				ctx,
				http.StatusBadRequest,
				response.Fail,
				"failed to authenticate request",
				fmt.Errorf("invalid api key"),
			)
			ctx.Abort()
			return
		}

		headerKey := splitted[1]

		if headerKey != env.AppEnv.ApiKey {
			response.SendErrResp(
				ctx,
				http.StatusBadRequest,
				response.Fail,
				"failed to authenticate request",
				fmt.Errorf("invalid api key"),
			)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
