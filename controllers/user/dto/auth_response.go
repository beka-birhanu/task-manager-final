package dto

import usersvc "github.com/beka-birhanu/service/user"

type AuthResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// FromAuthResult extracts the info for the login response from the given
// auth.Result and map them to new LoginResponse
func NewAuthResponse(authResult *usersvc.AuthResult) *AuthResponse {
	return &AuthResponse{ID: authResult.ID.String(), Username: authResult.Username}
}
