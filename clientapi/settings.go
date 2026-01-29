package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/internal/requester"
)

type settingsService struct {
	client           requester.Requester
	serverIdentifier string
}

func newSettingsService(client requester.Requester, serverIdentifier string) *settingsService {
	return &settingsService{client: client, serverIdentifier: serverIdentifier}
}

// Rename sends a request to change the server's name and optional description.
// A successful request returns a 204 No Content response.
func (s *settingsService) Rename(ctx context.Context, options api.RenameOptions) error {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return fmt.Errorf("failed to marshal rename options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/client/servers/%s/settings/rename", s.serverIdentifier)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return fmt.Errorf("failed to create rename request: %w", err)
	}

	// This endpoint returns 204 No Content on success, so we don't need to decode a response body.
	_, err = s.client.Do(ctx, req, nil)
	return err
}

// Reinstall sends a request to reinstall the server.
// A successful request returns a 204 No Content response.
func (s *settingsService) Reinstall(ctx context.Context) error {
	endpoint := fmt.Sprintf("/api/client/servers/%s/settings/reinstall", s.serverIdentifier)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create reinstall request: %w", err)
	}

	// This endpoint also returns 204 No Content on success.
	_, err = s.client.Do(ctx, req, nil)
	return err
}
