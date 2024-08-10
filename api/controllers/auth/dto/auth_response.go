package dto

import authresult "github.com/beka-birhanu/app/user/auth/common"

type AuthResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
}

// FromAuthResult extracts the info for the login response from the given
// auth.Result and map them to new LoginResponse
func NewAuthResponse(authResult *authresult.Result) *AuthResponse {
	return &AuthResponse{
		ID:       authResult.ID.String(),
		Username: authResult.Username,
		IsAdmin:  authResult.IsAdmin,
	}
}
