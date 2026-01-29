package api

import "time"

type User struct {
	ID         int       `json:"id"`
	ExternalID *string   `json:"external_id"`
	UUID       string    `json:"uuid"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Language   string    `json:"language"`
	RootAdmin  bool      `json:"root_admin"`
	TwoFA      bool      `json:"2fa"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserCreateOptions struct {
	Email      string `json:"email"`
	Username   string `json:"username"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Password   string `json:"password,omitempty"`
	RootAdmin  bool   `json:"root_admin,omitempty"`
	ExternalID string `json:"external_id,omitempty"`
	Language   string `json:"language,omitempty"`
}

type UserUpdateOptions struct {
	Email      string `json:"email,omitempty"`
	Username   string `json:"username,omitempty"`
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	Password   string `json:"password,omitempty"`
	RootAdmin  bool   `json:"root_admin,omitempty"`
	ExternalID string `json:"external_id,omitempty"`
	Language   string `json:"language,omitempty"`
}
