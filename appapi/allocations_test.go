package appapi

import (
	"context"
	"fmt"
	"github.com/vovamod/go-pterodactyl/internal/testutil"
	"strings"
	"testing"

	"github.com/vovamod/go-pterodactyl/api"
)

func TestAllocationsService_List(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		nodeID         int
		options        *api.PaginationOptions
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedCount  int
		expectedMethod string
		expectedPath   string
	}{
		{
			name:   "Successful list with pagination",
			nodeID: 1,
			options: &api.PaginationOptions{
				Page:    1,
				PerPage: 10,
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body: []byte(`{
					"object": "list",
					"data": [
						{
							"object": "allocation",
							"attributes": {
								"id": 1,
								"ip": "192.168.1.1",
								"port": 25565,
								"assigned": false
							}
						},
						{
							"object": "allocation",
							"attributes": {
								"id": 2,
								"ip": "192.168.1.1",
								"port": 25566,
								"assigned": true
							}
						}
					],
					"meta": {
						"pagination": {
							"total": 2,
							"count": 2,
							"per_page": 10,
							"current_page": 1,
							"total_pages": 1
						}
					}
				}`),
			},
			expectedError:  false,
			expectedCount:  2,
			expectedMethod: "GET",
			expectedPath:   "/api/application/nodes/1/allocations",
		},
		{
			name:    "Successful list without pagination",
			nodeID:  2,
			options: nil,
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body: []byte(`{
					"object": "list",
					"data": [],
					"meta": {
						"pagination": {
							"total": 0,
							"count": 0,
							"per_page": 100,
							"current_page": 1,
							"total_pages": 0
						}
					}
				}`),
			},
			expectedError:  false,
			expectedCount:  0,
			expectedMethod: "GET",
			expectedPath:   "/api/application/nodes/2/allocations",
		},
		{
			name:    "API error response",
			nodeID:  3,
			options: &api.PaginationOptions{Page: 1},
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body: []byte(`{
					"errors": [{
						"code": "NotFoundHttpException",
						"status": "404",
						"detail": "The requested node could not be located."
					}]
				}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/nodes/3/allocations",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{
				Responses: []testutil.MockResponse{tc.mockResponse},
			}

			service := newAllocationsService(mock, tc.nodeID)

			allocations, meta, err := service.List(context.Background(), tc.options)

			// Check error expectations
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check response expectations
			if len(allocations) != tc.expectedCount {
				t.Errorf("expected %d allocations, got %d", tc.expectedCount, len(allocations))
			}

			// Check request expectations
			if len(mock.Requests) != 1 {
				t.Fatalf("expected 1 request, got %d", len(mock.Requests))
			}

			req := mock.Requests[0]
			if req.Method != tc.expectedMethod {
				t.Errorf("expected Method %s, got %s", tc.expectedMethod, req.Method)
			}

			if req.Endpoint != tc.expectedPath {
				t.Errorf("expected path %s, got %s", tc.expectedPath, req.Endpoint)
			}

			// If we have pagination Options, verify they were passed correctly
			if tc.options != nil {
				if req.Options == nil {
					t.Error("expected pagination Options to be passed")
				} else if req.Options.Page != tc.options.Page {
					t.Errorf("expected page %d, got %d", tc.options.Page, req.Options.Page)
				}
			}

			// Verify meta is present for successful Responses
			if meta == nil {
				t.Error("expected meta to be non-nil")
			}
		})
	}
}

func TestAllocationsService_ListAll(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		nodeID        int
		mockResponses []testutil.MockResponse
		expectedError bool
		expectedCount int
		expectedCalls int
	}{
		{
			name:   "Single page of results",
			nodeID: 1,
			mockResponses: []testutil.MockResponse{
				{
					StatusCode: 200,
					Body: []byte(`{
						"object": "list",
						"data": [
							{
								"object": "allocation",
								"attributes": {
									"id": 1,
									"ip": "192.168.1.1",
									"port": 25565,
									"assigned": false
								}
							}
						],
						"meta": {
							"pagination": {
								"total": 1,
								"count": 1,
								"per_page": 100,
								"current_page": 1,
								"total_pages": 1
							}
						}
					}`),
				},
			},
			expectedError: false,
			expectedCount: 1,
			expectedCalls: 1,
		},
		{
			name:   "Multiple pages of results",
			nodeID: 2,
			mockResponses: []testutil.MockResponse{
				{
					StatusCode: 200,
					Body: []byte(`{
						"object": "list",
						"data": [
							{
								"object": "allocation",
								"attributes": {
									"id": 1,
									"ip": "192.168.1.1",
									"port": 25565,
									"assigned": false
								}
							}
						],
						"meta": {
							"pagination": {
								"total": 2,
								"count": 1,
								"per_page": 1,
								"current_page": 1,
								"total_pages": 2
							}
						}
					}`),
				},
				{
					StatusCode: 200,
					Body: []byte(`{
						"object": "list",
						"data": [
							{
								"object": "allocation",
								"attributes": {
									"id": 2,
									"ip": "192.168.1.1",
									"port": 25566,
									"assigned": true
								}
							}
						],
						"meta": {
							"pagination": {
								"total": 2,
								"count": 1,
								"per_page": 1,
								"current_page": 2,
								"total_pages": 2
							}
						}
					}`),
				},
			},
			expectedError: false,
			expectedCount: 2,
			expectedCalls: 2,
		},
		{
			name:   "API error on first page",
			nodeID: 3,
			mockResponses: []testutil.MockResponse{
				{
					StatusCode: 500,
					Body: []byte(`{
						"errors": [{
							"code": "InternalServerError",
							"status": "500",
							"detail": "Internal server error"
						}]
					}`),
				},
			},
			expectedError: true,
			expectedCalls: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{
				Responses: tc.mockResponses,
			}

			service := newAllocationsService(mock, tc.nodeID)

			allocations, err := service.ListAll(context.Background())

			// Check error expectations
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check response expectations
			if len(allocations) != tc.expectedCount {
				t.Errorf("expected %d allocations, got %d", tc.expectedCount, len(allocations))
			}

			// Check number of API calls
			if len(mock.Requests) != tc.expectedCalls {
				t.Errorf("expected %d API calls, got %d", tc.expectedCalls, len(mock.Requests))
			}

			// Verify all Requests were to the correct Endpoint
			for i, req := range mock.Requests {
				expectedPath := fmt.Sprintf("/api/application/nodes/%d/allocations", tc.nodeID)
				if req.Endpoint != expectedPath {
					t.Errorf("request %d: expected path %s, got %s", i, expectedPath, req.Endpoint)
				}

				if req.Method != "GET" {
					t.Errorf("request %d: expected Method GET, got %s", i, req.Method)
				}
			}
		})
	}
}

func TestAllocationsService_Create(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		nodeID        int
		options       api.AllocationCreateOptions
		mockResponse  testutil.MockResponse
		expectedError bool
		expectedBody  string
	}{
		{
			name:   "Successful creation",
			nodeID: 1,
			options: api.AllocationCreateOptions{
				IP:    "192.168.1.1",
				Ports: []string{"25565", "25566", "25567"},
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 204,
				Body:       []byte(""),
			},
			expectedError: false,
			expectedBody:  `{"ip":"192.168.1.1","ports":["25565","25566","25567"]}`,
		},
		{
			name:   "Single port creation",
			nodeID: 2,
			options: api.AllocationCreateOptions{
				IP:    "10.0.0.1",
				Ports: []string{"8080"},
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 204,
				Body:       []byte(""),
			},
			expectedError: false,
			expectedBody:  `{"ip":"10.0.0.1","ports":["8080"]}`,
		},
		{
			name:   "API error response",
			nodeID: 3,
			options: api.AllocationCreateOptions{
				IP:    "invalid-ip",
				Ports: []string{"99999"},
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 422,
				Body: []byte(`{
					"errors": [{
						"code": "ValidationHttpException",
						"status": "422",
						"detail": "The given data was invalid."
					}]
				}`),
			},
			expectedError: true,
			expectedBody:  `{"ip":"invalid-ip","ports":["99999"]}`,
		},
		{
			name:   "Network error",
			nodeID: 4,
			options: api.AllocationCreateOptions{
				IP:    "192.168.1.1",
				Ports: []string{"25565"},
			},
			mockResponse: testutil.MockResponse{
				Err: fmt.Errorf("network timeout"),
			},
			expectedError: true,
			expectedBody:  `{"ip":"192.168.1.1","ports":["25565"]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{
				Responses: []testutil.MockResponse{tc.mockResponse},
			}

			service := newAllocationsService(mock, tc.nodeID)

			err := service.Create(context.Background(), tc.options)

			// Check error expectations
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check request expectations
			if len(mock.Requests) != 1 {
				t.Fatalf("expected 1 request, got %d", len(mock.Requests))
			}

			req := mock.Requests[0]
			if req.Method != "POST" {
				t.Errorf("expected Method POST, got %s", req.Method)
			}

			expectedPath := fmt.Sprintf("/api/application/nodes/%d/allocations", tc.nodeID)
			if req.Endpoint != expectedPath {
				t.Errorf("expected path %s, got %s", expectedPath, req.Endpoint)
			}

			// Check request Body
			bodyStr := strings.TrimSpace(string(req.Body))
			if bodyStr != tc.expectedBody {
				t.Errorf("expected Body %s, got %s", tc.expectedBody, bodyStr)
			}
		})
	}
}

func TestAllocationsService_Delete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		nodeID         int
		allocationID   int
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:         "Successful deletion",
			nodeID:       1,
			allocationID: 123,
			mockResponse: testutil.MockResponse{
				StatusCode: 204,
				Body:       []byte(""),
			},
			expectedError:  false,
			expectedMethod: "DELETE",
			expectedPath:   "/api/application/nodes/1/allocations/123",
		},
		{
			name:         "Allocation not found",
			nodeID:       2,
			allocationID: 999,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body: []byte(`{
					"errors": [{
						"code": "NotFoundHttpException",
						"status": "404",
						"detail": "The requested allocation could not be located."
					}]
				}`),
			},
			expectedError:  true,
			expectedMethod: "DELETE",
			expectedPath:   "/api/application/nodes/2/allocations/999",
		},
		{
			name:         "Allocation in use",
			nodeID:       3,
			allocationID: 456,
			mockResponse: testutil.MockResponse{
				StatusCode: 422,
				Body: []byte(`{
					"errors": [{
						"code": "ValidationHttpException",
						"status": "422",
						"detail": "Cannot delete allocation that is currently assigned to a server."
					}]
				}`),
			},
			expectedError:  true,
			expectedMethod: "DELETE",
			expectedPath:   "/api/application/nodes/3/allocations/456",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{
				Responses: []testutil.MockResponse{tc.mockResponse},
			}

			service := newAllocationsService(mock, tc.nodeID)

			err := service.Delete(context.Background(), tc.allocationID)

			// Check error expectations
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check request expectations
			if len(mock.Requests) != 1 {
				t.Fatalf("expected 1 request, got %d", len(mock.Requests))
			}

			req := mock.Requests[0]
			if req.Method != tc.expectedMethod {
				t.Errorf("expected Method %s, got %s", tc.expectedMethod, req.Method)
			}

			if req.Endpoint != tc.expectedPath {
				t.Errorf("expected path %s, got %s", tc.expectedPath, req.Endpoint)
			}
		})
	}
}

func TestAllocationsService_Integration(t *testing.T) {
	t.Parallel()

	// Test that the service correctly handles the nodeID in all operations
	nodeID := 42
	mock := &testutil.MockRequester{
		Responses: []testutil.MockResponse{
			// List response
			{
				StatusCode: 200,
				Body: []byte(`{
					"object": "list",
					"data": [],
					"meta": {
						"pagination": {
							"total": 0,
							"count": 0,
							"per_page": 100,
							"current_page": 1,
							"total_pages": 0
						}
					}
				}`),
			},
			// ListAll response
			{
				StatusCode: 200,
				Body: []byte(`{
					"object": "list",
					"data": [],
					"meta": {
						"pagination": {
							"total": 0,
							"count": 0,
							"per_page": 100,
							"current_page": 1,
							"total_pages": 0
						}
					}
				}`),
			},
			// Create response
			{
				StatusCode: 204,
				Body:       []byte(""),
			},
			// Delete response
			{
				StatusCode: 204,
				Body:       []byte(""),
			},
		},
	}

	service := newAllocationsService(mock, nodeID)

	// Test List
	_, _, err := service.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	// Test ListAll
	_, err = service.ListAll(context.Background())
	if err != nil {
		t.Fatalf("ListAll failed: %v", err)
	}

	// Test Create
	createOptions := api.AllocationCreateOptions{
		IP:    "192.168.1.1",
		Ports: []string{"25565"},
	}
	err = service.Create(context.Background(), createOptions)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Test Delete
	err = service.Delete(context.Background(), 123)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify all Requests used the correct nodeID
	expectedBasePath := fmt.Sprintf("/api/application/nodes/%d/allocations", nodeID)
	for i, req := range mock.Requests {
		if !strings.HasPrefix(req.Endpoint, expectedBasePath) {
			t.Errorf("request %d: expected Endpoint to start with %s, got %s", i, expectedBasePath, req.Endpoint)
		}
	}

	// Verify we made exactly 4 Requests
	if len(mock.Requests) != 4 {
		t.Errorf("expected 4 Requests, got %d", len(mock.Requests))
	}
}

func TestAllocationsService_EdgeCases(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		nodeID        int
		allocationID  int
		options       api.AllocationCreateOptions
		mockResponse  testutil.MockResponse
		expectedError bool
		description   string
	}{
		{
			name:         "Zero node ID",
			nodeID:       0,
			allocationID: 123,
			mockResponse: testutil.MockResponse{
				StatusCode: 204,
				Body:       []byte(""),
			},
			expectedError: false,
			description:   "Should handle zero node ID gracefully",
		},
		{
			name:         "Negative allocation ID",
			nodeID:       1,
			allocationID: -1,
			mockResponse: testutil.MockResponse{
				StatusCode: 400,
				Body: []byte(`{
					"errors": [{
						"code": "BadRequestHttpException",
						"status": "400",
						"detail": "Invalid allocation ID"
					}]
				}`),
			},
			expectedError: true,
			description:   "Should handle negative allocation ID with error",
		},
		{
			name:   "Empty ports array",
			nodeID: 1,
			options: api.AllocationCreateOptions{
				IP:    "192.168.1.1",
				Ports: []string{},
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 422,
				Body: []byte(`{
					"errors": [{
						"code": "ValidationHttpException",
						"status": "422",
						"detail": "At least one port must be specified"
					}]
				}`),
			},
			expectedError: true,
			description:   "Should handle empty ports array with validation error",
		},
		{
			name:   "Empty IP address",
			nodeID: 1,
			options: api.AllocationCreateOptions{
				IP:    "",
				Ports: []string{"25565"},
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 422,
				Body: []byte(`{
					"errors": [{
						"code": "ValidationHttpException",
						"status": "422",
						"detail": "IP address is required"
					}]
				}`),
			},
			expectedError: true,
			description:   "Should handle empty IP address with validation error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{
				Responses: []testutil.MockResponse{tc.mockResponse},
			}

			service := newAllocationsService(mock, tc.nodeID)

			var err error
			if tc.options.IP != "" || len(tc.options.Ports) > 0 {
				// Test Create
				err = service.Create(context.Background(), tc.options)
			} else {
				// Test Delete
				err = service.Delete(context.Background(), tc.allocationID)
			}

			// Check error expectations
			if tc.expectedError {
				if err == nil {
					t.Errorf("expected error for %s but got none", tc.description)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %s: %v", tc.description, err)
				}
			}
		})
	}
}
