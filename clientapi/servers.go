package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/internal/requester"
)

type serverService struct {
	client     requester.Requester
	identifier string
}

// newServerService is the internal constructor.
func newServerService(client requester.Requester, identifier string) *serverService {
	return &serverService{client: client, identifier: identifier}
}

// GetDetails fetches the details for the configured server.
func (s *serverService) GetDetails(ctx context.Context) (*api.ClientServer, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s", s.identifier)
	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	res := &api.ListItem[api.ClientServer]{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}
	return res.Attributes, nil
}

// GetWebsocket fetches the details needed to connect to the server's console.
func (s *serverService) GetWebsocket(ctx context.Context) (*api.WebsocketDetails, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/websocket", s.identifier)
	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	res := &api.WebsocketResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}
	return &res.Data, nil
}

// GetResources fetches the current resource usage for the server.
func (s *serverService) GetResources(ctx context.Context) (*api.Resources, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/resources", s.identifier)
	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	// This endpoint returns the object directly, not wrapped in a ListItem
	res := &api.Resources{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SendCommand sends a command to the server's console.
func (s *serverService) SendCommand(ctx context.Context, command string) error {
	opts := api.SendCommandOptions{Command: command}
	jsonBytes, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("/api/client/servers/%s/command", s.identifier)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil) // Expects 204 No Content
	return err
}

// SetPowerState changes the power state of the server.
func (s *serverService) SetPowerState(ctx context.Context, signal string) error {
	opts := api.SetPowerStateOptions{Signal: signal}
	jsonBytes, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("/api/client/servers/%s/power", s.identifier)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil) // Expects 204 No Content
	return err
}

func (s *serverService) Databases() DatabasesService {
	return newDatabasesService(s.client, s.identifier)
}

func (s *serverService) Files() FileService {
	return newFilesService(s.client, s.identifier)
}

func (s *serverService) Schedules() ScheduleService {
	return newSchedulesService(s.client, s.identifier)
}

func (s *serverService) Network() NetworkService {
	return newNetworkService(s.client, s.identifier)
}

func (s *serverService) Users() UsersService {
	return newUsersService(s.client, s.identifier)
}

func (s *serverService) Backups() BackupService {
	return newBackupsService(s.client, s.identifier)
}

func (s *serverService) Startup() StartupService {
	return newStartupService(s.client, s.identifier)
}

func (s *serverService) Settings() SettingsService {
	return newSettingsService(s.client, s.identifier)
}
