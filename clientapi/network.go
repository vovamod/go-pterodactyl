package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/internal/requester"
)

type networkService struct {
	client           requester.Requester
	serverIdentifier string
}

// newNetworkService creates a new network service.
func newNetworkService(client requester.Requester, serverIdentifier string) *networkService {
	return &networkService{client: client, serverIdentifier: serverIdentifier}
}

// ListAllocations retrieves all network allocations for the server.
func (s *networkService) ListAllocations(ctx context.Context, options api.PaginationOptions) ([]*api.Allocation, *api.Meta, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/network/allocations", s.serverIdentifier)
	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil, &options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create list allocations request: %w", err)
	}

	res := &api.PaginatedResponse[api.Allocation]{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*api.Allocation, len(res.Data))
	for i, item := range res.Data {
		results[i] = item.Attributes
	}
	return results, &res.Meta, nil
}

// AssignAllocation requests that a new allocation be automatically assigned to the server.
func (s *networkService) AssignAllocation(ctx context.Context) (*api.Allocation, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/network/allocations", s.serverIdentifier)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create assign allocation request: %w", err)
	}

	res := &api.AllocationResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return res.Attributes, nil
}

// SetAllocationNote updates the notes for a specific allocation.
func (s *networkService) SetAllocationNote(ctx context.Context, allocationID int, options api.AllocationNoteOptions) (*api.Allocation, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal allocation note options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/client/servers/%s/network/allocations/%d", s.serverIdentifier, allocationID)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create set allocation note request: %w", err)
	}

	res := &api.AllocationResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return res.Attributes, nil
}

// SetPrimaryAllocation designates an allocation as the primary one for the server.
func (s *networkService) SetPrimaryAllocation(ctx context.Context, allocationID int) (*api.Allocation, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/network/allocations/%d/primary", s.serverIdentifier, allocationID)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create set primary allocation request: %w", err)
	}

	res := &api.AllocationResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return res.Attributes, nil
}

// UnassignAllocation removes a network allocation from the server.
func (s *networkService) UnassignAllocation(ctx context.Context, allocationID int) error {
	endpoint := fmt.Sprintf("/api/client/servers/%s/network/allocations/%d", s.serverIdentifier, allocationID)
	req, err := s.client.NewRequest(ctx, "DELETE", endpoint, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create unassign allocation request: %w", err)
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}
