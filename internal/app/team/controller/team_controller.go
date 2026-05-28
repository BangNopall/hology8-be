package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/contracts"
	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/domain/entity"
	"github.com/BangNopall/hology8-be/domain/enums"
	"github.com/BangNopall/hology8-be/internal/middlewares"
	"github.com/BangNopall/hology8-be/pkg/helpers/http/response"
)

type teamController struct {
	teamSvc contracts.TeamService
	logSvc  contracts.LogService
}

func InitTeamController(
	teamSvc contracts.TeamService,
	logSvc contracts.LogService,
	router *gin.Engine,
	middleware *middlewares.Middleware,
) {
	teamCtr := &teamController{
		teamSvc: teamSvc,
		logSvc:  logSvc,
	}

	teamRouter := router.Group("/api/v1/teams")
	teamRouter.GET("/:id", middleware.Authentication, teamCtr.FetchTeamData)
	teamRouter.GET("/user", middleware.Authentication, teamCtr.FetchUserTeams)
	teamRouter.GET("", middleware.Authentication, middleware.AuthorizationAdmin, teamCtr.FetchAllTeams)
	teamRouter.GET("/statistic", teamCtr.CountTeamNUniv)
	teamRouter.POST("", middleware.RateLimiter(), middleware.Authentication, middleware.VerifiedUser, teamCtr.CreateTeam)
	teamRouter.PUT("/:id/payment-proof", middleware.Authentication, middleware.AuthorizationTeam, teamCtr.UploadPaymentProof)
	teamRouter.PUT("/:id/twibbon-proof", middleware.Authentication, middleware.AuthorizationTeam, teamCtr.UploadTwibbonProof)
	teamRouter.PUT("/:id/proposal-doc", middleware.Authentication, middleware.AuthorizationTeam, teamCtr.UploadProposalDoc)
	teamRouter.PUT("/:id/statement-letter", middleware.Authentication, middleware.AuthorizationTeam, teamCtr.UploadStatementLetter)
	teamRouter.PUT("/:id", middleware.Authentication, middleware.AuthorizationTeam, teamCtr.UpdateTeam)
	teamRouter.PUT("/:id/status", middleware.Authentication, middleware.AuthorizationAdmin, teamCtr.UpdateTeamStatus)
	teamRouter.PUT("/join/:token", middleware.Authentication, middleware.VerifiedUser, teamCtr.JoinTeam)
	teamRouter.PUT("/:id/leader/:leaderId", middleware.Authentication, middleware.AuthorizationAdmin, teamCtr.UpdateLeader)
	teamRouter.DELETE("/:id/members/:memberId", middleware.Authentication, middleware.AuthorizationAdmin, teamCtr.DeleteMember)
}

// @Tags			Teams
// @Summary		Register a team
// @Description	Register a team to a competition
// @Produce		json
// @Param			TeamPayload	body		dto.TeamRegister		true	"Team Register Payload"
// @Success		200			{object}	response.Response		"OK"
// @Failure		401			{object}	response.ErrorResponse	"User's email is not verified yet"
// @Failure		400			{object}	response.ErrorResponse	"User already registered to same competition"
// @Failure		409			{object}	response.ErrorResponse	"Candidate leader is already a leader in another team"
// @Failure		500			{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/teams [post]
func (c *teamController) CreateTeam(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to create team"
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

	leaderIdString := ctx.GetString("id")
	leaderId, err := uuid.Parse(leaderIdString)
	if err != nil {
		return
	}

	var team dto.TeamRegister

	err = ctx.ShouldBindJSON(&team)
	if err != nil {

		return
	}

	err = c.teamSvc.CreateTeam(ctx, leaderId, team)
	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully create team"
}

// @Tags			Teams
// @Summary		Fetch Team Data
// @Description	Fetch team data
// @Produce		json
// @Param			id	path		string										true	"Team ID"
// @Success		200	{object}	response.Response{data=dto.TeamResponse}	"OK"
// @Failure		404	{object}	response.ErrorResponse						"Team Not Found"
// @Failure		500	{object}	response.ErrorResponse						"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/teams/{id} [get]
func (c *teamController) FetchTeamData(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to fetch team data"
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

	id := ctx.Param("id")

	res, err = c.teamSvc.FetchTeamData(ctx.Request.Context(), id)
	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully fetch team data"
}

// @Tags			Teams
// @Summary		Fetch All Teams (Admin Service)
// @Description	Fetch all teams
// @Produce		json
// @Param			competition_id	query		string												false	"Competition to filter"
// @Param			status			query		string												false	"Status to filter. Look up for enum Status to see what's available"
// @Param 			phase			query 		string												false 	"Phase to filter. Look up for enum Status to see what's available"
// @Param 			winner_place	query 		string												false 	"Winner place to filter. Look up for enum Status to see what's available"
// @Param			page			query		int													false	"Pagination Page"
// @Param			team_name		query		string												false	"Team name to filter"
// @Param			sort_by		query		string												false	"Sort by latest or oldest"
// @Success		200				{object}	response.Response{data=dto.TeamPaginationResponse}	"OK"
// @Failure		500				{object}	response.ErrorResponse								"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/teams [get]
func (c *teamController) FetchAllTeams(ctx *gin.Context) {
	var (
		err       error
		code      int = http.StatusBadRequest
		res       interface{}
		message   string = "failed to fetch all teams"
		pageParam *dto.PaginationRequest
		page      int
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

	compeId := ctx.Query("competition_id")
	queryPage := ctx.Query("page")
	queryStatus := ctx.Query("status")
	queryPhase := ctx.Query("phase")
	queryWinnerPlace := ctx.Query("winner_place")
	queryTeamName := ctx.Query("team_name")
	querySortBy := ctx.Query("sort_by")

	params := &dto.TeamParams{}

	params.Name = queryTeamName

	if compeId != "" {
		params.CompetitionID, err = strconv.Atoi(compeId)

		if err != nil {
			return
		}
	}

	if queryPage != "" {
		page, err = strconv.Atoi(queryPage)

		if err != nil {
			return
		}

		if page <= 0 {
			err = errors.New("invalid page number")
			return
		}

		pageParam = dto.NewPaginationRequest(page)
	}

	if queryStatus != "" {

		if !enums.IsValidStatus(queryStatus) {
			err = domain.ErrInvalidEnumInput
			return
		}

		params.Status = enums.Status(queryStatus)
	}

	if queryPhase != "" {
		if !enums.IsValidPhase(queryPhase) {
			err = domain.ErrInvalidEnumInput
			return
		}

		params.Phase = enums.Phase(queryPhase)
	}

	if queryWinnerPlace != "" {
		if !enums.IsValidWinnerPlace(queryWinnerPlace) {
			err = domain.ErrInvalidEnumInput
			return
		}

		params.WinnerPlace = enums.WinnerPlace(queryWinnerPlace)
	}

	if querySortBy != "latest" && querySortBy != "oldest" {
		params.SortBy = "latest"
	} else {
		params.SortBy = querySortBy
	}

	res, err = c.teamSvc.FetchAll(ctx, params, pageParam)
	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully fetch all teams"
}

// @Tags			Teams
// @Summary		Fetch User Teams
// @Description	Fetch the team user has joined or created
// @Produce		json
// @Success		200	{object}	response.Response{data=dto.UserTeamsResponse}	"OK"
// @Failure		500	{object}	response.ErrorResponse							"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/teams/user [get]
func (c *teamController) FetchUserTeams(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to fetch user's teams"
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

	userId := uuid.MustParse(ctx.GetString("id"))

	res, err = c.teamSvc.FetchUserTeams(ctx.Request.Context(), userId)
	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully fetch user's teams"
}

// @Tags			Teams
// @Summary		Upload payment proof
// @Description	Upload team payment proof
// @Produce		json
// @Param			id				path		string					true	"Team ID"
// @Param			payment_file	formData	file					true	"File to upload"
// @Success		200				{object}	response.Response		"OK"
// @Failure		401				{object}	response.ErrorResponse	"User is not part of the team"
// @Failure		413				{object}	response.ErrorResponse	"Requested file is too large"
// @Failure		500				{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/teams/{id}/payment-proof [put]
func (c *teamController) UploadPaymentProof(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to upload team payment proof"
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

	teamId := ctx.Param("id")

	file, err := ctx.FormFile("payment_file")

	if err != nil {
		return
	}

	err = c.teamSvc.UploadPaymentProof(ctx.Request.Context(), teamId, ctx.GetString("id"), file)

	code = domain.GetCode(err)

	if err != nil {

		return
	}

	message = "successfully upload payment proof"
}

// @Tags			Teams
// @Summary		Upload twibbon proof
// @Description	Upload team twibbon proof
// @Produce		json
// @Param			id				path		string					true	"Team ID"
// @Param			twibbon_file	formData	file					true	"File to upload"
// @Success		200				{object}	response.Response		"OK"
// @Failure		401				{object}	response.ErrorResponse	"User is not part of the team"
// @Failure		413				{object}	response.ErrorResponse	"Requested file is too large"
// @Failure		500				{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/teams/{id}/twibbon-proof [put]
func (c *teamController) UploadTwibbonProof(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to upload team twibbon proof"
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

	teamId := ctx.Param("id")

	file, err := ctx.FormFile("twibbon_file")

	if err != nil {

		return
	}

	err = c.teamSvc.UploadTwibbonProof(ctx.Request.Context(), teamId, ctx.GetString("id"), file)
	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully upload team twibbon proof"
}

// @Tags			Teams
// @Summary		Upload proposal doc
// @Description	Upload team proposal doc
// @Produce		json
// @Param			id				path		string					true	"Team ID"
// @Param			proposal_file	formData	file					true	"File to upload"
// @Success		200				{object}	response.Response		"OK"
// @Failure		401				{object}	response.ErrorResponse	"User is not part of the team"
// @Failure		500				{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/teams/{id}/proposal-doc [put]
func (c *teamController) UploadProposalDoc(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to upload team proposal doc"
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

	teamId := ctx.Param("id")

	file, err := ctx.FormFile("proposal_file")

	if err != nil {

		return
	}

	err = c.teamSvc.UploadProposalDoc(ctx.Request.Context(), teamId, ctx.GetString("id"), file)
	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully upload team proposal doc"
}

// @Tags			Teams
// @Summary		Upload statement letter
// @Description	Upload team statement letter
// @Produce		json
// @Param			id				path		string					true	"Team ID"
// @Param			statement_file	formData	file					true	"File to upload"
// @Success		200				{object}	response.Response		"OK"
// @Failure		401				{object}	response.ErrorResponse	"User is not part of the team"
// @Failure		500				{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/teams/{id}/statement-letter [put]
func (c *teamController) UploadStatementLetter(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to upload team statement letter"
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

	teamId := ctx.Param("id")

	file, err := ctx.FormFile("statement_file")

	if err != nil {

		return
	}

	err = c.teamSvc.UploadStatementLetter(ctx.Request.Context(), teamId, ctx.GetString("id"), file)
	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully upload team statement letter"
}

// @Tags			Teams
// @Summary		Update Team Data
// @Description	Upload team data
// @Produce		json
// @Param			id			path		string					true	"Team ID"
// @Param			TeamPayload	body		dto.TeamUpdate			true	"Team"
// @Success		200			{object}	response.Response		"OK"
// @Failure		401			{object}	response.ErrorResponse	"User is not part of the team"
// @Failure		403			{object}	response.ErrorResponse	"Forbidden to update some attributes"
// @Failure		404			{object}	response.ErrorResponse	"Team not found"
// @Failure		500			{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/teams/{id} [put]
func (c *teamController) UpdateTeam(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to update team data"
		team    dto.TeamUpdate
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

	id := ctx.Param("id")

	userId := ctx.GetString("id")

	if err = ctx.ShouldBindJSON(&team); err != nil {
		return
	}

	err = c.teamSvc.UpdateTeamData(ctx.Request.Context(), id, userId, &team)

	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully update team data"
}

// @Tags			Teams
// @Summary		Update Team Status, Phase or WinnerPlace (Admin Service)
// @Description	Update team status, phase, or winner place
// @Produce		json
// @Param			id			path		string					true	"Team ID"
// @Param			TeamPayload	body		dto.TeamUpdate			true	"Team"
// @Success		200			{object}	response.Response		"OK"
// @Failure		400			{object}	response.ErrorResponse	"Invalid enum status, phase or winner place update"
// @Failure		404			{object}	response.ErrorResponse	"Team not found"
// @Failure		500			{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/teams/{id}/status [put]
func (c *teamController) UpdateTeamStatus(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to update team status"
		team    dto.TeamUpdate
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

	id := ctx.Param("id")

	if err = ctx.ShouldBindJSON(&team); err != nil {
		response.SendErrResp(
			ctx,
			400,
			response.Fail,
			"failed to update team status",
			err,
		)
		return
	}

	err = c.teamSvc.UpdateTeamStatus(ctx.Request.Context(), id, &team)
	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully update team status"

	//logging
	c.logSvc.InsertLog(ctx.Request.Context(), &dto.LogRequest{
		AdminID: ctx.GetString("id"),
		Action:  "update team status with team id" + ctx.Param("id"),
	})
}

// @Tags			Teams
// @Summary		Delete Member (Admin Service)
// @Description	Delete a member from a team
// @Produce		json
// @Param			id			path		string					true	"Team ID"
// @Param			memberId	path		string					true	"Member ID"
// @Success		200			{object}	response.Response		"OK"
// @Failure		404			{object}	response.ErrorResponse	"Team or Member not found"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/teams/{id}/members/{memberId} [delete]
func (c *teamController) DeleteMember(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to delete member"
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

	teamId := ctx.Param("id")
	memberIdStr := ctx.Param("memberId")

	memberId, err := uuid.Parse(memberIdStr)

	if err != nil {
		return
	}

	err = c.teamSvc.RemoveMember(ctx.Request.Context(), &entity.DetailTeams{UserID: memberId, TeamID: teamId})

	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully delete member"

	//logging
	c.logSvc.InsertLog(ctx.Request.Context(), &dto.LogRequest{
		AdminID: ctx.GetString("id"),
		Action:  "remove user with id " + memberIdStr + " from team id " + teamId,
	})
}

// @Tags			Teams
// @Summary		Update Team Leader (Admin Service)
// @Description	Update team leader
// @Produce		json
// @Param			id			path		string					true	"Team ID"
// @Param			leaderId	path		string					true	"Leader ID"
// @Success		200			{object}	response.Response		"OK"
// @Failure		400			{object}	response.ErrorResponse	"New leader already registered to same competition"
// @Failure		404			{object}	response.ErrorResponse	"Team not found"
// @Failure		409			{object}	response.ErrorResponse	"Candidate leader is already a leader in another team"
// @Failure		500			{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/teams/{id}/leader/{leaderId} [put]
func (c *teamController) UpdateLeader(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to update leader"
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

	leaderIdStr := ctx.Param("leaderId")
	teamId := ctx.Param("id")

	leaderId, err := uuid.Parse(leaderIdStr)

	if err != nil {
		return
	}

	err = c.teamSvc.UpdateLeader(ctx.Request.Context(), teamId, leaderId)

	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully update leader"

	//logging
	c.logSvc.InsertLog(ctx.Request.Context(), &dto.LogRequest{
		AdminID: ctx.GetString("id"),
		Action:  "update leader from team with team id " + teamId + " with user " + leaderIdStr,
	})
}

// @Tags			Teams
// @Summary		Join team
// @Description	Join team by token
// @Produce		json
// @Param			token	path		string	true	"Token"
// @Success		200		{object}	response.Response{data=nil}
// @Failure		400		{object}	response.ErrorResponse	"Bad request"
// @Failure		401		{object}	response.ErrorResponse	"User's email is not verified yet"
// @Failure		404		{object}	response.ErrorResponse	"Team not found"
// @Failure		500		{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/teams/join/{token} [put]
func (c *teamController) JoinTeam(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to join team"
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

	token := ctx.Param("token")
	userIdString := ctx.GetString("id")
	userId, err := uuid.Parse(userIdString)

	if err != nil {

		return
	}

	err = c.teamSvc.JoinTeam(ctx.Request.Context(), token, userId)

	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully join team"
}

// @Tags			Teams
// @Summary		Count Team and University
// @Description	Count team and university
// @Produce		json
// @Success		200	{object}	response.Response{data=dto.TeamNUnivCounter}	"OK"
// @Failure		500	{object}	response.ErrorResponse							"Internal server error"
// @Security		ApiKeyAuth
// @Router			/teams/statistic [get]
func (c *teamController) CountTeamNUniv(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to fetch team statistic"
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

	res, err = c.teamSvc.CountTeamNUniv(ctx.Request.Context())
	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully fetch team statistic"
}
