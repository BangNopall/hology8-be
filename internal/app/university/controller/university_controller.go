package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/dto"
	response "github.com/hology8/hology-be/pkg/helpers/http/response"
)

type universityController struct {
	uniSvc contracts.UniversityService
}

func InitUniversityController(uniSvc contracts.UniversityService, router *gin.Engine) {
	uniCtr := &universityController{uniSvc}

	uniRouter := router.Group("/api/v1/universities")

	uniRouter.GET("", uniCtr.FetchAll)
	uniRouter.GET("/:id", uniCtr.FetchByID)
}

// @Tags			Universities
// @Summary		Fetch all universities data
// @Description	Fetch all available universities data
// @Produce		json
// @Param			name	query		string	false	"University name"
// @Success		200		{object}	response.Response{data=[]dto.UniversityResponse}
// @Failure		500		{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Router			/universities [get]
func (c *universityController) FetchAll(ctx *gin.Context) {
	var (
		err     error
		code    int
		res     interface{}
		message string = "failed to fetch universities data"
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

	univName := ctx.Query("name")

	defer sendResp()

	res, err = c.uniSvc.FetchAll(ctx.Request.Context(), &dto.UniversityParam{Name: univName})

	code = domain.GetCode(err)

	if err != nil {

		return
	}

	message = "successfully fetch universities data"
}

// @Tags			Universities
// @Summary		Fetch universities data by ID
// @Description	Fetch one university data by ID
// @Produce		json
// @Param			id	path		int	true	"University ID"
// @Success		200	{object}	response.Response{data=[]dto.UniversityResponse}
// @Failure		404	{object}	response.ErrorResponse	"University with requested ID not found"
// @Failure		500	{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Router			/universities/{id} [get]
func (c *universityController) FetchByID(ctx *gin.Context) {
	var (
		err     error
		code    int
		res     interface{}
		message string = "failed to fetch university"
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
		code = http.StatusBadRequest
		message = ""
		return
	}

	res, err = c.uniSvc.FetchByID(ctx.Request.Context(), id)

	code = domain.GetCode(err)

	if err != nil {

		return
	}

	message = "successfully fetch university data"
}
