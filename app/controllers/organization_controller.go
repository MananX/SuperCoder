package controllers

import (
	"ai-developer/app/config"
	"ai-developer/app/models"
	"ai-developer/app/repositories"
	"ai-developer/app/services"
	"ai-developer/app/types/request"
	"ai-developer/app/types/response"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrganizationController struct {
	jwtService           *services.JWTService
	userService          *services.UserService
	organizationService  *services.OrganisationService
	organisationUserRepo *repositories.OrganisationUserRepository
	appRedirectUrl       string
}

func (controller *OrganizationController) FetchOrganizationUsers(c *gin.Context) {
	var users []*response.UsersResponse
	organizationID, err := controller.getOrganisationIDFromUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, &response.FetchOrganisationUserResponse{Success: false, Error: err.Error(), Users: nil})
	}
	fmt.Println("Fetching org users: ", organizationID)
	users, err = controller.organizationService.GetOrganizationUsers(organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.FetchOrganisationUserResponse{Success: false, Error: err.Error(), Users: nil})
		return
	}

	c.JSON(http.StatusOK, &response.FetchOrganisationUserResponse{Success: true, Error: nil, Users: users})
}

func (controller *OrganizationController) InviteUserToOrganisation(c *gin.Context) {
	var inviteUserRequest request.InviteUserRequest
	if err := c.ShouldBindJSON(&inviteUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	organizationID, err := controller.getOrganisationIDFromUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, &response.SendEmailResponse{
			Success:   false,
			MessageId: "",
			Error:     err.Error(),
		})
		return
	}
	sendEmailResponse, err := controller.organizationService.InviteUserToOrganization(int(organizationID), inviteUserRequest.Email, inviteUserRequest.CurrentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, &response.SendEmailResponse{
			Success:   false,
			MessageId: "",
			Error:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, sendEmailResponse)
}

func (controller *OrganizationController) HandleUserInvite(c *gin.Context) {
	var inviteToken = c.Query("invite_token")
	email, _, err := controller.jwtService.DecodeInviteToken(inviteToken)
	if err != nil {
		redirectUrl := controller.appRedirectUrl + "?error_msg=INVALID_TOKEN"
		c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		return
	}
	redirectUrl := controller.appRedirectUrl + "?user_email=" + email + "&invite_token=" + inviteToken
	c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

func (controller *OrganizationController) RemoveUserFromOrganisation(c *gin.Context) {
	var removeOrgUserRequest request.RemoveOrgUserRequest
	if err := c.ShouldBindJSON(&removeOrgUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	_, err := controller.getOrganisationIDFromUserID(c)
	if err != nil {
		c.JSON(http.StatusForbidden, &response.FetchOrganisationUserResponse{Success: false, Error: "OrganisationID mismatch", Users: nil})
		return
	}
	user, err := controller.userService.GetUserByID(uint(removeOrgUserRequest.UserID))
	if user == nil {
		c.JSON(http.StatusBadRequest, &response.FetchOrganisationUserResponse{Success: false, Error: "User not found"})
		return
	}
	organisation := &models.Organisation{
		Name: controller.organizationService.CreateOrganisationName(),
	}
	organisation, err = controller.organizationService.CreateOrganisation(organisation)
	if err != nil {
		fmt.Println("Error while creating organization: ", err.Error())
		c.JSON(http.StatusInternalServerError, &response.FetchOrganisationUserResponse{Success: false, Error: err.Error()})
		return
	}
	user.OrganisationID = organisation.ID
	err = controller.userService.UpdateUserByEmail(user.Email, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.FetchOrganisationUserResponse{Success: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, &response.FetchOrganisationUserResponse{Success: true, Error: nil})
}

func (controller *OrganizationController) getOrganisationIDFromUserID(context *gin.Context) (uint, error) {
	userID, exists := context.Get("user_id")
	if !exists {
		return 0, errors.New("userId not found in context")
	}
	userIDInt, ok := userID.(int)
	if !ok {
		context.JSON(http.StatusBadRequest, gin.H{"error": "User ID is not of type int"})
		return 0, errors.New("userId is not of type int")
	}
	organisationIdByUserID, err := controller.userService.FetchOrganisationIDByUserID(uint(userIDInt))
	if err != nil {
		return 0, err
	}
	return organisationIdByUserID, nil
}

func NewOrganizationController(
	jwtService *services.JWTService,
	userService *services.UserService,
	organizationService *services.OrganisationService,
	organisationUserRepo *repositories.OrganisationUserRepository,
) *OrganizationController {
	return &OrganizationController{
		jwtService:           jwtService,
		userService:          userService,
		organizationService:  organizationService,
		appRedirectUrl:       config.AppUrl(),
		organisationUserRepo: organisationUserRepo,
	}
}
