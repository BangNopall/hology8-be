package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/contracts"
	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/internal/middlewares"
	"github.com/BangNopall/hology8-be/pkg/helpers/http/response"
)

type competitionController struct {
	competitionSvc contracts.CompetitionService
	logSvc         contracts.LogService
}

func InitCompetitionController(
	competitionSvc contracts.CompetitionService,
	logSvc contracts.LogService,
	router *gin.Engine,
	middleware *middlewares.Middleware,
) {
	compeCtr := &competitionController{
		competitionSvc: competitionSvc,
		logSvc:         logSvc,
	}

	compeRouter := router.Group("/api/v1/competitions")

	compeRouter.GET("", compeCtr.FetchAll)
	compeRouter.GET("/resource/admin", middleware.Authentication, middleware.AuthorizationAdmin, compeCtr.FetchAllAdmin)
	compeRouter.GET("/:id", compeCtr.FetchOne)
	compeRouter.GET("/resource/admin/:id", middleware.Authentication, middleware.AuthorizationAdmin, compeCtr.FetchOneAdmin)
	compeRouter.POST("", middleware.Authentication, middleware.AuthorizationAdmin, compeCtr.InsertCompe)
	compeRouter.PUT("/:id", middleware.Authentication, middleware.AuthorizationAdmin, compeCtr.UpdateCompe)
	compeRouter.DELETE("/:id", middleware.Authentication, middleware.AuthorizationAdmin, compeCtr.DeleteCompe)
}

// @Tags			Competitions
// @Summary		Fetch all competitions data
// @Description	Fetch all available competitions data
// @Produce		json
// @Success		200	{object}	response.Response{data=[]dto.CompetitionResponse}
// @Failure		500	{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Router			/competitions [get]
func (c *competitionController) FetchAll(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to fetch competitions"
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

	res, err = c.competitionSvc.FetchAll(ctx.Request.Context(), "")
	code = domain.GetCode(err)

	if err != nil {

		return
	}

	message = "successfully fetch competitions"
}

// @Tags			Competitions
// @Summary		Fetch all competitions data
// @Description	Fetch all available competitions data
// @Produce		json
// @Param			relation	query		string	false	"Relation to loaded. Available: team,announcement. Usage: ?relation=team"
// @Success		200			{object}	response.Response{data=[]dto.CompetitionResponse}
// @Failure		500			{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/competitions/resource/admin [get]
func (c *competitionController) FetchAllAdmin(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to fetch competitions"
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

	queryRelation := ctx.Query("relation")

	res, err = c.competitionSvc.FetchAll(ctx.Request.Context(), queryRelation)
	code = domain.GetCode(err)

	if err != nil {

		return
	}

	message = "successfully fetch competitions"
}

// @Tags			Competitions
// @Summary		Fetch competition data by ID
// @Description	Fetch one competition data by ID
// @Produce		json
// @Param			id	path		int												true	"Competition ID"
// @Success		200	{object}	response.Response{data=dto.CompetitionResponse}	"OK"
// @Failure		404	{object}	response.ErrorResponse							"Competition with requested ID not found"
// @Failure		500	{object}	response.ErrorResponse							"Internal server error"
// @Security		ApiKeyAuth
// @Router			/competitions/{id} [get]
func (c *competitionController) FetchOne(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to fetch competition"
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

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {

		return
	}

	res, err = c.competitionSvc.FetchOne(ctx.Request.Context(), id, "")

	code = domain.GetCode(err)

	if err != nil {

		return
	}

	message = "successfully fetch competition"
}

// @Tags			Competitions
// @Summary		Fetch competition data by ID
// @Description	Fetch one competition data by ID
// @Produce		json
// @Param			id			path		int												true	"Competition ID"
// @Param			relation	query		string											false	"Relation to loaded. Available: team,announcement. Usage: ?relation=team"
// @Success		200			{object}	response.Response{data=dto.CompetitionResponse}	"OK"
// @Failure		404			{object}	response.ErrorResponse							"Competition with requested ID not found"
// @Failure		500			{object}	response.ErrorResponse							"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/competitions/resource/admin{id} [get]
func (c *competitionController) FetchOneAdmin(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to fetch competition"
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

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {

		return
	}

	queryRelation := ctx.Query("relation")

	res, err = c.competitionSvc.FetchOne(ctx.Request.Context(), id, queryRelation)

	code = domain.GetCode(err)

	if err != nil {

		return
	}

	message = "successfully fetch competition"
}

func (c *competitionController) InsertCompe(ctx *gin.Context) {
	var compe dto.CompetitionRequest

	if err := ctx.ShouldBindJSON(&compe); err != nil {
		response.SendErrResp(ctx, http.StatusBadRequest, response.Fail, "failed to insert competition", err)
		return
	}

	err := c.competitionSvc.InsertCompe(ctx.Request.Context(), &compe)

	code := domain.GetCode(err)
	status := response.GetStatus(code)

	if err != nil {
		response.SendErrResp(ctx, code, status, "failed to insert competition", err)
		return
	}

	response.SendResp(ctx, code, status, "successfully insert competition", nil)
}

// @Tags			Competitions
// @Summary		Update competition data  (Admin Service)
// @Description	Update competition data
// @Accept			json
// @Produce		json
// @Param			id					path		int							true	"Competition ID"
// @Param			CompetitionPayload	body		dto.CompetitionRequest		true	"Competition Payload"
// @Success		200					{object}	response.Response{data=nil}	"OK"
// @Failure		400					{object}	response.ErrorResponse		"Bad data request"
// @Failure		401					{object}	response.ErrorResponse		"Unauthorized Admin"
// @Failure		404					{object}	response.ErrorResponse		"Competition with requested ID not found"
// @Failure		500					{object}	response.ErrorResponse		"Internal server error"
// @Security		UserAuth
// @Security		ApiKeyAuth
// @Router			/competitions/{id} [put]
func (c *competitionController) UpdateCompe(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to update competition"
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

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {

		return
	}

	var compe dto.CompetitionRequest
	compe.ID = id

	if err = ctx.ShouldBindJSON(&compe); err != nil {

		return
	}

	err = c.competitionSvc.UpdateCompe(ctx.Request.Context(), &compe)

	code = domain.GetCode(err)

	if err != nil {

		return
	}

	message = "successfully update competition"

	//logging
	c.logSvc.InsertLog(ctx.Request.Context(), &dto.LogRequest{
		AdminID: ctx.GetString("id"),
		Action:  "update compe id #" + ctx.Param("id"),
	})
}

func (c *competitionController) DeleteCompe(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		response.SendErrResp(ctx, http.StatusBadRequest, response.Fail, "failed to fetch competition", err)
		return
	}

	err = c.competitionSvc.DeleteCompe(ctx, id)

	code := domain.GetCode(err)
	status := response.GetStatus(code)

	if err != nil {
		response.SendErrResp(ctx, code, status, "failed to delete competition", err)
		return
	}

	response.SendResp(ctx, code, status, "successfully delete competition", nil)
}
