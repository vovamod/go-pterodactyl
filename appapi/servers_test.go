package appapi

import (
	"context"
	"github.com/vovamod/go-pterodactyl/internal/testutil"
	"testing"

	"github.com/vovamod/go-pterodactyl/api"
)

func TestServersService_List(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		options        api.PaginationOptions
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedCount  int
		expectedMethod string
		expectedPath   string
	}{
		{
			name:    "Successful list with pagination",
			options: api.PaginationOptions{Page: 1, PerPage: 10},
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body: []byte(`{"object": "list", "data": [
					{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "Server1", "description": "desc1", "suspended": false, "user": 1, "node": 1, "allocation": 1, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}}
				], "meta": {"pagination": {"total": 1, "count": 1, "per_page": 10, "current_page": 1, "total_pages": 1}}}`),
			},
			expectedError:  false,
			expectedCount:  1,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers",
		},
		{
			name:    "API error response",
			options: api.PaginationOptions{Page: 1},
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			servers, meta, err := service.List(context.Background(), tc.options)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(servers) != tc.expectedCount {
				t.Errorf("expected %d servers, got %d", tc.expectedCount, len(servers))
			}
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
			if meta == nil {
				t.Error("expected meta to be non-nil")
			}
			if tc.expectedCount > 0 && len(servers) > 0 {
				server := servers[0]
				if server.ID == 0 {
					t.Error("expected server ID to be non-zero")
				}
				if server.Name == "" {
					t.Error("expected server name to be non-empty")
				}
			}
		})
	}
}

func TestServersService_ListAll(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedCount  int
		expectedMethod string
		expectedPath   string
	}{
		{
			name: "Successful list all",
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body: []byte(`{"object": "list", "data": [
					{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "Server1", "description": "desc1", "suspended": false, "user": 1, "node": 1, "allocation": 1, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}},
					{"object": "server", "attributes": {"id": 2, "uuid": "uuid-2", "identifier": "id2", "name": "Server2", "description": "desc2", "suspended": true, "user": 2, "node": 1, "allocation": 2, "nest": 1, "egg": 1, "created_at": "2023-01-02T00:00:00Z"}}
				], "meta": {"pagination": {"total": 2, "count": 2, "per_page": 100, "current_page": 1, "total_pages": 1}}}`),
			},
			expectedError:  false,
			expectedCount:  2,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers",
		},
		{
			name: "Empty list",
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body:       []byte(`{"object": "list", "data": [], "meta": {"pagination": {"total": 0, "count": 0, "per_page": 100, "current_page": 1, "total_pages": 0}}}`),
			},
			expectedError:  false,
			expectedCount:  0,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers",
		},
		{
			name: "API error response",
			mockResponse: testutil.MockResponse{
				StatusCode: 500,
				Body:       []byte(`{"errors": [{"code": "InternalServerError", "status": "500", "detail": "Internal server error."}]}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			servers, err := service.ListAll(context.Background())
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(servers) != tc.expectedCount {
				t.Errorf("expected %d servers, got %d", tc.expectedCount, len(servers))
			}
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
			if tc.expectedCount > 0 && len(servers) > 0 {
				server := servers[0]
				if server.ID == 0 {
					t.Error("expected server ID to be non-zero")
				}
				if server.Name == "" {
					t.Error("expected server name to be non-empty")
				}
			}
		})
	}
}

func TestServersService_Get(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		serverID       int
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:     "Successful get",
			serverID: 1,
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body:       []byte(`{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "TestServer", "description": "Test description", "suspended": false, "user": 1, "node": 1, "allocation": 1, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}}`),
			},
			expectedError:  false,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers/1",
		},
		{
			name:     "Server not found",
			serverID: 999,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Server not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers/999",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			server, err := service.Get(context.Background(), tc.serverID)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if server == nil {
				t.Fatal("expected server to be non-nil")
			}
			if server.ID != tc.serverID {
				t.Errorf("expected server ID %d, got %d", tc.serverID, server.ID)
			}
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

func TestServersService_GetExternal(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		externalID     string
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:       "Successful get external",
			externalID: "external-123",
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body:       []byte(`{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "TestServer", "description": "Test description", "suspended": false, "user": 1, "node": 1, "allocation": 1, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}}`),
			},
			expectedError:  false,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers/external/external-123",
		},
		{
			name:       "External server not found",
			externalID: "nonexistent",
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "External server not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers/external/nonexistent",
		},
		{
			name:       "URL encoding test",
			externalID: "external/with/slashes",
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body:       []byte(`{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "TestServer", "description": "Test description", "suspended": false, "user": 1, "node": 1, "allocation": 1, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}}`),
			},
			expectedError:  false,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers/external/external%2Fwith%2Fslashes",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			server, err := service.GetExternal(context.Background(), tc.externalID)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if server == nil {
				t.Fatal("expected server to be non-nil")
			}
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

func TestServersService_Create(t *testing.T) {
	t.Parallel()
	description := "Test server description"
	nodeID := 1
	allocation := &struct {
		Default    int   `json:"default"`
		Additional []int `json:"additional,omitempty"`
	}{
		Default: 1,
	}

	testCases := []struct {
		name           string
		options        api.ServerCreateOptions
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name: "Successful create",
			options: api.ServerCreateOptions{
				Name:        "TestServer",
				User:        1,
				NodeID:      &nodeID,
				Allocation:  allocation,
				Nest:        1,
				Egg:         1,
				Description: &description,
				Limits: api.ServerLimits{
					Memory: 1024,
					Swap:   512,
					Disk:   10000,
					IO:     500,
					CPU:    100,
				},
				FeatureLimits: api.ServerFeatureLimits{
					Databases:   5,
					Allocations: 1,
					Backups:     2,
				},
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 201,
				Body:       []byte(`{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "TestServer", "description": "Test server description", "suspended": false, "user": 1, "node": 1, "allocation": 1, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}}`),
			},
			expectedError:  false,
			expectedMethod: "POST",
			expectedPath:   "/api/application/servers",
		},
		{
			name: "API error response",
			options: api.ServerCreateOptions{
				Name:       "TestServer",
				User:       1,
				NodeID:     &nodeID,
				Allocation: allocation,
				Nest:       1,
				Egg:        1,
				Limits: api.ServerLimits{
					Memory: 1024,
					Swap:   512,
					Disk:   10000,
					IO:     500,
					CPU:    100,
				},
				FeatureLimits: api.ServerFeatureLimits{
					Databases:   5,
					Allocations: 1,
					Backups:     2,
				},
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 422,
				Body:       []byte(`{"errors": [{"code": "ValidationException", "status": "422", "detail": "Validation failed."}]}`),
			},
			expectedError:  true,
			expectedMethod: "POST",
			expectedPath:   "/api/application/servers",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			server, err := service.Create(context.Background(), tc.options)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if server == nil {
				t.Fatal("expected server to be non-nil")
			}
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
			// Verify request Body contains expected data
			if req.Body == nil {
				t.Error("expected request Body to be non-nil")
			}
		})
	}
}

func TestServersService_UpdateDetails(t *testing.T) {
	t.Parallel()
	description := "Updated description"
	externalID := "external-123"

	testCases := []struct {
		name           string
		serverID       int
		options        api.ServerUpdateDetailsOptions
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:     "Successful update details",
			serverID: 1,
			options: api.ServerUpdateDetailsOptions{
				Name:        "UpdatedServer",
				Description: &description,
				ExternalID:  &externalID,
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body:       []byte(`{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "UpdatedServer", "description": "Updated description", "suspended": false, "user": 1, "node": 1, "allocation": 1, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}}`),
			},
			expectedError:  false,
			expectedMethod: "PATCH",
			expectedPath:   "/api/application/servers/1/details",
		},
		{
			name:     "Server not found",
			serverID: 999,
			options: api.ServerUpdateDetailsOptions{
				Name: "UpdatedServer",
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Server not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "PATCH",
			expectedPath:   "/api/application/servers/999/details",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			server, err := service.UpdateDetails(context.Background(), tc.serverID, tc.options)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if server == nil {
				t.Fatal("expected server to be non-nil")
			}
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

func TestServersService_UpdateBuild(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		serverID       int
		options        api.ServerUpdateBuildOptions
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:     "Successful update build",
			serverID: 1,
			options: api.ServerUpdateBuildOptions{
				Allocation: 2,
				Memory:     1024,
				Swap:       512,
				Disk:       10000,
				IO:         500,
				CPU:        100,
				FeatureLimits: api.ServerFeatureLimits{
					Databases:   5,
					Allocations: 1,
					Backups:     2,
				},
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body:       []byte(`{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "TestServer", "description": "Test description", "suspended": false, "user": 1, "node": 1, "allocation": 2, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}}`),
			},
			expectedError:  false,
			expectedMethod: "PATCH",
			expectedPath:   "/api/application/servers/1/build",
		},
		{
			name:     "Invalid build configuration",
			serverID: 1,
			options: api.ServerUpdateBuildOptions{
				Allocation: 1,
				Memory:     -1, // Invalid memory value
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 422,
				Body:       []byte(`{"errors": [{"code": "ValidationException", "status": "422", "detail": "Invalid memory allocation."}]}`),
			},
			expectedError:  true,
			expectedMethod: "PATCH",
			expectedPath:   "/api/application/servers/1/build",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			server, err := service.UpdateBuild(context.Background(), tc.serverID, tc.options)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if server == nil {
				t.Fatal("expected server to be non-nil")
			}
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

func TestServersService_UpdateStartup(t *testing.T) {
	t.Parallel()
	environment := map[string]string{
		"JAVA_MEMORY": "1024M",
		"PORT":        "25565",
	}

	testCases := []struct {
		name           string
		serverID       int
		options        api.ServerUpdateStartupOptions
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:     "Successful update startup",
			serverID: 1,
			options: api.ServerUpdateStartupOptions{
				Startup:     "java -jar server.jar",
				Environment: &environment,
				Egg:         1,
				Image:       "openjdk:11",
				SkipScripts: false,
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body:       []byte(`{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "TestServer", "description": "Test description", "suspended": false, "user": 1, "node": 1, "allocation": 1, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}}`),
			},
			expectedError:  false,
			expectedMethod: "PATCH",
			expectedPath:   "/api/application/servers/1/startup",
		},
		{
			name:     "Invalid startup configuration",
			serverID: 1,
			options: api.ServerUpdateStartupOptions{
				Startup: "", // Empty startup command
				Egg:     1,
				Image:   "openjdk:11",
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 422,
				Body:       []byte(`{"errors": [{"code": "ValidationException", "status": "422", "detail": "Startup command cannot be empty."}]}`),
			},
			expectedError:  true,
			expectedMethod: "PATCH",
			expectedPath:   "/api/application/servers/1/startup",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			server, err := service.UpdateStartup(context.Background(), tc.serverID, tc.options)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if server == nil {
				t.Fatal("expected server to be non-nil")
			}
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

func TestServersService_Suspend(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		serverID       int
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:     "Successful suspend",
			serverID: 1,
			mockResponse: testutil.MockResponse{
				StatusCode: 204,
				Body:       []byte(``),
			},
			expectedError:  false,
			expectedMethod: "POST",
			expectedPath:   "/api/application/servers/1/suspend",
		},
		{
			name:     "Server not found",
			serverID: 999,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Server not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "POST",
			expectedPath:   "/api/application/servers/999/suspend",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			err := service.Suspend(context.Background(), tc.serverID)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
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

func TestServersService_Unsuspend(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		serverID       int
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:     "Successful unsuspend",
			serverID: 1,
			mockResponse: testutil.MockResponse{
				StatusCode: 204,
				Body:       []byte(``),
			},
			expectedError:  false,
			expectedMethod: "POST",
			expectedPath:   "/api/application/servers/1/unsuspend",
		},
		{
			name:     "Server not found",
			serverID: 999,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Server not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "POST",
			expectedPath:   "/api/application/servers/999/unsuspend",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			err := service.Unsuspend(context.Background(), tc.serverID)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
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

func TestServersService_Reinstall(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		serverID       int
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:     "Successful reinstall",
			serverID: 1,
			mockResponse: testutil.MockResponse{
				StatusCode: 204,
				Body:       []byte(``),
			},
			expectedError:  false,
			expectedMethod: "POST",
			expectedPath:   "/api/application/servers/1/reinstall",
		},
		{
			name:     "Server not found",
			serverID: 999,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Server not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "POST",
			expectedPath:   "/api/application/servers/999/reinstall",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			err := service.Reinstall(context.Background(), tc.serverID)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
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

func TestServersService_Delete(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		serverID       int
		force          bool
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:     "Successful delete",
			serverID: 1,
			force:    false,
			mockResponse: testutil.MockResponse{
				StatusCode: 204,
				Body:       []byte(``),
			},
			expectedError:  false,
			expectedMethod: "DELETE",
			expectedPath:   "/api/application/servers/1",
		},
		{
			name:     "Successful force delete",
			serverID: 1,
			force:    true,
			mockResponse: testutil.MockResponse{
				StatusCode: 204,
				Body:       []byte(``),
			},
			expectedError:  false,
			expectedMethod: "DELETE",
			expectedPath:   "/api/application/servers/1",
		},
		{
			name:     "Server not found",
			serverID: 999,
			force:    false,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Server not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "DELETE",
			expectedPath:   "/api/application/servers/999",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewServersService(mock)
			err := service.Delete(context.Background(), tc.serverID, tc.force)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
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
			// Check if force delete includes request Body
			if tc.force && req.Body == nil {
				t.Error("expected request Body for force delete")
			}
		})
	}
}

func TestServersService_Databases(t *testing.T) {
	t.Parallel()
	service := NewServersService(&testutil.MockRequester{})

	// Test that Databases returns a non-nil DatabaseService
	databaseService := service.Databases(context.Background(), 1)
	if databaseService == nil {
		t.Fatal("expected DatabaseService to be non-nil")
	}

}

// Test data validation and edge cases
func TestServersService_DataValidation(t *testing.T) {
	t.Parallel()

	t.Run("Empty server list response", func(t *testing.T) {
		mock := &testutil.MockRequester{Responses: []testutil.MockResponse{
			{
				StatusCode: 200,
				Body:       []byte(`{"object": "list", "data": [], "meta": {"pagination": {"total": 0, "count": 0, "per_page": 100, "current_page": 1, "total_pages": 0}}}`),
			},
		}}
		service := NewServersService(mock)
		servers, err := service.ListAll(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(servers) != 0 {
			t.Errorf("expected 0 servers, got %d", len(servers))
		}
	})

	t.Run("Invalid JSON response", func(t *testing.T) {
		mock := &testutil.MockRequester{Responses: []testutil.MockResponse{
			{
				StatusCode: 200,
				Body:       []byte(`invalid json`),
			},
		}}
		service := NewServersService(mock)
		_, err := service.ListAll(context.Background())
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("Server with minimal data", func(t *testing.T) {
		mock := &testutil.MockRequester{Responses: []testutil.MockResponse{
			{
				StatusCode: 200,
				Body:       []byte(`{"object": "server", "attributes": {"id": 1, "uuid": "uuid-1", "identifier": "id1", "name": "MinimalServer", "suspended": false, "user": 1, "node": 1, "allocation": 1, "nest": 1, "egg": 1, "created_at": "2023-01-01T00:00:00Z"}}`),
			},
		}}
		service := NewServersService(mock)
		server, err := service.Get(context.Background(), 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if server == nil {
			t.Fatal("expected server to be non-nil")
		}
		if server.ID != 1 {
			t.Errorf("expected server ID 1, got %d", server.ID)
		}
		if server.Name != "MinimalServer" {
			t.Errorf("expected server name 'MinimalServer', got '%s'", server.Name)
		}
	})
}

// Test constructor
func TestNewServersService(t *testing.T) {
	t.Parallel()
	mock := &testutil.MockRequester{}
	service := NewServersService(mock)
	if service == nil {
		t.Fatal("expected service to be non-nil")
	}
}
