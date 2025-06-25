package dto

import "errors"

var (
	Err_INTERNAL_FAILED_BUILD_QUERY = errors.New("failed to build query")
	Err_INTERNAL_FAILED_SCAN_USER   = errors.New("failed to scan user")
	Err_INTERNAL_FAILED_INSERT_USER = errors.New("failed to insert user")
	Err_INTERNAL_FAILED_UPDATE_USER = errors.New("failed to update user")

	Err_NOTFOUND_USER_NOT_FOUND = errors.New("user not found")

	Err_UNAUTHORIZED_USER_ID_NOTFOUND = errors.New("user_id is not found")
)

type (
	GetProfileResponse struct {
		FullName         string `json:"full_name"`
		Image            string `json:"image"`
		Email            string `json:"email"`
		Verified         bool   `json:"verified"`
		TwoFactorEnabled bool   `json:"two_factor_enabled"`
	}
)
