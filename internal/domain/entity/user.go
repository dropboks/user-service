package entity

type User struct {
	ID               string
	FullName         string
	Image            *string
	Email            string
	Password         string
	Verified         bool
	TwoFactorEnabled bool
}
