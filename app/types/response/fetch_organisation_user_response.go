package response

type FetchOrganisationUserResponse struct {
	Success bool             `json:"success"`
	Error   interface{}      `json:"error,omitempty"`
	Users   []*UsersResponse `json:"users,omitempty"`
}

type UsersResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	OrganisationID uint   `json:"organisation_id"`
}
