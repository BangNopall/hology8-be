package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/contracts"
	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/internal/middlewares"
	"github.com/BangNopall/hology8-be/pkg/helpers/http/response"
)

type AdminController struct {
	adminSvc   contracts.AdminService
	logSvc     contracts.LogService
	middleware *middlewares.Middleware
}

func InitAdminController(
	adminSvc contracts.AdminService,
	logSvc contracts.LogService,
	router *gin.Engine,
	middleware *middlewares.Middleware,
) {
	adminController := AdminController{
		adminSvc:   adminSvc,
		logSvc:     logSvc,
		middleware: middleware,
	}

	adminRouter := router.Group("api/v1/admins")
	adminRouter.POST("/login", adminController.Login)
	adminRouter.POST("/send-email", middleware.Authentication, middleware.AuthorizationAdmin, adminController.SendEmail)
}

// @Summary		Login
// @Description	Login for admin
// @Tags			Admins
// @Accept			json
// @Produce		json
// @Param			login	body		dto.AdminLogin			true	"Admin login"
// @Success		200		{object}	response.Response		"Login success"
// @Failure		400		{object}	response.ErrorResponse	"Status Bad Request"
// @Failure		401		{object}	response.ErrorResponse	"Status Unauthorized"
// @Failure		408		{object}	response.ErrorResponse	"Status Request Timeout"
// @Failure		500		{object}	response.ErrorResponse	"Status Internal Server Error"
// @Security		ApiKeyAuth
// @Router			/admins/login [post]
func (c *AdminController) Login(ctx *gin.Context) {
	var (
		err        error
		code       int    = http.StatusBadRequest
		message    string = "failed to login"
		res        interface{}
		adminLogin dto.AdminLogin
	)

	sendResp := func() {
		response.Send(
			ctx,
			code,
			message,
			res,
			err,
		)
	}

	defer sendResp()

	err = ctx.ShouldBindJSON(&adminLogin)
	if err != nil {
		return
	}

	res, err = c.adminSvc.Login(ctx, adminLogin)
	code = domain.GetCode(err)

	if err != nil {

		return
	}

	message = "successfully login with admin account"
}

// @Summary		Send Email
// @Description	Send email to all users/user who join a competition/a leader of a team
// @Tags			Admins
// @Accept			json
// @Produce		json
// @Param			to		query		string											true	"Relation to loaded. Available: all, team, competition. Usage example: ?to=all"
// @Param			email	body		dto.EmailMessage								true	"Email message. name is required if to is team/competition"
// @Success		200		{object}	response.Response{data=dto.AdminLoginResponse}	"Email sent"
// @Failure		400		{object}	response.ErrorResponse							"Status Bad Request"
// @Failure		401		{object}	response.ErrorResponse							"Status Unauthorized"
// @Failure		408		{object}	response.ErrorResponse							"Status Request Timeout"
// @Failure		500		{object}	response.ErrorResponse							"Status Internal Server Error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/admins/send-email [post]
func (c *AdminController) SendEmail(ctx *gin.Context) {
	var (
		err          error
		code         int    = http.StatusBadRequest
		message      string = "failed to send email"
		res          interface{}
		emailMessage dto.EmailMessage
	)

	sendResp := func() {
		response.Send(
			ctx,
			code,
			message,
			res,
			err,
		)
	}

	defer sendResp()

	err = ctx.ShouldBindJSON(&emailMessage)

	to := ctx.Request.FormValue("to")

	if err != nil {
		return
	}

	err = c.adminSvc.SendEmail(ctx, to, emailMessage)
	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully send email"

	c.logSvc.InsertLog(ctx.Request.Context(), &dto.LogRequest{
		AdminID: ctx.GetString("id"),
		Action:  "send email to " + to,
	})
}
