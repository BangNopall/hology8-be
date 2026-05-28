package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/internal/middlewares"
	"github.com/hology8/hology-be/pkg/helpers/http/response"
)

type presenceController struct {
	presenceSvc contracts.PresenceService
	logSvc      contracts.LogService
}

func InitPresenceController(
	presenceSvc contracts.PresenceService,
	logSvc contracts.LogService,
	router *gin.Engine,
	middleware *middlewares.Middleware,
) {
	presenceCtr := &presenceController{
		presenceSvc: presenceSvc,
		logSvc:      logSvc,
	}

	presences := router.Group("/api/v1/presences")

	presences.POST("", middleware.Authentication, middleware.AuthorizationAdmin, presenceCtr.CreatePresence)

	presences.GET("/:user_id", middleware.Authentication, presenceCtr.CheckPresence)

	presences.GET("", middleware.Authentication, middleware.AuthorizationAdmin, presenceCtr.ListPresences)
}

// @Tags         Presences
// @Summary      Create Presence
// @Description  Create presence for a user (admin only)
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.PresenceCreateRequest  true  "Create presence payload"
// @Success      200      {object}  response.Response{data=dto.PresenceCreateResponse}  "successfully create presence"
// @Example      {json}   200-success  {
// @Example      {json}   200-success  "code": 200,
// @Example      {json}   200-success  "status": "success",
// @Example      {json}   200-success  "message": "successfully create presence",
// @Example      {json}   200-success  "data": {
// @Example      {json}   200-success    "user_id": "3bf354d4-9915-495a-a362-18ba15ea19ae",
// @Example      {json}   200-success    "fullname": "Lovely Ito Panjaitan",
// @Example      {json}   200-success    "team_name": "ITO PASTI MENANG",
// @Example      {json}   200-success    "created_at": "2025-10-09T23:23:26.894274Z"
// @Example      {json}   200-success  }
// @Example      {json}   200-success}
// @Failure      400      {object}  response.Response{data=any}  "Bad request"
// @Failure      401      {object}  response.Response{data=any}  "Unauthorized"
// @Failure      404      {object}  response.Response{data=any}  "User not found"
// @Failure      409      {object}  response.Response{data=any}  "Presence already exists"
// @Failure      500      {object}  response.Response{data=any}  "Internal server error"
// @Security     ApiKeyAuth
// @Router       /presences [post]
func (c *presenceController) CreatePresence(ctx *gin.Context) {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to create presence"
		payload dto.PresenceCreateRequest
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

	// bind and validate request body
	if err = ctx.ShouldBindJSON(&payload); err != nil {
		return
	}

	// call service
	result, svcErr := c.presenceSvc.CreatePresence(ctx.Request.Context(), payload.UserID)
	err = svcErr
	code = domain.GetCode(err)
	if err != nil {
		return
	}

	message = "successfully create presence"
	res = result

	// logging (admin action)
	adminID := ctx.GetString("id")
	if _, parseErr := uuid.Parse(adminID); parseErr == nil {
		c.logSvc.InsertLog(ctx.Request.Context(), &dto.LogRequest{
			AdminID: adminID,
			Action:  "create presence for user " + payload.UserID.String(),
		})
	}
}

// @Tags         Presences
// @Summary      Check user's presence
// @Description  Check whether a user has presenced (by user_id)
// @Produce      json
// @Param        user_id  path      string  true  "User ID (UUID)"
// @Success      200      {object}  response.Response{data=dto.PresenceCheckResponse}  "presence status"
// @Example      {json}   200-success  {
// @Example      {json}   200-success  "code": 200,
// @Example      {json}   200-success  "status": "success",
// @Example      {json}   200-success  "message": "presence status",
// @Example      {json}   200-success  "data": {
// @Example      {json}   200-success    "user_id": "3bf354d4-9915-495a-a362-18ba15ea19ae",
// @Example      {json}   200-success    "exists": true,
// @Example      {json}   200-success    "created_at": "2025-10-09T22:46:20.320507Z"
// @Example      {json}   200-success  }
// @Example      {json}   200-success}
// @Failure      400      {object}  response.Response{data=any}                        "Bad request"
// @Failure      404      {object}  response.Response{data=dto.PresenceCheckResponse}  "Presence not found"
// @Security     ApiKeyAuth
// @Router       /presences/{user_id} [get]
func (c *presenceController) CheckPresence(ctx *gin.Context) {
	var (
		err     error
		code    = http.StatusOK
		res     interface{}
		message = "presence status"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	userIDStr := ctx.Param("user_id")
	userID, parseErr := uuid.Parse(userIDStr)
	if parseErr != nil {
		err = parseErr
		code = http.StatusBadRequest
		message = "invalid user_id"
		return
	}

	result, svcErr := c.presenceSvc.CheckPresence(ctx.Request.Context(), userID)
	err = svcErr
	code = domain.GetCode(err)
	if err != nil {
		return
	}

	// If not found, set HTTP 404
	if !result.Exists {
		code = http.StatusNotFound
		message = "presence not found"
		res = result
		return
	}

	res = result
}

// @Tags         Presences
// @Summary      List presences (admin)
// @Description  Paginated list sorted by created_at DESC (newest first)
// @Produce      json
// @Param        page  query     int  false  "Page number (default 1)"
// @Success      200   {object}  response.Response{data=dto.PresenceListResponse}  "successfully fetch presences"
// @Example      {json} 200-success  {
// @Example      {json} 200-success  "code": 200,
// @Example      {json} 200-success  "status": "success",
// @Example      {json} 200-success  "message": "successfully fetch presences",
// @Example      {json} 200-success  "data": {
// @Example      {json} 200-success    "presences": [
// @Example      {json} 200-success      {
// @Example      {json} 200-success        "user_id": "3bf354d4-9915-495a-a362-18ba15ea19ae",
// @Example      {json} 200-success        "fullname": "Lovely Ito Panjaitan",
// @Example      {json} 200-success        "team_name": "ITO PASTI MENANG",
// @Example      {json} 200-success        "created_at": "2025-10-09T23:23:26.894274Z"
// @Example      {json} 200-success      }
// @Example      {json} 200-success    ],
// @Example      {json} 200-success    "pagination": { "total_pages": 1, "page": 1 }
// @Example      {json} 200-success  }
// @Example      {json} 200-success}
// @Failure      401   {object}  response.Response{data=any}  "Unauthorized"
// @Security     ApiKeyAuth
// @Router       /presences [get]
func (c *presenceController) ListPresences(ctx *gin.Context) {
	var (
		err     error
		code    = http.StatusOK
		res     interface{}
		message = "successfully fetch presences"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}

	pReq := dto.NewPaginationRequest(page)

	items, pagination, svcErr := c.presenceSvc.FetchPresences(ctx.Request.Context(), pReq)
	err = svcErr
	code = domain.GetCode(err)
	if err != nil {
		return
	}

	res = dto.PresenceListResponse{
		Presences:  items,
		Pagination: pagination,
	}
}
