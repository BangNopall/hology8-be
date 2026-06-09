package middlewares

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://elb-ccs-304817220.us-east-1.elb.amazonaws.com"},
		AllowMethods:     []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodHead, http.MethodDelete, http.MethodOptions, http.MethodPut},
		AllowHeaders:     []string{"Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization", "X-API-Key", "X-Cursor", "Token-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}
