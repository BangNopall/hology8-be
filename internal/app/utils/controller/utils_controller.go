package controller

import (

	"github.com/gin-gonic/gin"

	response "github.com/hology8/hology-be/pkg/helpers/http/response"
)

type utilsController struct {
}

func InitUtilsController(router *gin.Engine) {
	utilsCtr := new(utilsController)

	router.GET("/health", utilsCtr.GetHealth)
}

func (c *utilsController) GetHealth(ctx *gin.Context) {
	response.SendResp(ctx, 200, response.Success, "server is running ok", nil)
}
