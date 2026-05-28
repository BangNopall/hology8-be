package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/contracts"
	"github.com/BangNopall/hology8-be/pkg/helpers/http/response"
)

type provinceController struct {
	provinceSvc contracts.ProvinceService
}

func InitProvinceController(provinceSvc contracts.ProvinceService, router *gin.Engine) {
	provCtr := &provinceController{provinceSvc}

	provRouter := router.Group("/api/v1/provinces")

	provRouter.GET("", provCtr.FetchAll)
	provRouter.GET("/:id", provCtr.FetchByID)
}

// @Tags			Provinces
// @Summary		Fetch all provinces data
// @Description	Fetch all available provinces data
// @Produce		json
// @Success		200	{object}	response.Response{data=[]dto.ProvinceResponse}
// @Failure		500	{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Router			/provinces [get]
func (c *provinceController) FetchAll(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to fetch provinces"
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

	res, err = c.provinceSvc.FetchAll(ctx.Request.Context())

	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "succcessfully fetch provinces"
}

// @Tags			Provinces
// @Summary		Fetch provinces data by ID
// @Description	Fetch one province data by ID
// @Produce		json
// @Param			id	path		int	true	"Province ID"
// @Success		200	{object}	response.Response{data=dto.ProvinceResponse}
// @Failure		404	{object}	response.ErrorResponse	"Province with requested ID not found"
// @Failure		500	{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Router			/provinces/{id} [get]
func (c *provinceController) FetchByID(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to fetch province"
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

	res, err = c.provinceSvc.FetchByID(ctx.Request.Context(), id)

	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "succcessfully fetch province"
}
