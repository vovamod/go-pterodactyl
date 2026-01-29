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

type LocationService struct {
	client requester.Requester
}

func NewLocationService(client requester.Requester) *LocationService {
	return &LocationService{client: client}
}

func (s *LocationService) List(ctx context.Context, options *api.PaginationOptions) ([]*api.Location, *api.Meta, error) {
	return crud.List[api.Location](ctx, s.client, "/api/application/locations", options)
}

func (s *LocationService) ListAll(ctx context.Context) ([]*api.Location, error) {
	return crud.ListAll[api.Location](ctx, s.client, "/api/application/locations", 100)
}

func (s *LocationService) Get(ctx context.Context, id int) (*api.Location, error) {
	return crud.Get[api.Location](ctx, s.client, "/api/application/locations", id)
}

func (s *LocationService) Create(ctx context.Context, options api.LocationCreateOptions) (*api.Location, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create location Options: %w", err)
	}

	req, err := s.client.NewRequest(ctx, "POST", "/api/application/locations", bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new location request: %w", err)
	}

	response := &api.ListItem[api.Location]{}
	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, err
	}

	return response.Attributes, nil
}

func (s *LocationService) Update(ctx context.Context, id int, options api.LocationUpdateOptions) (*api.Location, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update location Options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/application/locations/%d", id)
	req, err := s.client.NewRequest(ctx, "PATCH", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create update location request: %w", err)
	}

	response := &api.ListItem[api.Location]{}
	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, err
	}

	return response.Attributes, nil
}

func (s *LocationService) Delete(ctx context.Context, id int) error {
	return crud.Delete[api.Location](ctx, s.client, "/api/application/locations", id)
}
