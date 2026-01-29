package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/internal/requester"
)

type startupService struct {
	client           requester.Requester
	serverIdentifier string
}

// newStartupService creates a new startup service.
func newStartupService(client requester.Requester, serverIdentifier string) *startupService {
	return &startupService{client: client, serverIdentifier: serverIdentifier}
}

// ListVariables retrieves all startup variables for the server.
func (s *startupService) ListVariables(ctx context.Context, options api.PaginationOptions) ([]*api.StartupVariable, *api.Meta, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/startup", s.serverIdentifier)
	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil, &options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create list startup variables request: %w", err)
	}

	// The response is a standard paginated list of objects.
	res := &api.PaginatedResponse[api.StartupVariable]{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, nil, err
	}

	// Extract the attributes from each item in the response data.
	results := make([]*api.StartupVariable, len(res.Data))
	for i, item := range res.Data {
		results[i] = item.Attributes
	}

	return results, &res.Meta, nil
}

// UpdateVariable updates the value of a single startup variable.
func (s *startupService) UpdateVariable(ctx context.Context, options api.UpdateVariableOptions) (*api.StartupVariable, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update variable options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/client/servers/%s/startup/variable", s.serverIdentifier)
	req, err := s.client.NewRequest(ctx, "PUT", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create update startup variable request: %w", err)
	}

	// The API returns the updated variable, wrapped in a standard object with attributes.
	res := &api.UpdateVariableResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return res.Attributes, nil
}
