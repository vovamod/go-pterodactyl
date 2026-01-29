package appapi

import (
	"context"
	"fmt"
	"github.com/vovamod/go-pterodactyl/internal/testutil"
	"strings"
	"testing"
	"time"

	"github.com/vovamod/go-pterodactyl/api"
)

func TestEggsService_List(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		nestID         int
		options        *api.PaginationOptions
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedCount  int
		expectedMethod string
		expectedPath   string
	}{
		{
			name:   "Successful list with pagination",
			nestID: 1,
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
							"object": "egg",
							"attributes": {
								"id": 1,
								"uuid": "12345678-1234-1234-1234-123456789abc",
								"nest_id": 1,
								"author": "Pterodactyl Team",
								"description": "Minecraft Server",
								"docker_images": {
									"java": "ghcr.io/pterodactyl/yolks:java_17",
									"bedrock": "ghcr.io/pterodactyl/yolks:bedrock"
								},
								"config": {
									"files": {},
									"startup": "java -Xms128M -Xmx{{SERVER_MEMORY}}M -jar {{SERVER_JARFILE}}",
									"stop": "^C",
									"logs": {}
								},
								"startup": "java -Xms128M -Xmx{{SERVER_MEMORY}}M -jar {{SERVER_JARFILE}}",
								"script": {
									"privileged": false,
									"install": "echo \"Installing Minecraft...\"",
									"entry": "java -Xms128M -Xmx{{SERVER_MEMORY}}M -jar {{SERVER_JARFILE}}",
									"container": "ghcr.io/pterodactyl/yolks:java_17",
									"extends": null
								},
								"created_at": "2023-01-01T00:00:00Z",
								"updated_at": "2023-01-01T00:00:00Z"
							}
						},
						{
							"object": "egg",
							"attributes": {
								"id": 2,
								"uuid": "87654321-4321-4321-4321-cba987654321",
								"nest_id": 1,
								"author": "Pterodactyl Team",
								"description": "Source Dedicated Server",
								"docker_images": {
									"source": "ghcr.io/pterodactyl/yolks:source"
								},
								"config": {
									"files": {},
									"startup": "./srcds_run -game {{SRCDS_GAME}} -console -port {{SERVER_PORT}} +map {{SRCDS_MAP}}",
									"stop": "^C",
									"logs": {}
								},
								"startup": "./srcds_run -game {{SRCDS_GAME}} -console -port {{SERVER_PORT}} +map {{SRCDS_MAP}}",
								"script": {
									"privileged": false,
									"install": "echo \"Installing Source Dedicated Server...\"",
									"entry": "./srcds_run -game {{SRCDS_GAME}} -console -port {{SERVER_PORT}} +map {{SRCDS_MAP}}",
									"container": "ghcr.io/pterodactyl/yolks:source",
									"extends": null
								},
								"created_at": "2023-01-02T00:00:00Z",
								"updated_at": "2023-01-02T00:00:00Z"
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
			expectedPath:   "/api/application/nests/1/eggs",
		},
		{
			name:    "Successful list without pagination",
			nestID:  2,
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
			expectedPath:   "/api/application/nests/2/eggs",
		},
		{
			name:    "API error response",
			nestID:  3,
			options: &api.PaginationOptions{Page: 1},
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body: []byte(`{
					"errors": [{
						"code": "NotFoundHttpException",
						"status": "404",
						"detail": "The requested nest could not be located."
					}]
				}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/nests/3/eggs",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{
				Responses: []testutil.MockResponse{tc.mockResponse},
			}

			service := NewEggsService(mock, tc.nestID)

			eggs, meta, err := service.List(context.Background(), tc.options)

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
			if len(eggs) != tc.expectedCount {
				t.Errorf("expected %d eggs, got %d", tc.expectedCount, len(eggs))
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

			// Verify egg data if we have any
			if tc.expectedCount > 0 && len(eggs) > 0 {
				egg := eggs[0]
				if egg.ID == 0 {
					t.Error("expected egg ID to be non-zero")
				}
				if egg.Description == "" {
					t.Error("expected egg description to be non-empty")
				}
				if egg.UUID == "" {
					t.Error("expected egg UUID to be non-empty")
				}
			}
		})
	}
}

func TestEggsService_ListAll(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		nestID        int
		mockResponses []testutil.MockResponse
		expectedError bool
		expectedCount int
		expectedCalls int
	}{
		{
			name:   "Single page of results",
			nestID: 1,
			mockResponses: []testutil.MockResponse{
				{
					StatusCode: 200,
					Body: []byte(`{
						"object": "list",
						"data": [
							{
								"object": "egg",
								"attributes": {
									"id": 1,
									"uuid": "12345678-1234-1234-1234-123456789abc",
									"nest_id": 1,
									"author": "Pterodactyl Team",
									"description": "Minecraft Server",
									"docker_images": {
										"java": "ghcr.io/pterodactyl/yolks:java_17"
									},
									"config": {},
									"startup": "java -Xms128M -Xmx{{SERVER_MEMORY}}M -jar {{SERVER_JARFILE}}",
									"script": {},
									"created_at": "2023-01-01T00:00:00Z",
									"updated_at": "2023-01-01T00:00:00Z"
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
			nestID: 2,
			mockResponses: []testutil.MockResponse{
				{
					StatusCode: 200,
					Body: []byte(`{
						"object": "list",
						"data": [
							{
								"object": "egg",
								"attributes": {
									"id": 1,
									"uuid": "12345678-1234-1234-1234-123456789abc",
									"nest_id": 2,
									"author": "Pterodactyl Team",
									"description": "Page 1 Egg",
									"docker_images": {},
									"config": {},
									"startup": "echo 'page 1'",
									"script": {},
									"created_at": "2023-01-01T00:00:00Z",
									"updated_at": "2023-01-01T00:00:00Z"
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
								"object": "egg",
								"attributes": {
									"id": 2,
									"uuid": "87654321-4321-4321-4321-cba987654321",
									"nest_id": 2,
									"author": "Pterodactyl Team",
									"description": "Page 2 Egg",
									"docker_images": {},
									"config": {},
									"startup": "echo 'page 2'",
									"script": {},
									"created_at": "2023-01-02T00:00:00Z",
									"updated_at": "2023-01-02T00:00:00Z"
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
			nestID: 3,
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

			service := NewEggsService(mock, tc.nestID)

			eggs, err := service.ListAll(context.Background())

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
			if len(eggs) != tc.expectedCount {
				t.Errorf("expected %d eggs, got %d", tc.expectedCount, len(eggs))
			}

			// Check number of API calls
			if len(mock.Requests) != tc.expectedCalls {
				t.Errorf("expected %d API calls, got %d", tc.expectedCalls, len(mock.Requests))
			}

			// Verify all Requests were to the correct Endpoint
			for i, req := range mock.Requests {
				expectedPath := fmt.Sprintf("/api/application/nests/%d/eggs", tc.nestID)
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

func TestEggsService_Get(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		nestID         int
		eggID          int
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:   "Successful get",
			nestID: 1,
			eggID:  123,
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body: []byte(`{
					"object": "egg",
					"attributes": {
						"id": 123,
						"uuid": "12345678-1234-1234-1234-123456789abc",
						"nest_id": 1,
						"author": "Pterodactyl Team",
						"description": "Minecraft Server",
						"docker_images": {
							"java": "ghcr.io/pterodactyl/yolks:java_17",
							"bedrock": "ghcr.io/pterodactyl/yolks:bedrock"
						},
						"config": {
							"files": {},
							"startup": "java -Xms128M -Xmx{{SERVER_MEMORY}}M -jar {{SERVER_JARFILE}}",
							"stop": "^C",
							"logs": {}
						},
						"startup": "java -Xms128M -Xmx{{SERVER_MEMORY}}M -jar {{SERVER_JARFILE}}",
						"script": {
							"privileged": false,
							"install": "echo \"Installing Minecraft...\"",
							"entry": "java -Xms128M -Xmx{{SERVER_MEMORY}}M -jar {{SERVER_JARFILE}}",
							"container": "ghcr.io/pterodactyl/yolks:java_17",
							"extends": null
						},
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z"
					}
				}`),
			},
			expectedError:  false,
			expectedMethod: "GET",
			expectedPath:   "/api/application/nests/1/eggs/123",
		},
		{
			name:   "Egg not found",
			nestID: 2,
			eggID:  999,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body: []byte(`{
					"errors": [{
						"code": "NotFoundHttpException",
						"status": "404",
						"detail": "The requested egg could not be located."
					}]
				}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/nests/2/eggs/999",
		},
		{
			name:   "Nest not found",
			nestID: 999,
			eggID:  123,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body: []byte(`{
					"errors": [{
						"code": "NotFoundHttpException",
						"status": "404",
						"detail": "The requested nest could not be located."
					}]
				}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/nests/999/eggs/123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{
				Responses: []testutil.MockResponse{tc.mockResponse},
			}

			service := NewEggsService(mock, tc.nestID)

			egg, err := service.Get(context.Background(), tc.eggID)

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

			// Verify egg data
			if egg == nil {
				t.Error("expected egg to be non-nil")
			} else {
				if egg.ID != tc.eggID {
					t.Errorf("expected egg ID %d, got %d", tc.eggID, egg.ID)
				}
				if egg.Description == "" {
					t.Error("expected egg description to be non-empty")
				}
				if egg.UUID == "" {
					t.Error("expected egg UUID to be non-empty")
				}
			}
		})
	}
}

func TestEggsService_Integration(t *testing.T) {
	t.Parallel()

	// Test that the service correctly handles the nestID in all operations
	nestID := 42
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
			// Get response
			{
				StatusCode: 200,
				Body: []byte(`{
					"object": "egg",
					"attributes": {
						"id": 123,
						"uuid": "12345678-1234-1234-1234-123456789abc",
						"nest_id": 42,
						"author": "Test Author",
						"description": "Test Egg",
						"docker_images": {},
						"config": {},
						"startup": "echo 'test'",
						"script": {},
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z"
					}
				}`),
			},
		},
	}

	service := NewEggsService(mock, nestID)

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

	// Test Get
	_, err = service.Get(context.Background(), 123)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	// Verify all Requests used the correct nestID
	expectedBasePath := fmt.Sprintf("/api/application/nests/%d/eggs", nestID)
	for i, req := range mock.Requests {
		if !strings.HasPrefix(req.Endpoint, expectedBasePath) {
			t.Errorf("request %d: expected Endpoint to start with %s, got %s", i, expectedBasePath, req.Endpoint)
		}
	}

	// Verify we made exactly 3 Requests
	if len(mock.Requests) != 3 {
		t.Errorf("expected 3 Requests, got %d", len(mock.Requests))
	}
}

func TestEggsService_EdgeCases(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		nestID        int
		eggID         int
		mockResponse  testutil.MockResponse
		expectedError bool
		description   string
	}{
		{
			name:   "Zero nest ID",
			nestID: 0,
			eggID:  123,
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body: []byte(`{
					"object": "egg",
					"attributes": {
						"id": 123,
						"uuid": "12345678-1234-1234-1234-123456789abc",
						"nest_id": 0,
						"author": "Test Author",
						"description": "Test Egg",
						"docker_images": {},
						"config": {},
						"startup": "echo 'test'",
						"script": {},
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z"
					}
				}`),
			},
			expectedError: false,
			description:   "Should handle zero nest ID gracefully",
		},
		{
			name:   "Negative egg ID",
			nestID: 1,
			eggID:  -1,
			mockResponse: testutil.MockResponse{
				StatusCode: 400,
				Body: []byte(`{
					"errors": [{
						"code": "BadRequestHttpException",
						"status": "400",
						"detail": "Invalid egg ID"
					}]
				}`),
			},
			expectedError: true,
			description:   "Should handle negative egg ID with error",
		},
		{
			name:   "Zero egg ID",
			nestID: 1,
			eggID:  0,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body: []byte(`{
					"errors": [{
						"code": "NotFoundHttpException",
						"status": "404",
						"detail": "The requested egg could not be located."
					}]
				}`),
			},
			expectedError: true,
			description:   "Should handle zero egg ID with error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{
				Responses: []testutil.MockResponse{tc.mockResponse},
			}

			service := NewEggsService(mock, tc.nestID)

			_, err := service.Get(context.Background(), tc.eggID)

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

func TestEggsService_DataValidation(t *testing.T) {
	t.Parallel()

	// Test that egg data is properly parsed
	mock := &testutil.MockRequester{
		Responses: []testutil.MockResponse{
			{
				StatusCode: 200,
				Body: []byte(`{
					"object": "egg",
					"attributes": {
						"id": 123,
						"uuid": "12345678-1234-1234-1234-123456789abc",
						"nest_id": 1,
						"author": "Test Author",
						"description": "Test Egg Description",
						"docker_images": {
							"java": "ghcr.io/pterodactyl/yolks:java_17",
							"bedrock": "ghcr.io/pterodactyl/yolks:bedrock"
						},
						"config": {
							"files": {},
							"startup": "java -Xms128M -Xmx{{SERVER_MEMORY}}M -jar {{SERVER_JARFILE}}",
							"stop": "^C",
							"logs": {}
						},
						"startup": "java -Xms128M -Xmx{{SERVER_MEMORY}}M -jar {{SERVER_JARFILE}}",
						"script": {
							"privileged": false,
							"install": "echo \"Installing...\"",
							"entry": "java -Xms128M -Xmx{{SERVER_MEMORY}}M -jar {{SERVER_JARFILE}}",
							"container": "ghcr.io/pterodactyl/yolks:java_17",
							"extends": null
						},
						"created_at": "2023-01-01T12:00:00Z",
						"updated_at": "2023-01-02T12:00:00Z"
					}
				}`),
			},
		},
	}

	service := NewEggsService(mock, 1)

	egg, err := service.Get(context.Background(), 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify all fields are correctly parsed
	if egg.ID != 123 {
		t.Errorf("expected ID 123, got %d", egg.ID)
	}
	if egg.UUID != "12345678-1234-1234-1234-123456789abc" {
		t.Errorf("expected UUID '12345678-1234-1234-1234-123456789abc', got '%s'", egg.UUID)
	}
	if egg.NestID != 1 {
		t.Errorf("expected NestID 1, got %d", egg.NestID)
	}
	if egg.Author != "Test Author" {
		t.Errorf("expected Author 'Test Author', got '%s'", egg.Author)
	}
	if egg.Description != "Test Egg Description" {
		t.Errorf("expected Description 'Test Egg Description', got '%s'", egg.Description)
	}

	// Verify docker images
	if len(egg.DockerImages) != 2 {
		t.Errorf("expected 2 docker images, got %d", len(egg.DockerImages))
	}
	if egg.DockerImages["java"] != "ghcr.io/pterodactyl/yolks:java_17" {
		t.Errorf("expected java image 'ghcr.io/pterodactyl/yolks:java_17', got '%s'", egg.DockerImages["java"])
	}
	if egg.DockerImages["bedrock"] != "ghcr.io/pterodactyl/yolks:bedrock" {
		t.Errorf("expected bedrock image 'ghcr.io/pterodactyl/yolks:bedrock', got '%s'", egg.DockerImages["bedrock"])
	}

	// Verify timestamps are parsed correctly
	expectedCreatedAt, _ := time.Parse(time.RFC3339, "2023-01-01T12:00:00Z")
	if !egg.CreatedAt.Equal(expectedCreatedAt) {
		t.Errorf("expected CreatedAt %v, got %v", expectedCreatedAt, egg.CreatedAt)
	}

	expectedUpdatedAt, _ := time.Parse(time.RFC3339, "2023-01-02T12:00:00Z")
	if !egg.UpdatedAt.Equal(expectedUpdatedAt) {
		t.Errorf("expected UpdatedAt %v, got %v", expectedUpdatedAt, egg.UpdatedAt)
	}
}

func TestEggsService_Constructor(t *testing.T) {
	t.Parallel()

	// Test the constructor function
	mock := &testutil.MockRequester{}
	nestID := 42

	service := NewEggsService(mock, nestID)

	if service == nil {
		t.Fatal("expected service to be non-nil")
	}

	// Verify the service has the correct nestID
	// Note: We can't directly access the nestID field since it's unexported,
	// but we can verify it's working by making a request and checking the Endpoint
	mock.Responses = []testutil.MockResponse{
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
	}

	_, _, err := service.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(mock.Requests) != 1 {
		t.Fatalf("expected 1 request, got %d", len(mock.Requests))
	}

	expectedPath := fmt.Sprintf("/api/application/nests/%d/eggs", nestID)
	if mock.Requests[0].Endpoint != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, mock.Requests[0].Endpoint)
	}
}
