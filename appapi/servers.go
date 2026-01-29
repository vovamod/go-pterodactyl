package appapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/internal/crud"
	"github.com/vovamod/go-pterodactyl/internal/requester"
	"io"
	"net/url"
)

type serversService struct {
	client requester.Requester
}

// NewServersService is the exported constructor.
func NewServersService(client requester.Requester) ServersService {
	return &serversService{client: client}
}

func (s *serversService) List(ctx context.Context, options api.PaginationOptions) ([]*api.Server, *api.Meta, error) {
	return crud.List[api.Server](ctx, s.client, "/api/application/servers", &options)
}

func (s *serversService) ListAll(ctx context.Context) ([]*api.Server, error) {
	return crud.ListAll[api.Server](ctx, s.client, "/api/application/servers", 100)
}

func (s *serversService) Get(ctx context.Context, id int) (*api.Server, error) {
	endpoint := "/api/application/servers"
	return crud.Get[api.Server](ctx, s.client, endpoint, id)
}

func (s *serversService) GetExternal(ctx context.Context, externalID string) (*api.Server, error) {
	endpoint := fmt.Sprintf("/api/application/servers/external/%s", url.PathEscape(externalID))
	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get external server request: %w", err)
	}

	response := &api.ListItem[api.Server]{}
	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, err
	}
	return response.Attributes, nil
}

func (s *serversService) Create(ctx context.Context, options api.ServerCreateOptions) (*api.Server, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create server Options: %w", err)
	}
	req, err := s.client.NewRequest(ctx, "POST", "/api/application/servers", bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, err
	}
	response := &api.ListItem[api.Server]{}
	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, err
	}
	return response.Attributes, nil
}

func (s *serversService) UpdateDetails(ctx context.Context, serverID int, options api.ServerUpdateDetailsOptions) (*api.Server, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("/api/application/servers/%d/details", serverID)
	req, err := s.client.NewRequest(ctx, "PATCH", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, err
	}
	response := &api.ListItem[api.Server]{}
	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, err
	}
	return response.Attributes, nil
}

func (s *serversService) UpdateBuild(ctx context.Context, serverID int, options api.ServerUpdateBuildOptions) (*api.Server, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("/api/application/servers/%d/build", serverID)
	req, err := s.client.NewRequest(ctx, "PATCH", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, err
	}
	response := &api.ListItem[api.Server]{}
	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, err
	}
	return response.Attributes, nil
}

func (s *serversService) UpdateStartup(ctx context.Context, serverID int, options api.ServerUpdateStartupOptions) (*api.Server, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("/api/application/servers/%d/startup", serverID)
	req, err := s.client.NewRequest(ctx, "PATCH", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, err
	}
	response := &api.ListItem[api.Server]{}
	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, err
	}
	return response.Attributes, nil
}

func (s *serversService) Suspend(ctx context.Context, serverID int) error {
	endpoint := fmt.Sprintf("/api/application/servers/%d/suspend", serverID)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, nil, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(ctx, req, nil)
	return err
}

func (s *serversService) Unsuspend(ctx context.Context, serverID int) error {
	endpoint := fmt.Sprintf("/api/application/servers/%d/unsuspend", serverID)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, nil, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(ctx, req, nil)
	return err
}

func (s *serversService) Reinstall(ctx context.Context, serverID int) error {
	endpoint := fmt.Sprintf("/api/application/servers/%d/reinstall", serverID)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, nil, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(ctx, req, nil)
	return err
}

func (s *serversService) Delete(ctx context.Context, serverID int, force bool) error {
	var body io.Reader
	if force {
		jsonBytes, err := json.Marshal(api.ServerDeleteOptions{Force: true})
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(jsonBytes)
	}
	endpoint := fmt.Sprintf("/api/application/servers/%d", serverID)
	req, err := s.client.NewRequest(ctx, "DELETE", endpoint, body, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(ctx, req, nil)
	return err
}

// Databases returns a specialized service for managing databases for a specific server.
func (s *serversService) Databases(ctx context.Context, serverID int) DatabaseService {
	return newDatabaseService(s.client, serverID)
}
