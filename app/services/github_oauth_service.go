package services

import (
	"ai-developer/app/models"
	"ai-developer/app/repositories"
	"context"
	"errors"
	"fmt"
	goGithub "github.com/google/go-github/github"
	"golang.org/x/oauth2"
	oauthGithub "golang.org/x/oauth2/github"
	"strings"
	"time"
)

type GithubOauthService struct {
	jwtService           *JWTService
	userService          *UserService
	organisationService  *OrganisationService
	organisationUserRepo *repositories.OrganisationUserRepository
	clientID             string
	clientSecret         string
	redirectURL          string
}

func (s *GithubOauthService) HandleGithubCallback(code string, state string) (string, string, string, string, int, error) {
	var githubOauthConfig = &oauth2.Config{
		ClientID:     s.clientID,
		ClientSecret: s.clientSecret,
		RedirectURL:  s.redirectURL,
		Scopes:       []string{"user:email"},
		Endpoint:     oauthGithub.Endpoint,
	}
	token, err := githubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return "", "", "", "", 0, err
	}

	client := goGithub.NewClient(githubOauthConfig.Client(context.Background(), token))
	githubUser, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		return "", "", "", "", 0, err
	}
	emails, _, err := client.Users.ListEmails(context.Background(), nil)
	if err != nil {
		return "", "", "", "", 0, err
	}
	var primaryEmail string
	for _, email := range emails {
		if email.GetPrimary() {
			primaryEmail = email.GetEmail()
			break
		}
	}
	var name string
	if githubUser.Login != nil {
		name = *githubUser.Login
	} else {
		name = "N/A"
	}
	var userExists string
	user, err := s.userService.GetUserByEmail(primaryEmail)
	if err != nil {
		if user == nil {
			hashedPassword, err := s.userService.HashUserPassword(s.userService.CreatePassword())
			if err != nil {
				return "", "", "", "", 0, errors.New("error while hashing user password")
			}
			user = &models.User{
				Name:     name,
				Email:    primaryEmail,
				Password: hashedPassword,
			}
			if err != nil {
				return "", "", "", "", 0, err
			}
			user, err = s.userService.CreateUser(user)
			if err != nil {
				return "", "", "", "", 0, err
			}
			_, err = s.createOrganisationUser(user)
			if err != nil {
				return "", "", "", "", 0, err
			}
			userExists = "no"
		}
	}
	user.Name = name
	userEmail, inviteOrgId, err := s.DecodeInviteToken(state)
	if err != nil {
		return "", "", "", "", 0, err
	}
	if userEmail != nil && inviteOrgId != nil {
		user, err = s.handleUserInvite(user, inviteOrgId, userEmail, primaryEmail)
	}
	if user.OrganisationID == 0 {
		organisation := &models.Organisation{
			Name: s.organisationService.CreateOrganisationName(),
		}
		organisation, err = s.organisationService.CreateOrganisation(organisation)
		user.OrganisationID = organisation.ID
	}

	err = s.userService.UpdateUserByEmail(primaryEmail, user)
	if err != nil {
		return "", "", "", "", 0, err
	}
	userExists = "yes"
	accessToken, err := s.jwtService.GenerateToken(int(user.ID), primaryEmail)
	if err != nil {
		return "", "", "", "", 0, err
	}

	return accessToken, name, primaryEmail, userExists, int(user.OrganisationID), nil
}

func (s *GithubOauthService) handleUserInvite(user *models.User, inviteOrgId *int, userEmail *string, primaryEmail string) (*models.User, error) {
	if userEmail != nil && *userEmail == primaryEmail {
		user.OrganisationID = uint(*inviteOrgId)
		_, err := s.createOrganisationUser(user)
		if err != nil {
			fmt.Println("Error while creating Organisation User: ", err.Error())
		}
	}
	return user, nil
}

func (s *GithubOauthService) createOrganisationUser(user *models.User) (*models.OrganisationUser, error) {
	orgUser, err := s.organisationUserRepo.GetOrganisationUserByUserIDAndOrganisationID(user.ID, user.OrganisationID)
	if err != nil {
		return nil, err
	}
	if orgUser == nil {
		return s.organisationUserRepo.CreateOrganisationUser(s.organisationUserRepo.GetDB(), &models.OrganisationUser{
			OrganisationID: user.OrganisationID,
			UserID:         user.ID,
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		})
	}
	return orgUser, nil
}

func (s *GithubOauthService) DecodeInviteToken(state string) (*string, *int, error) {
	if strings.HasPrefix(state, "token:") {
		tokenValue := strings.TrimPrefix(state, "token:")
		userEmail, inviteOrgId, err := s.jwtService.DecodeInviteToken(tokenValue)
		if err != nil {
			return nil, nil, err
		}
		return userEmail, inviteOrgId, nil
	}
	return nil, nil, nil
}

func NewGithubOauthService(userService *UserService, jwtService *JWTService, organisationService *OrganisationService, organisationUserRepo *repositories.OrganisationUserRepository, clientID string, clientSecret string, redirectURL string) *GithubOauthService {
	return &GithubOauthService{
		userService:          userService,
		jwtService:           jwtService,
		organisationService:  organisationService,
		organisationUserRepo: organisationUserRepo,
		clientID:             clientID,
		clientSecret:         clientSecret,
		redirectURL:          redirectURL,
	}
}
