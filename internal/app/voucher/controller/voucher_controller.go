package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/internal/middlewares"
	"github.com/hology8/hology-be/pkg/helpers/http/response"
)

type voucherController struct {
	voucherSvc contracts.VoucherService
}

func InitVoucherController(
	voucherSvc contracts.VoucherService,
	router *gin.Engine,
	middleware *middlewares.Middleware,
) {
	voucherCtr := &voucherController{
		voucherSvc: voucherSvc,
	}

	voucherRouter := router.Group("/api/v1/vouchers")
	voucherRouter.GET("", middleware.Authentication, middleware.AuthorizationAdmin, voucherCtr.FetchAll)
	voucherRouter.GET("/:id", middleware.Authentication, middleware.AuthorizationAdmin, voucherCtr.FetchByID)
	voucherRouter.POST("", middleware.Authentication, middleware.AuthorizationAdmin, voucherCtr.InsertVoucher)
	voucherRouter.POST("/redeem", middleware.RateLimiter(), middleware.Authentication, middleware.VerifiedUser, voucherCtr.RedeemVoucher)
}

// @Tags			Vouchers
// @Summary		Fetch all vouchers data (Admin Service)
// @Description	Fetch all available vouchers data
// @Produce		json
// @Success		200	{object}	response.Response{data=[]dto.VoucherResponse}
// @Failure		500	{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/vouchers [get]
func (c *voucherController) FetchAll(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to fetch vouchers"
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

	res, err = c.voucherSvc.FetchAll(ctx.Request.Context())
	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully fetch vouchers"
}

// @Tags			Vouchers
// @Summary		Fetch voucher data by ID (Admin Service)
// @Description	Fetch one voucher data by ID
// @Produce		json
// @Param			id	path		string	true	"Voucher ID"
// @Success		200	{object}	response.Response{data=dto.VoucherResponse}
// @Failure		404	{object}	response.ErrorResponse	"Voucher with requested ID not found"
// @Failure		500	{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/vouchers/{id} [get]
func (c *voucherController) FetchByID(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to fetch voucher"
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

	res, err = c.voucherSvc.FetchByID(ctx.Request.Context(), id)
	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully fetch voucher"
}

// @Tags			Vouchers
// @Summary		Insert voucher (Admin Service)
// @Description	Insert voucher to be redeemed by a team
// @Accept		json
// @Produce		json
// @Param			VoucherPayload body			dto.VoucherRequest		true 	"Voucher creation request"
// @Success		200	{object}	response.Response		"OK"
// @Failure		409 {object}	response.ErrorResponse	"Voucher ID already been used"
// @Failure		500	{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/vouchers [post]
func (c *voucherController) InsertVoucher(ctx *gin.Context) {
	var voucher dto.VoucherRequest

	if err := ctx.ShouldBindBodyWithJSON(&voucher); err != nil {
		response.SendErrResp(ctx, http.StatusBadRequest, response.Fail, "failed to insert voucher", err)
		return
	}

	err := c.voucherSvc.InsertVoucher(ctx, &voucher)
	code := domain.GetCode(err)
	status := response.GetStatus(code)

	if err != nil {
		response.SendErrResp(ctx, code, status, "failed to insert voucher", err)
		return
	}

	response.SendResp(ctx, code, status, "successfully insert voucher", nil)
}

// @Tags			Vouchers
// @Summary		Redeems voucher
// @Description	Redeems voucher that are provided by team
// @Accept		json
// @Produce		json
// @Param			VoucherPayload body			dto.VoucherRedeem		true 	"Voucher redeem request"
// @Success		200	{object}	response.Response		"OK"
// @Failure 	400 {object} 	response.ErrorResponse  "Bad request"
// @Failure 	404 {object} 	response.ErrorResponse 	"Voucher not found"
// @Failure		409 {object}	response.ErrorResponse	"Voucher already been redeemed"
// @Failure		500	{object}	response.ErrorResponse	"Internal server error"
// @Security		ApiKeyAuth
// @Security		UserAuth
// @Router			/vouchers/redeem [post]
func (c *voucherController) RedeemVoucher(ctx *gin.Context) {
	var voucher dto.VoucherRedeem

	if err := ctx.ShouldBindBodyWithJSON(&voucher); err != nil {
		response.SendErrResp(ctx, http.StatusBadRequest, response.Fail, "failed to redeem voucher", err)
		return
	}

	err := c.voucherSvc.RedeemVoucher(ctx.Request.Context(), &voucher)
	code := domain.GetCode(err)
	status := response.GetStatus(code)

	if err != nil {
		response.SendErrResp(ctx, code, status, "failed to redeem voucher", err)
		return
	}

	response.SendResp(ctx, code, status, "successfully redeem voucher", nil)
}
