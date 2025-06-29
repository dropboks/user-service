package dto

import "mime/multipart"

type (
	UserData struct {
		UserId string `json:"user_id"`
	}

	UpdateUserRequest struct {
		FullName         string                `form:"full_name"`
		Image            *multipart.FileHeader `form:"image"`
		TwoFactorEnabled bool                  `form:"two_factor_enabled"`
	}
)
