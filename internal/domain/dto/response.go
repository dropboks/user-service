package dto

import "errors"

var (
	SUCCESS_GET_PROFILE    = "success get profile data"
	SUCCESS_UPDATE_PROFILE = "success update profile data"
)

var (
	Err_INTERNAL_FAILED_BUILD_QUERY = errors.New("failed to build query")
	Err_INTERNAL_FAILED_SCAN_USER   = errors.New("failed to scan user")
	Err_INTERNAL_FAILED_INSERT_USER = errors.New("failed to insert user")
	Err_INTERNAL_FAILED_UPDATE_USER = errors.New("failed to update user")
	Err_INTERNAL_CONVERT_IMAGE      = errors.New("error processing image")

	Err_NOTFOUND_USER_NOT_FOUND = errors.New("user not found")

	Err_UNAUTHORIZED_USER_ID_NOTFOUND = errors.New("user_id is not found")

	Err_BAD_REQUEST_WRONG_EXTENTION     = errors.New("error file extension, support jpg, jpeg, and png")
	Err_BAD_REQUEST_LIMIT_SIZE_EXCEEDED = errors.New("max size exceeded: 6mb")
)

type (
	GetProfileResponse struct {
		FullName         string  `json:"full_name"`
		Image            *string `json:"image"`
		Email            string  `json:"email"`
		Verified         bool    `json:"verified"`
		TwoFactorEnabled bool    `json:"two_factor_enabled"`
	}
)
