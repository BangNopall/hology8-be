package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/internal/middlewares"
	"github.com/hology8/hology-be/pkg/helpers/http/response"
	"github.com/hology8/hology-be/pkg/oauth"
	"github.com/hology8/hology-be/pkg/redis"
)

type userController struct {
	userSvc contracts.UserService
	oauth   oauth.OauthInterface
	redis   redis.RedisInterface
}

func InitUserController(
	userSvc contracts.UserService,
	router *gin.Engine,
	oauth oauth.OauthInterface,
	middleware *middlewares.Middleware,
	redis redis.RedisInterface,
) {
	userController := &userController{
		userSvc: userSvc,
		oauth:   oauth,
		redis:   redis,
	}

	userRouter := router.Group("/api/v1/users")
	userRouter.GET("/oauth", userController.Oauth)
	userRouter.GET("/oauth/callback", userController.OauthCallBack)
	userRouter.GET("/profile", middleware.Authentication, userController.FetchProfile)
	userRouter.GET("", middleware.Authentication, middleware.AuthorizationAdmin, userController.FetchAll)
	userRouter.POST("/register", middleware.RateLimiter(), userController.Register)
	userRouter.GET("/verify-email/:email/:emailVerPass", userController.VerifyEmail)
	userRouter.POST("/login", middleware.RateLimiter(), userController.LoginWithEmail)
	userRouter.POST("/forgot-password", middleware.RateLimiter(), userController.ForgotPassword)
	userRouter.PUT("/forgot-password/:token", userController.ResetPassword)
	userRouter.PUT("", middleware.Authentication, userController.UpdateUser)
	userRouter.POST("/ktm-image", middleware.RateLimiter(), middleware.Authentication, userController.UploadKtmImage)
	userRouter.POST("/proof", middleware.RateLimiter(), middleware.Authentication, userController.UploadProofImage)
	userRouter.POST("/logout", middleware.Authentication, userController.Logout)
}

// @Description	OAuth is used to authenticate user with Google
// @Tags			Users
// @Accept			json
// @Produce		json
// @Success		302	{object}	response.Response{data=dto.UserOuathLink}	"ok"
// @Security		ApiKeyAuth
// @Router			/users/oauth [get]
func (c *userController) Oauth(ctx *gin.Context) {
	url := c.oauth.GetConfig().AuthCodeURL("code", oauth2.AccessTypeOffline)

	resp := dto.UserOuathLink{
		RedirectLink: url,
	}

	c.redis.Set(ctx, "referer", ctx.Request.Referer(), time.Hour)

	response.SendResp(ctx, 200, response.Success, "please redirect to this url", resp)
}

// @Summary		Oauth Callback
// @Description	OAuth callback is used when Google needs to call my endpoint to process user information after authentication
// @Tags			Users
// @Accept			json
// @Produce		json
// @Success		200	{object}	response.Response		"ok"
// @Failure		404	{object}	response.ErrorResponse	"not found"
// @Failure		408	{object}	response.ErrorResponse	"request timeout"
// @Failure		500	{object}	response.ErrorResponse	"internal server error"
// @Security		ApiKeyAuth
// @Router			/users/oauth/callback [get]
func (c *userController) OauthCallBack(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusInternalServerError
		res     dto.UserLoginResponse
		message string = "failed to get user info from oauth"
	)

	defer func() {

		val, _ := c.redis.Get(ctx, "referer")
		c.redis.Delete(ctx, "referer")

		if err != nil {
			ctx.Redirect(
				301,
				fmt.Sprintf(
					"%soauth/login/redirect?code=%v&message=%v",
					val,
					code,
					message,
				),
			)
		} else {
			ctx.Redirect(
				301,
				fmt.Sprintf(
					"%soauth/login/redirect?token=%v&code=%v&message=%v",
					val,
					res.Token,
					code,
					message,
				),
			)

		}

	}()

	resp, err := c.oauth.GetUserInfo(ctx)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	var user dto.UserOauth

	err = json.NewDecoder(resp.Body).Decode(&user)

	if err != nil {
		return
	}

	res, err = c.userSvc.LoginRegisterOauth(ctx.Request.Context(), user)
	code = domain.GetCode(err)

	if err != nil {
		message = "failed to register/login user"
		return
	}

	message = "successfully login with google"
}

// @Summary		Register
// @Description	Register for user
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			register	body		dto.UserRegister		true	"User register"
// @Success		200			{object}	response.Response		"ok"
// @Failure		400			{object}	response.ErrorResponse	"bad request"
// @Failure		401			{object}	response.ErrorResponse	"unauthorized"
// @Failure		408			{object}	response.ErrorResponse	"request timeout"
// @Failure		500			{object}	response.ErrorResponse	"internal server error"
// @Security		ApiKeyAuth
// @Router			/users/register [post]
func (c *userController) Register(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to register account"
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

	var user dto.UserRegister

	err = ctx.ShouldBindBodyWithJSON(&user)

	if err != nil {
		return
	}

	err = c.userSvc.Register(ctx.Request.Context(), user, ctx.Request.Referer())
	code = domain.GetCode(err)

	if err != nil {

		return
	}

	message = "success to register account, please check your email"
}

// @Summary		Verify Email
// @Description	Verify email for user
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			email			path		string					true	"Email"
// @Param			emailVerPass	path		string					true	"Email verification password"
// @Success		200				{object}	response.Response		"ok"
// @Failure		401				{object}	response.ErrorResponse	"unauthorized"
// @Failure		408				{object}	response.ErrorResponse	"request timeout"
// @Failure		500				{object}	response.ErrorResponse	"internal server error"
// @Security		ApiKeyAuth
// @Router			/users/verify-email/{email}/{emailVerPass} [get]
func (c *userController) VerifyEmail(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to verify email"
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

	email := ctx.Param("email")
	emailVerPass := ctx.Param("emailVerPass")

	err = c.userSvc.VerifyEmail(ctx.Request.Context(), email, emailVerPass)
	code = domain.GetCode(err)

	if err != nil {

		return
	}

	message = "success to verify email"
}

// @Summary		Login with email
// @Description	Login with email for user
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			login	body		dto.UserLogin			true	"User login"
// @Success		200		{object}	dto.UserLoginResponse	"ok"
// @Failure		400		{object}	response.ErrorResponse	"bad request"
// @Failure		401		{object}	response.ErrorResponse	"unauthorized"
// @Failure		408		{object}	response.ErrorResponse	"request timeout"
// @Failure		500		{object}	response.ErrorResponse	"internal server error"
// @Security		ApiKeyAuth
// @Router			/users/login [post]
func (c *userController) LoginWithEmail(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to login with email"
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

	var user dto.UserLogin

	err = ctx.ShouldBindBodyWithJSON(&user)
	if err != nil {

		return
	}

	res, err = c.userSvc.LoginWithEmail(ctx.Request.Context(), user)
	code = domain.GetCode(err)

	if err != nil {

		return
	}

	message = "successfully login with email"
}

// @Summary		Reset Password
// @Description	This api is used to reset user password
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			token			path		string											true	"Token"
// @Param			resetpassword	body		dto.UserResetPassword							true	"User reset password"
// @Success		200				{object}	response.Response{data=dto.UserLoginResponse}	"ok"
// @Failure		400				{object}	response.ErrorResponse							"bad request"
// @Failure		401				{object}	response.ErrorResponse							"unauthorized"
// @Failure		408				{object}	response.ErrorResponse							"request timeout"
// @Failure		500				{object}	response.ErrorResponse							"internal server error"
// @Security		ApiKeyAuth
// @Router			/users/forgot-password/{token} [put]
func (c *userController) ResetPassword(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to reset password"
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
	var user dto.UserResetPassword

	err = ctx.ShouldBindBodyWithJSON(&user)
	if err != nil {

		return
	}

	err = c.userSvc.ResetPassword(ctx.Request.Context(), user, token)
	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "success to reset password"
}

// @Summary		Forgot Password
// @Description	This api is used to reset password, user will get email to reset passwor after hit this endpoint
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			forgotpassword	body		dto.UserForgotPassword	true	"User forgot password"
// @Success		200				{object}	response.Response		"ok"
// @Failure		400				{object}	response.ErrorResponse	"bad request"
// @Failure		404				{object}	response.ErrorResponse	"not found"
// @Failure		408				{object}	response.ErrorResponse	"request timeout"
// @Failure		500				{object}	response.ErrorResponse	"internal server error"
// @Security		ApiKeyAuth
// @Router			/users/forgot-password [post]
func (c *userController) ForgotPassword(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to create forgot password token"
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

	var user dto.UserForgotPassword

	err = ctx.ShouldBindBodyWithJSON(&user)
	if err != nil {

		return
	}

	err = c.userSvc.ForgotPassword(ctx.Request.Context(), user, ctx.Request.Referer())
	code = domain.GetCode(err)

	if err != nil {

		return
	}

	message = "successfully create forgot password token, please check your email"
}

// @Summary		Update User
// @Description	Update user profile
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			updateUser	body		dto.UserUpdate			true	"User update"
// @Success		200			{object}	response.Response		"ok"
// @Failure		400			{object}	response.ErrorResponse	"bad request"
// @Failure		408			{object}	response.ErrorResponse	"request timeout"
// @Failure		500			{object}	response.ErrorResponse	"internal server error"
// @Security		UserAuth
// @Security		ApiKeyAuth
// @Router			/users [put]
func (c *userController) UpdateUser(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to update user"
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

	userIdString := ctx.GetString("id")
	userId, err := uuid.Parse(userIdString)
	if err != nil {

		return
	}

	var updateUser dto.UserUpdate

	err = ctx.ShouldBindBodyWithJSON(&updateUser)
	if err != nil {

		return
	}

	err = c.userSvc.UpdateUser(ctx.Request.Context(), userId, updateUser)
	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully update user"
}

// @Summary		Upload KTM Image
// @Description	Upload KTM image for user
// @Tags			Users
// @Accept			multipart/form-data
// @Produce		json
// @Param			ktm_image_link	formData	file					true	"KTM image"
// @Success		200				{object}	response.Response		"ok"
// @Failure		400				{object}	response.ErrorResponse	"bad request"
// @Failure		408				{object}	response.ErrorResponse	"request timeout"
// @Failure		413				{object}	response.ErrorResponse	"request entity too large"
// @Failure		500				{object}	response.ErrorResponse	"internal server error"
// @Security		UserAuth
// @Security		ApiKeyAuth
// @Router			/users/ktm-image [post]
func (c *userController) UploadKtmImage(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to upload ktm image"
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

	userIdString := ctx.GetString("id")
	userId, err := uuid.Parse(userIdString)
	if err != nil {

		return
	}

	ktmImage, err := ctx.FormFile("ktm_image_link")
	if err != nil {

		return
	}

	err = c.userSvc.UploadKtmImage(ctx.Request.Context(), userId, ktmImage)
	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "success to upload ktm image"
}

// @Summary		Upload Proof Image
// @Description	Upload proof image for user
// @Tags			Users
// @Accept			multipart/form-data
// @Produce		json
// @Param			follow-proof	formData	file					false	"Follow Proof Image"
// @Param			share-proof		formData	file					false	"Share Proof Image"
// @Param			proof			query		string					false	"Proof type. Available: follow-proof, share-proof"
// @Success		200				{object}	response.Response		"ok"
// @Failure		400				{object}	response.ErrorResponse	"bad request"
// @Failure		408				{object}	response.ErrorResponse	"request timeout"
// @Failure		413				{object}	response.ErrorResponse	"request entity too large"
// @Failure		500				{object}	response.ErrorResponse	"internal server error"
// @Security		UserAuth
// @Security		ApiKeyAuth
// @Router			/users/proof [post]
func (c *userController) UploadProofImage(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to upload proof image"
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

	userIdString := ctx.GetString("id")
	userId, err := uuid.Parse(userIdString)
	if err != nil {

		return
	}

	proofQuery := ctx.Query("proof")

	if proofQuery == "" {
		err = domain.ErrInvalidProofType
		return
	}

	fileImage, err := ctx.FormFile(proofQuery)
	if err != nil {

		return
	}

	err = c.userSvc.UploadProofImage(ctx.Request.Context(), userId, proofQuery, fileImage)
	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "success to upload proof image"
}

// @Summary		Fetch All Users (Admin Service)
// @Description	Fetch all users data
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			page		query		int													false	"Pagination Page"
// @Success		200			{object}	response.Response{data=dto.UserPaginationResponse}	"ok"
// @Failure		400			{object}	response.ErrorResponse								"Bad request"
// @Failure		500			{object}	response.ErrorResponse								"internal server error"
// @Security		UserAuth
// @Security		ApiKeyAuth
// @Router			/users [get]
func (c *userController) FetchAll(ctx *gin.Context) {
	var (
		err      error
		code     int = http.StatusBadRequest
		res      interface{}
		pageResp dto.PaginationResponse
		message  string = "failed to fetch all users"
		page     int

		userParam dto.UserParam
		pageParam *dto.PaginationRequest
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

	queryPage := ctx.Query("page")

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

	users, pageResp, err := c.userSvc.FetchAll(ctx, &userParam, pageParam)

	code = domain.GetCode(err)

	if err != nil {
		return
	}

	res = dto.UserPaginationResponse{
		Users:      users,
		Pagination: pageResp,
	}

	message = "successfully fetch all users"
}

// @Summary		Fetch Profile
// @Description	Fetch user profile
// @Tags			Users
// @Accept			json
// @Produce		json
// @Success		200	{object}	response.Response{data=dto.UserResponse}	"ok"
// @Failure		400	{object}	response.ErrorResponse						"bad request"
// @Failure		408	{object}	response.ErrorResponse						"request timeout"
// @Failure		500	{object}	response.ErrorResponse						"internal server error"
// @Security		UserAuth
// @Security		ApiKeyAuth
// @Router			/users/profile [get]
func (c *userController) FetchProfile(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to fetch user profile"
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

	id := ctx.GetString("id")

	uuid, err := uuid.Parse(id)

	if err != nil {

		return
	}

	res, err = c.userSvc.FetchByParam(ctx.Request.Context(), &dto.UserParam{ID: uuid})

	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully fetch user profile"
}

// @Summary		Logout
// @Description	Logout user
// @Tags			Users
// @Accept			json
// @Produce		json
// @Success		200	{object}	response.Response		"ok"
// @Failure		400	{object}	response.ErrorResponse	"bad request"
// @Failure		500	{object}	response.ErrorResponse	"internal server error"
// @Security		UserAuth
// @Security		ApiKeyAuth
// @Router			/users/logout [post]
func (c *userController) Logout(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to logout user"
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

	bearerToken := ctx.GetHeader("Authorization")

	token := strings.Split(bearerToken, " ")

	err = c.userSvc.Logout(ctx.Request.Context(), token[1])

	code = domain.GetCode(err)

	if err != nil {
		return
	}

	message = "successfully logout user"
}
