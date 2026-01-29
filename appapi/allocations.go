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

type allocationsService struct {
	client requester.Requester
	nodeID int
}

func newAllocationsService(client requester.Requester, nodeID int) *allocationsService {
	return &allocationsService{client: client, nodeID: nodeID}
}

func (s *allocationsService) List(ctx context.Context, options *api.PaginationOptions) ([]*api.Allocation, *api.Meta, error) {
	endpoint := fmt.Sprintf("/api/application/nodes/%d/allocations", s.nodeID)
	return crud.List[api.Allocation](ctx, s.client, endpoint, options)
}

func (s *allocationsService) ListAll(ctx context.Context) ([]*api.Allocation, error) {
	endpoint := fmt.Sprintf("/api/application/nodes/%d/allocations", s.nodeID)
	return crud.ListAll[api.Allocation](ctx, s.client, endpoint, 100)
}

func (s *allocationsService) Create(ctx context.Context, options api.AllocationCreateOptions) error {
	// Marshal the Options struct into JSON.
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return fmt.Errorf("failed to marshal create allocation Options: %w", err)
	}

	// Construct the Endpoint for creating allocations on this node.
	endpoint := fmt.Sprintf("/api/application/nodes/%d/allocations", s.nodeID)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return fmt.Errorf("failed to create new allocation request: %w", err)
	}

	// Execute the request. We pass `nil` for the decoding target because we
	// expect a 204 No Content response.
	_, err = s.client.Do(ctx, req, nil)
	return err
}

// Delete deletes a specific allocation from the configured node.
func (s *allocationsService) Delete(ctx context.Context, allocationID int) error {
	// Construct the specific Endpoint for the allocation to be deleted.
	endpoint := fmt.Sprintf("/api/application/nodes/%d/allocations", s.nodeID)
	return crud.Delete[api.Allocation](ctx, s.client, endpoint, allocationID)
}
