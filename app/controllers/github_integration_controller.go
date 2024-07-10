package controllers

import (
	"ai-developer/app/client/github"
	"ai-developer/app/config"
	"ai-developer/app/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GithubIntegrationController struct {
	authService              *services.AuthService
	githubIntegrationService *services.GithubIntegrationService
	githubClient             *github.GithubClient
	clientID                 string
	clientSecret             string
	redirectURL              string
}

func (controller *GithubIntegrationController) OauthCallback(c *gin.Context) {
	fmt.Println("Inside Github Integration")

	var env = config.Get("app.env")
	fmt.Println("ENV : ", env)

	//if env == "development" {
	//	fmt.Println("Handling Skip Authentication Token.....")
	//	redirectURL, err := controller.authService.HandleDefaultAuth()
	//	if err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get default user token"})
	//		return
	//	}
	//	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
	//
	//}

	var code = c.Query("code")
	fmt.Println("CODE : ", code)

	var response, err = controller.githubClient.FetchAccessToken(code)
	fmt.Println(response)
	fmt.Println(err)

	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000")

	//var githubOauthConfig = &oauth2.Config{
	//	RedirectURL:  controller.redirectURL,
	//	ClientID:     controller.clientID,
	//	ClientSecret: controller.clientSecret,
	//	Scopes:       []string{"user:email"},
	//	Endpoint:     oauthGithub.Endpoint,
	//}
	//callback := githubOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOnline)
	//c.Redirect(http.StatusTemporaryRedirect, callback)
}

func NewGithubIntegrationController(
	githubIntegrationService *services.GithubIntegrationService,
	githubClient *github.GithubClient,
	authService *services.AuthService,
	clientID string,
	clientSecret string,
	redirectURL string,
) *GithubIntegrationController {
	return &GithubIntegrationController{
		authService:              authService,
		githubIntegrationService: githubIntegrationService,
		githubClient:             githubClient,
		clientID:                 clientID,
		clientSecret:             clientSecret,
		redirectURL:              redirectURL,
	}
}
