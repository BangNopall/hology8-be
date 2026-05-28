package oauth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	oauth2 "golang.org/x/oauth2"
	google "golang.org/x/oauth2/google"

	env "github.com/hology8/hology-be/internal/infra/env"
	log "github.com/hology8/hology-be/pkg/log"
)

type OauthInterface interface {
	GetConfig() *oauth2.Config
	GetUserInfo(c *gin.Context) (*http.Response, error)
}

type OauthStruct struct {
	Config *oauth2.Config
}

var Oauth = getOauth()

func getOauth() OauthInterface {
	redirectUrl := "https://api.hology.id/api/v1/users/oauth/callback"

	config := &oauth2.Config{
		ClientID:     env.AppEnv.GoogleClientID,
		ClientSecret: env.AppEnv.GoogleClientSecret,
		RedirectURL:  redirectUrl,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	return &OauthStruct{Config: config}
}

func (o *OauthStruct) GetConfig() *oauth2.Config {
	return o.Config
}

func (o *OauthStruct) GetUserInfo(ctx *gin.Context) (*http.Response, error) {
	code := ctx.Query("code")

	token, err := o.Config.Exchange(ctx, code)
	if err != nil {
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[OAUTH][GetUserInfo] failed to exchange code for token")

		return nil, err
	}

	client := o.Config.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[OAUTH][GetUserInfo] failed to get user info")

		return nil, err
	}

	return resp, nil
}
