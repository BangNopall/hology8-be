package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/internal/middlewares"
	"github.com/hology8/hology-be/pkg/helpers/http/response"
)

type logController struct {
	logSvc contracts.LogService
}

func InitLogController(logSvc contracts.LogService, router *gin.Engine, middleware *middlewares.Middleware) {
	logController := &logController{
		logSvc: logSvc,
	}

	logRouter := router.Group("/api/v1/logs")
	logRouter.GET("", middleware.Authentication, middleware.AuthorizationAdmin, logController.FetchAllLogs)
}

// @Tags			Logs
// @Summary		Fetch all applications logs only by Admin
// @Description	Fetch all applications logs
// @Produce		json
// @Success		200	{object}	response.Response{data=[]dto.LogResponse}
// @Failure		500	{object}	response.ErrorResponse	"Internal server error"
// @Security 		UserAuth
// @Security		ApiKeyAuth
// @Router			/logs [get]
func (c *logController) FetchAllLogs(ctx *gin.Context) {
	logs, err := c.logSvc.FetchAllLogs(ctx.Request.Context())

	code := domain.GetCode(err)
	if err != nil {
		response.SendErrResp(
			ctx,
			code,
			response.Error,
			"failed to fetch logs",
			err,
		)

		return
	}

	response.SendResp(
		ctx,
		code,
		response.Success,
		"success to fetch logs",
		logs,
	)
}
