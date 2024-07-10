package services

type GithubIntegrationService struct {
	jwtService          *JWTService
	userService         *UserService
	organisationService *OrganisationService
	clientID            string
	clientSecret        string
	redirectURL         string
}

func NewGithubIntegrationService(userService *UserService, jwtService *JWTService, organisationService *OrganisationService, clientID string, clientSecret string, redirectURL string) *GithubIntegrationService {
	return &GithubIntegrationService{
		userService:         userService,
		jwtService:          jwtService,
		organisationService: organisationService,
		clientID:            clientID,
		clientSecret:        clientSecret,
		redirectURL:         redirectURL,
	}
}
