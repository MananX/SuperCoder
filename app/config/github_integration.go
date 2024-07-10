package config

func GithubIntegrationClientId() string { return config.String("github.integration.client.id") }

func GithubIntegrationClientSecret() string { return config.String("github.integration.client.secret") }

func GithubIntegrationRedirectURL() string { return config.String("github.integration.redirect.url") }
