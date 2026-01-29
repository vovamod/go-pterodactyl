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

type nodesService struct {
	client requester.Requester
}

func NewNodesService(client requester.Requester) NodesService {
	return &nodesService{client: client}
}

func (s *nodesService) List(ctx context.Context, options *api.PaginationOptions) ([]*api.Node, *api.Meta, error) {
	return crud.List[api.Node](ctx, s.client, "/api/application/nodes", options)
}

func (s *nodesService) ListAll(ctx context.Context) ([]*api.Node, error) {
	return crud.ListAll[api.Node](ctx, s.client, "/api/application/nodes", 100)
}

func (s *nodesService) Get(ctx context.Context, id int) (*api.Node, error) {
	return crud.Get[api.Node](ctx, s.client, "/api/application/nodes", id)
}

func (s *nodesService) GetConfiguration(ctx context.Context, nodeID int) (*api.NodeConfiguration, error) {
	endpoint := fmt.Sprintf("/api/application/nodes/%d/configuration", nodeID)
	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get node configuration request: %w", err)
	}

	response := &api.NodeConfiguration{}
	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *nodesService) Create(ctx context.Context, options api.NodeCreateOptions) (*api.Node, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create node Options: %w", err)
	}

	req, err := s.client.NewRequest(ctx, "POST", "/api/application/nodes", bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new node request: %w", err)
	}

	response := &api.ListItem[api.Node]{}
	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, err
	}

	return response.Attributes, nil
}

func (s *nodesService) Update(ctx context.Context, id int, options api.NodeUpdateOptions) (*api.Node, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update node Options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/application/nodes/%d", id)
	req, err := s.client.NewRequest(ctx, "PATCH", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create update node request: %w", err)
	}

	response := &api.ListItem[api.Node]{}
	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, err
	}

	return response.Attributes, nil
}

func (s *nodesService) Delete(ctx context.Context, id int) error {
	return crud.Delete[api.Node](ctx, s.client, "/api/application/nodes", id)

}

func (s *nodesService) Allocations(ctx context.Context, nodeID int) AllocationsService {
	return newAllocationsService(s.client, nodeID)
}
