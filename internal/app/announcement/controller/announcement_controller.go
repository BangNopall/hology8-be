package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/contracts"
	"github.com/BangNopall/hology8-be/domain/dto"

	"github.com/BangNopall/hology8-be/internal/middlewares"
	"github.com/BangNopall/hology8-be/pkg/helpers/http/response"
)

type announcementController struct {
	announcementSvc contracts.AnnouncementService
	logSvc          contracts.LogService
}

func InitAnnouncementController(
	announcementSvc contracts.AnnouncementService,
	logSvc contracts.LogService,
	router *gin.Engine,
	middleware *middlewares.Middleware,
) {
	announceController := &announcementController{
		announcementSvc,
		logSvc,
	}

	announcements := router.Group("/api/v1/announcements")
	announcements.GET("", middleware.Authentication, middleware.AuthorizationAdmin, announceController.FetchAnnouncementByTo)
	announcements.POST("/", middleware.Authentication, middleware.AuthorizationAdmin, announceController.CreateAnnouncement)
	announcements.PUT("/:id", middleware.Authentication, middleware.AuthorizationAdmin, announceController.UpdateAnnouncement)
	announcements.DELETE("/:id", middleware.Authentication, middleware.AuthorizationAdmin, announceController.DeleteAnnouncement)
}

// @Tags			Announcements
// @Summary		Fetch announcements by To
// @Description	Fetch announcements by to. If team_id and competition_id is empty, then it will fetch all announcements to all users
// @Produce		json
// @Param			team_id			query		string	false	"Team ID"
// @Param			competition_id	query		string	false	"Competition ID"
// @Success		200				{object}	response.Response{data=[]dto.AnnouncementResponse}
// @Failure		400				{object}	response.ErrorResponse	"Bad request"
// @Failure		500				{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/announcements [get]
func (announcementController *announcementController) FetchAnnouncementByTo(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to fetch announcements"
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

	teamID := ctx.Query("team_id") // if not exit then ""
	competitionID := ctx.Query("competition_id")
	var competitionIDInt int // if not exit then 0

	if competitionID != "" {
		competitionIDIntTemp, err := strconv.Atoi(competitionID)

		if err != nil {
			return
		}

		competitionIDInt = competitionIDIntTemp
	}

	res, err = announcementController.announcementSvc.FetchAnnouncementByTo(ctx.Request.Context(), teamID, competitionIDInt)
	code = domain.GetCode(err)

	if err != nil {

		return
	}

	message = "successfully fetch announcements"
}

// @Tags			Announcements
// @Summary		Create announcement
// @Description	Create announcement only by admin
// @Accept			json
// @Produce		json
// @Param			AnnouncementRequest	body		dto.AnnouncementRequest	true	"Announcement request body"
// @Success		200					{object}	response.Response{data=nil}
// @Failure		400					{object}	response.ErrorResponse	"Bad request"
// @Failure		404					{object}	response.ErrorResponse	"Team/Competition not found"
// @Failure		500					{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/announcements [post]
func (announcementController *announcementController) CreateAnnouncement(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to create announcement"
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

	var payload dto.AnnouncementRequest

	if err = ctx.ShouldBindJSON(&payload); err != nil {
		return
	}

	adminID, err := uuid.Parse(ctx.GetString("id"))

	if err != nil {
		return
	}

	payload.AdminID = adminID

	err = announcementController.announcementSvc.CreateAnnouncement(ctx.Request.Context(), &payload)

	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully create announcements"

	// logging
	var logTo string
	var logIdTo interface{}
	if payload.TeamID == "" {
		logTo = "Compe ID"
		logIdTo = payload.CompetitionID
	} else {
		logTo = "Team ID"
		logIdTo = payload.TeamID
	}
	announcementController.logSvc.InsertLog(ctx.Request.Context(), &dto.LogRequest{
		AdminID: ctx.GetString("id"),
		Action: fmt.Sprintf(`
		create announcement with description %v to %v %v
		`, payload.Description, logTo, logIdTo),
	})
}

// @Tags			Announcements
// @Summary		Update announcement
// @Description	Update announcement only by admin
// @Produce		json
// @Param			id					path		int						true	"Announcement ID"
// @Param			AnnouncementRequest	body		dto.AnnouncementRequest	true	"Announcement request body"
// @Success		200					{object}	response.Response{data=nil}
// @Failure		400					{object}	response.ErrorResponse	"Bad request"
// @Failure		404					{object}	response.ErrorResponse	"Announcement/Team/Competition not found"
// @Failure		500					{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/announcements/{id} [put]
func (announcementController *announcementController) UpdateAnnouncement(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to update announcement"
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
	var payload dto.AnnouncementRequest

	if err != nil {
		return
	}

	if err = ctx.ShouldBindJSON(&payload); err != nil {
		return
	}

	adminID, err := uuid.Parse(ctx.GetString("id"))

	if err != nil {
		return
	}

	payload.ID = id
	payload.AdminID = adminID

	err = announcementController.announcementSvc.UpdateAnnouncement(ctx.Request.Context(), &payload)
	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully update announcements"

	//logging
	announcementController.logSvc.InsertLog(ctx.Request.Context(), &dto.LogRequest{
		AdminID: ctx.GetString("id"),
		Action:  "update announcement id #" + ctx.Param("id"),
	})
}

// @Tags			Announcements
// @Summary		Delete announcement
// @Description	Delete announcement only by admin
// @Produce		json
// @Param			id	path		int	true	"Announcement ID"
// @Success		200	{object}	response.Response{data=nil}
// @Failure		400	{object}	response.ErrorResponse	"Bad request"
// @Failure		404	{object}	response.ErrorResponse	"Announcement/Team/Competition not found"
// @Failure		500	{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/announcements/{id} [delete]
func (announcementController *announcementController) DeleteAnnouncement(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to delete announcement"
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

	err = announcementController.announcementSvc.DeleteAnnouncement(ctx.Request.Context(), id)

	code = domain.GetCode(err)

	if err != nil {

		return
	}

	message = "successfully delete announcements"

	//logging
	announcementController.logSvc.InsertLog(ctx.Request.Context(), &dto.LogRequest{
		AdminID: ctx.GetString("id"),
		Action:  "delete announcement id #" + ctx.Param("id"),
	})
}
