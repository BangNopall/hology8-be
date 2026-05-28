package response

import "github.com/gin-gonic/gin"

type Status string

const (
	// Success thrown when all went well
	Success Status = "success"

	/* Failed thrown when there is a problem from request data
	ex: (validation error, failed to bind json data, etc) */
	Fail Status = "fail"

	/* Error thrown when an error occured in processing the request aka server error */
	Error Status = "error"
)

type Response struct {
	Code    int         `json:"code"`
	Status  Status      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Status  Status `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func SendResp(
	ctx *gin.Context,
	code int,
	status Status,
	message string,
	data interface{},
) {
	ctx.JSON(code, Response{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func SendErrResp(
	ctx *gin.Context,
	code int,
	status Status,
	message string,
	err error,
) {
	errMsg := ""

	if err != nil {
		errMsg = err.Error()
	}

	ctx.JSON(code, ErrorResponse{
		Code:    code,
		Status:  status,
		Message: message,
		Error:   errMsg,
	})
}

func Send(
	ctx *gin.Context,
	code int,
	message string,
	data interface{},
	err error,
) {
	status := GetStatus(code)

	if err != nil {
		SendErrResp(
			ctx,  
			code, 
			status,
			message,
			err,
		)	
	} else {
		SendResp(
			ctx,
			code,
			status,
			message,
			data,
		)
	}

}

func GetStatus(code int) Status {
	if code >= 500 {
		return Error
	}

	if code >= 400 {
		return Fail
	}

	return Success
}
