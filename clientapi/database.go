package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/internal/crud"
	"github.com/vovamod/go-pterodactyl/internal/requester"
)

type databasesService struct {
	client           requester.Requester
	serverIdentifier string
}

func newDatabasesService(client requester.Requester, serverIdentifier string) *databasesService {
	return &databasesService{client: client, serverIdentifier: serverIdentifier}
}

func (s *databasesService) List(ctx context.Context, options api.PaginationOptions) ([]*api.ClientDatabase, *api.Meta, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/databases", s.serverIdentifier)
	return crud.List[api.ClientDatabase](ctx, s.client, endpoint, &options)
}

func processCreateOrRotateResponse(res *api.ClientDatabaseCreateResponse) *api.ClientDatabase {
	db := res.Attributes
	if res.Relationships != nil && res.Relationships.Password != nil && res.Relationships.Password.Attributes != nil {
		db.Password = res.Relationships.Password.Attributes.Password
	}
	return db
}

func (s *databasesService) Create(ctx context.Context, options api.ClientDatabaseCreateOptions) (*api.ClientDatabase, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create database options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/client/servers/%s/databases", s.serverIdentifier)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new database request: %w", err)
	}

	// Use the special response struct to capture the password from relationships.
	res := &api.ClientDatabaseCreateResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return processCreateOrRotateResponse(res), nil
}

func (s *databasesService) RotatePassword(ctx context.Context, databaseID string) (*api.ClientDatabase, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/databases/%s/rotate-password", s.serverIdentifier, databaseID)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create rotate password request: %w", err)
	}

	res := &api.ClientDatabaseCreateResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return processCreateOrRotateResponse(res), nil
}

func (s *databasesService) Delete(ctx context.Context, databaseID string) error {
	endpoint := fmt.Sprintf("/api/client/servers/%s/databases/%s", s.serverIdentifier, databaseID)
	req, err := s.client.NewRequest(ctx, "DELETE", endpoint, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete database request: %w", err)
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}
