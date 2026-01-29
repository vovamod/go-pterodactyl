package appapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/internal/crud"
	"github.com/vovamod/go-pterodactyl/internal/requester"
)

type usersService struct {
	client requester.Requester
}

func NewUsersService(client requester.Requester) *usersService {
	return &usersService{client: client}
}

func (s *usersService) List(ctx context.Context, options *api.PaginationOptions) ([]*api.User, *api.Meta, error) {
	return crud.List[api.User](ctx, s.client, "/api/application/users", options)
}

func (s *usersService) ListAll(ctx context.Context) ([]*api.User, error) {
	return crud.ListAll[api.User](ctx, s.client, "/api/application/users", 100)
}

func (s *usersService) Get(ctx context.Context, id int) (*api.User, error) {
	return crud.Get[api.User](ctx, s.client, "/api/application/users", id)
}

func (s *usersService) GetExternalID(ctx context.Context, externalId string) (*api.User, error) {
	endpoint := fmt.Sprintf("/api/application/users/external/%s", externalId)

	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user get request: %w", err)
	}
	response := &api.ListItem[api.User]{}

	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, err
	}

	return response.Attributes, nil

}

func (s *usersService) Create(ctx context.Context, options api.UserCreateOptions) (*api.User, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create user Options: %w", err)
	}

	req, err := s.client.NewRequest(ctx, "POST", "/api/application/users", bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new user request: %w", err)
	}

	response := &api.ListItem[api.User]{}
	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, err
	}
	return response.Attributes, nil
}

func (s *usersService) Update(ctx context.Context, id int, options api.UserUpdateOptions) (*api.User, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update user Options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/application/users/%d", id)
	req, err := s.client.NewRequest(ctx, "PATCH", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create update user request: %w", err)
	}

	response := &api.ListItem[api.User]{}
	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, err
	}

	return response.Attributes, nil
}

func (s *usersService) Delete(ctx context.Context, id int) error {
	return crud.Delete[api.User](ctx, s.client, "/api/application/users", id)
}
