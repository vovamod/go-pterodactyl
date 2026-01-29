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

type databaseService struct {
	client   requester.Requester
	serverID int
}

func newDatabaseService(client requester.Requester, serverID int) *databaseService {
	return &databaseService{client: client, serverID: serverID}
}

func (s *databaseService) List(ctx context.Context, options api.PaginationOptions) ([]*api.Database, *api.Meta, error) {
	endpoint := fmt.Sprintf("/api/application/servers/%d/databases", s.serverID)
	return crud.List[api.Database](ctx, s.client, endpoint, &options)
}

// Get fetches a single database by its ID for the configured server.
func (s *databaseService) Get(ctx context.Context, databaseID int) (*api.Database, error) {
	endpoint := fmt.Sprintf("/api/application/servers/%d/databases", s.serverID)
	return crud.Get[api.Database](ctx, s.client, endpoint, databaseID)
}

// Create creates a new database for the configured server.
func (s *databaseService) Create(ctx context.Context, options api.DatabaseCreateOptions) (*api.Database, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create database Options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/application/servers/%d/databases", s.serverID)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new database request: %w", err)
	}

	response := &api.ListItem[api.Database]{}
	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, err
	}
	return response.Attributes, nil
}

// ResetPassword Requests a password reset for a specific database.
func (s *databaseService) ResetPassword(ctx context.Context, databaseID int) error {
	endpoint := fmt.Sprintf("/api/application/servers/%d/databases/%d/reset-password", s.serverID, databaseID)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create reset database password request: %w", err)
	}

	// This Endpoint returns 204 No Content.
	_, err = s.client.Do(ctx, req, nil)
	return err
}

// Delete deletes a specific database from the configured server.
func (s *databaseService) Delete(ctx context.Context, databaseID int) error {
	endpoint := fmt.Sprintf("/api/application/servers/%d/databases", s.serverID)
	return crud.Delete[api.Database](ctx, s.client, endpoint, databaseID)
}
