package controller

import (
	"net/http"
	"strconv"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/contracts"
	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/internal/middlewares"
	"github.com/BangNopall/hology8-be/pkg/helpers/http/response"
	"github.com/gin-gonic/gin"
)

type partnerController struct {
	partnerService contracts.PartnerService
}

func InitPartnerController(partnerService contracts.PartnerService, router *gin.Engine, middleware *middlewares.Middleware) {
	controller := &partnerController{partnerService}

	partnerRouter := router.Group("/api/v1/partners")

	partnerRouter.GET("", controller.FetchAll)
	partnerRouter.GET("/types", controller.FetchAllTypes)
	partnerRouter.GET("/:id", controller.FetchByID)
	partnerRouter.POST("", middleware.Authentication, middleware.AuthorizationAdmin, controller.Create)
	partnerRouter.PUT("/:id", middleware.Authentication, middleware.AuthorizationAdmin, controller.Update)
	partnerRouter.DELETE("/:id", middleware.Authentication, middleware.AuthorizationAdmin, controller.Delete)
}

// @Tags			Partners
// @Summary		Fetch all partners data
// @Description	Fetch all available partners data
// @Produce		json
// @Param			partner_type_id	query		string	false	"Partner type ID"
// @Success		200				{object}	response.Response{data=dto.PartnersResponse}
// @Failure		500				{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Router			/partners [get]
func (c *partnerController) FetchAll(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to fetch partners data"

		partnerParams dto.PartnerParams
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

	partnerTypeID := ctx.Query("partner_type_id")

	if partnerTypeID != "" {
		partnerParams.PartnerTypeID, err = strconv.Atoi(partnerTypeID)

		if err != nil {
			return
		}
	}

	res, err = c.partnerService.FetchAll(ctx.Request.Context(), &partnerParams)

	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully fetched partners data"
}

// @Tags			Partners
// @Summary		Fetch all partner types data
// @Description	Fetch all available partner types data
// @Produce		json
// @Success		200	{object}	response.Response{data=[]dto.PartnerTypeResponse}
// @Failure		500	{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Router			/partners/types [get]
func (c *partnerController) FetchAllTypes(ctx *gin.Context) {
	var (
		err     error
		code    int
		res     interface{}
		message string = "failed to fetch partner types data"
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

	res, err = c.partnerService.FetchAllTypes(ctx.Request.Context())
	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully fetch partner types data"
}

// @Tags			Partners
// @Summary		Fetch partner data by ID
// @Description	Fetch one partner data by ID
// @Produce		json
// @Param			id	path		int	true	"Partner ID"
// @Success		200	{object}	response.Response{data=dto.PartnerResponse}
// @Failure		404	{object}	response.ErrorResponse	"Partner with requested ID not found"
// @Failure		500	{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Router			/partners/{id} [get]
func (c *partnerController) FetchByID(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to fetch partner data"
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

	res, err = c.partnerService.FetchOneByID(ctx.Request.Context(), id)

	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully fetched partner data"
}

// @Tags			Partners
// @Summary		Create new partner
// @Description	Create new partner
// @Produce		json
// @Param			name			formData	string	true	"Partner name"
// @Param			partner_type_id	formData	int		true	"Partner type ID"
// @Param			image_file		formData	file	true	"Partner image file"
// @Success		201				{object}	response.Response{data=dto.PartnerResponse}
// @Failure		400				{object}	response.ErrorResponse	"Invalid request payload"
// @Failure		500				{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security 		UserAuth
// @Router			/partners [post]
func (c *partnerController) Create(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to create partner"
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

	var partnerCreate dto.PartnerCreate

	partnerCreate.Name = ctx.PostForm("name")
	partnerCreate.PartnerTypeID, err = strconv.Atoi(ctx.PostForm("partner_type_id"))
	if err != nil {
		return
	}

	file, err := ctx.FormFile("image_file")
	if err != nil {
		return
	}

	err = c.partnerService.CreatePartner(ctx.Request.Context(), &partnerCreate, file)

	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully created partner"
}

// @Tags			Partners
// @Summary		Update partner data
// @Description	Update partner data
// @Produce		json
// @Param			id				path		int		true	"Partner ID"
// @Param			name			formData	string	true	"Partner name"
// @Param			partner_type_id	formData	int		true	"Partner type ID"
// @Param			image_file		formData	file	true	"Partner image file"
// @Success		200				{object}	response.Response{data=dto.PartnerResponse}
// @Failure		400				{object}	response.ErrorResponse	"Invalid request payload"
// @Failure		404				{object}	response.ErrorResponse	"Partner with requested ID not found"
// @Failure		500				{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security 		UserAuth
// @Router			/partners/{id} [put]
func (c *partnerController) Update(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to update partner"
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

	var partnerUpdate dto.PartnerUpdate

	partnerUpdate.Name = ctx.PostForm("name")
	partnerUpdate.PartnerTypeID, err = strconv.Atoi(ctx.PostForm("partner_type_id"))
	if err != nil {
		return
	}

	file, err := ctx.FormFile("image_file")
	if err != nil {
		return
	}

	err = c.partnerService.UpdatePartner(ctx.Request.Context(), id, &partnerUpdate, file)

	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully updated partner"
}

// @Tags			Partners
// @Summary		Delete partner data
// @Description	Delete partner data
// @Produce		json
// @Param			id	path		int	true	"Partner ID"
// @Success		204	{object}	response.Response
// @Failure		404	{object}	response.ErrorResponse	"Partner with requested ID not found"
// @Failure		500	{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/partners/{id} [delete]
func (c *partnerController) Delete(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to delete partner"
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

	err = c.partnerService.DeletePartner(ctx.Request.Context(), id)

	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully deleted partner"
}
