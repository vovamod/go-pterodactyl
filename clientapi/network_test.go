package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/errors"
	"github.com/vovamod/go-pterodactyl/internal/testutil"
)

func TestNetworkService_ListAllocations(t *testing.T) {
	alias := "test-alias"
	notes := "test-notes"
	expectedAllocs := []*api.Allocation{
		{
			ID:       1,
			IP:       "127.0.0.1",
			Alias:    &alias,
			Port:     8080,
			Notes:    &notes,
			Assigned: true,
		},
		{
			ID:       2,
			IP:       "127.0.0.1",
			Port:     8081,
			Assigned: false,
		},
	}
	data := make([]*api.ListItem[api.Allocation], len(expectedAllocs))
	for i, a := range expectedAllocs {
		data[i] = &api.ListItem[api.Allocation]{Object: "allocation", Attributes: a}
	}
	meta := api.Meta{Pagination: api.Pagination{Total: 2, PerPage: 25, CurrentPage: 1, TotalPages: 1}}
	res := api.PaginatedResponse[api.Allocation]{
		Object: "list",
		Data:   data,
		Meta:   meta,
	}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newNetworkService(mock, testServerIdentifier)
		allocs, m, err := s.ListAllocations(context.Background(), api.PaginationOptions{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(allocs, expectedAllocs) {
			t.Errorf("expected allocations %+v, got %+v", expectedAllocs, allocs)
		}
		if !reflect.DeepEqual(m, &meta) {
			t.Errorf("expected meta %+v, got %+v", &meta, m)
		}
	})
}

func TestNetworkService_AssignAllocation(t *testing.T) {
	expectedAlloc := &api.Allocation{
		ID:       3,
		IP:       "127.0.0.1",
		Port:     8082,
		Assigned: true,
	}
	res := api.AllocationResponse{Object: "allocation", Attributes: expectedAlloc}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newNetworkService(mock, testServerIdentifier)
		alloc, err := s.AssignAllocation(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(alloc, expectedAlloc) {
			t.Errorf("expected allocation %+v, got %+v", expectedAlloc, alloc)
		}
		req := mock.Requests[0]
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/network/allocations", testServerIdentifier)
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
		if req.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", req.Method)
		}
	})
}

func TestNetworkService_SetAllocationNote(t *testing.T) {
	allocID := 4
	notes := "new-notes"
	options := api.AllocationNoteOptions{Notes: &notes}
	jsonOptions, _ := json.Marshal(options)
	expectedAlloc := &api.Allocation{
		ID:    allocID,
		IP:    "127.0.0.1",
		Port:  8083,
		Notes: &notes,
	}
	res := api.AllocationResponse{Object: "allocation", Attributes: expectedAlloc}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newNetworkService(mock, testServerIdentifier)
		alloc, err := s.SetAllocationNote(context.Background(), allocID, options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(alloc, expectedAlloc) {
			t.Errorf("expected allocation %+v, got %+v", expectedAlloc, alloc)
		}
		req := mock.Requests[0]
		if !bytes.Equal(req.Body, jsonOptions) {
			t.Errorf("expected body %s, got %s", jsonOptions, req.Body)
		}
	})
}

func TestNetworkService_SetPrimaryAllocation(t *testing.T) {
	allocID := 5
	expectedAlloc := &api.Allocation{ID: allocID, IP: "127.0.0.1", Port: 8084}
	res := api.AllocationResponse{Object: "allocation", Attributes: expectedAlloc}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newNetworkService(mock, testServerIdentifier)
		alloc, err := s.SetPrimaryAllocation(context.Background(), allocID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(alloc, expectedAlloc) {
			t.Errorf("expected allocation %+v, got %+v", expectedAlloc, alloc)
		}
		req := mock.Requests[0]
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/network/allocations/%d/primary", testServerIdentifier, allocID)
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
	})
}

func TestNetworkService_UnassignAllocation(t *testing.T) {
	allocID := 6
	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusNoContent}},
		}
		s := newNetworkService(mock, testServerIdentifier)
		err := s.UnassignAllocation(context.Background(), allocID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		req := mock.Requests[0]
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/network/allocations/%d", testServerIdentifier, allocID)
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
		if req.Method != http.MethodDelete {
			t.Errorf("expected method DELETE, got %s", req.Method)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{Err: &errors.APIError{HTTPStatusCode: http.StatusNotFound}}},
		}
		s := newNetworkService(mock, testServerIdentifier)
		err := s.UnassignAllocation(context.Background(), allocID)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}
