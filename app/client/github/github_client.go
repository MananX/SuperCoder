package github

import (
	"ai-developer/app/client"
	"ai-developer/app/config"
	"ai-developer/app/monitoring"
	"ai-developer/app/types/request"
	"ai-developer/app/types/response"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type GithubClient struct {
	baseURL    string
	httpClient *client.HttpClient
	slackAlert *monitoring.SlackAlert
	logger     *zap.Logger
}

func NewGithubClient(
	httpClient *client.HttpClient,
	logger *zap.Logger,
	slackAlert *monitoring.SlackAlert,
) *GithubClient {
	return &GithubClient{
		baseURL:    "https://github.com",
		httpClient: httpClient,
		logger:     logger.Named("GithubClient"),
		slackAlert: slackAlert,
	}
}

func (c *GithubClient) FetchAccessToken(code string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/login/oauth/access_token", c.baseURL)

	payload := request.GithubAccessTokenRequest{
		ClientId:     config.GithubIntegrationClientId(),
		ClientSecret: config.GithubIntegrationClientSecret(),
		Code:         code,
		RedirectURI:  config.GithubIntegrationRedirectURL(),
	}

	headers := map[string]string{
		"Accept":       "*/*",
		"Content-Type": "application/json",
	}

	res, err := c.httpClient.Post(url, payload, headers)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create project, status code: %d", res.StatusCode)
	}

	var accessTokenResponse response.GithubAccessTokenResponse
	err = json.NewDecoder(res.Body).Decode(&accessTokenResponse)
	if err != nil {
		return nil, err
	}

	responseMap, err := structToMap(accessTokenResponse)
	if err != nil {
		return nil, err
	}

	return responseMap, nil
}

func structToMap(data interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
