package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/internal/requester"
)

type usersService struct {
	client           requester.Requester
	serverIdentifier string
}

func newUsersService(client requester.Requester, serverIdentifier string) *usersService {
	return &usersService{client: client, serverIdentifier: serverIdentifier}
}

func (s *usersService) List(ctx context.Context, options api.PaginationOptions) ([]*api.Subuser, *api.Meta, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/users", s.serverIdentifier)
	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil, &options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create list subusers request: %w", err)
	}

	res := &api.PaginatedResponse[api.Subuser]{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*api.Subuser, len(res.Data))
	for i, item := range res.Data {
		results[i] = item.Attributes
	}
	return results, &res.Meta, nil
}

// Create sends a request to add a new subuser to the server.
func (s *usersService) Create(ctx context.Context, options api.SubuserCreateOptions) (*api.Subuser, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create subuser options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/client/servers/%s/users", s.serverIdentifier)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new subuser request: %w", err)
	}

	res := &api.SubuserResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return res.Attributes, nil
}

// Details retrieves the details of a specific subuser by their UUID.
func (s *usersService) Details(ctx context.Context, uuid string) (*api.Subuser, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/users/%s", s.serverIdentifier, uuid)
	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create subuser details request: %w", err)
	}

	res := &api.SubuserResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return res.Attributes, nil
}

// Update modifies the permissions of an existing subuser.
// Note: The Pterodactyl API uses POST for this update operation.
func (s *usersService) Update(ctx context.Context, uuid string, options api.SubuserUpdateOptions) (*api.Subuser, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update subuser options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/client/servers/%s/users/%s", s.serverIdentifier, uuid)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create update subuser request: %w", err)
	}

	res := &api.SubuserResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return res.Attributes, nil
}

// Delete removes a subuser from the server.
func (s *usersService) Delete(ctx context.Context, uuid string) error {
	endpoint := fmt.Sprintf("/api/client/servers/%s/users/%s", s.serverIdentifier, uuid)
	req, err := s.client.NewRequest(ctx, "DELETE", endpoint, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete subuser request: %w", err)
	}

	// A successful deletion returns a 204 No Content response.
	_, err = s.client.Do(ctx, req, nil)
	return err
}
