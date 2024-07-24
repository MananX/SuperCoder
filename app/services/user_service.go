package services

import (
	"ai-developer/app/models"
	"ai-developer/app/repositories"
	"ai-developer/app/types/request"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

type UserService struct {
	userRepo             *repositories.UserRepository
	organisationUserRepo *repositories.OrganisationUserRepository
	orgService           *OrganisationService
	jwtService           *JWTService
}

func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	return s.userRepo.GetUserByID(userID)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetUserByEmail(email)
}

func (s *UserService) CreateUser(user *models.User) (*models.User, error) {
	return s.userRepo.CreateUser(user)
}

func (s *UserService) CreatePassword() string {
	length := 8
	b := make([]byte, length)
	for i := range b {
		switch rand.Intn(3) {
		case 0:
			b[i] = byte(rand.Intn(10)) + '0' // digits
		case 1:
			b[i] = byte(rand.Intn(26)) + 'A' // uppercase letters
		case 2:
			b[i] = byte(rand.Intn(26)) + 'a' // lowercase letters
		}
	}
	return string(b)
}

func (s *UserService) UpdateUserByEmail(email string, user *models.User) error {
	return s.userRepo.UpdateUserByEmail(email, user)
}

func (s *UserService) HandleUserSignUp(request request.CreateUserRequest, inviteToken string) (*models.User, string, error) {
	var err error
	var inviteOrganisationId *int
	var inviteEmail *string
	hashedPassword, err := s.HashUserPassword(request.Password)
	if err != nil {
		fmt.Println("Error while hashing password: ", err.Error())
		return nil, "", err
	}
	newUser := &models.User{
		Name:     request.Email,
		Email:    request.Email,
		Password: hashedPassword,
	}
	if inviteToken != "" {
		inviteEmail, inviteOrganisationId, err = s.jwtService.DecodeInviteToken(inviteToken)
		if err != nil {
			return nil, "", err
		}
		if inviteEmail != nil && inviteOrganisationId != nil {
			newUser, err = s.HandleUserInvite(newUser, inviteOrganisationId, inviteEmail, request.Email)
			if err != nil {
				return nil, "", err
			}
		}
	}
	if newUser.OrganisationID == 0 {
		organisation := &models.Organisation{
			Name: s.orgService.CreateOrganisationName(),
		}
		organisation, err = s.orgService.CreateOrganisation(organisation)
		newUser.OrganisationID = organisation.ID
	}
	newUser, err = s.CreateUser(newUser)
	if err != nil {
		fmt.Println("Error while creating user: ", err.Error())
		return nil, "", err
	}
	_, err = s.createOrganisationUser(newUser)
	var accessToken, jwtErr = s.jwtService.GenerateToken(int(newUser.ID), newUser.Email)
	if jwtErr != nil {
		fmt.Println(" Jwt error: ", accessToken, jwtErr.Error())
		return nil, "", nil
	}
	return newUser, accessToken, nil
}

func (s *UserService) HashUserPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *UserService) VerifyUserPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *UserService) HandleUserInvite(user *models.User, inviteOrgId *int, userEmail *string, primaryEmail string) (*models.User, error) {
	if *userEmail == primaryEmail {
		user.OrganisationID = uint(*inviteOrgId)
		_, err := s.createOrganisationUser(user)
		if err != nil {
			fmt.Println("Error while creating Organisation User: ", err.Error())
		}
	}
	return user, nil
}

func (s *UserService) createOrganisationUser(user *models.User) (*models.OrganisationUser, error) {
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

func NewUserService(userRepo *repositories.UserRepository, orgService *OrganisationService, jwtService *JWTService,
	organisationUserRepo *repositories.OrganisationUserRepository) *UserService {
	return &UserService{
		userRepo:             userRepo,
		orgService:           orgService,
		jwtService:           jwtService,
		organisationUserRepo: organisationUserRepo,
	}
}

func (s *UserService) FetchOrganisationIDByUserID(userID uint) (uint, error) {
	return s.userRepo.FetchOrganisationIDByUserID(userID)
}

func (s *UserService) GetDefaultUser() (*models.User, error) {
	defaultUser, err := s.GetUserByEmail("supercoder@superagi.com")
	if err != nil {
		return nil, err
	}
	return defaultUser, nil
}
