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

func TestDatabaseService_List(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		serverID       int
		options        api.PaginationOptions
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedCount  int
		expectedMethod string
		expectedPath   string
	}{
		{
			name:     "Successful list with pagination",
			serverID: 1,
			options: api.PaginationOptions{
				Page:    1,
				PerPage: 10,
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body: []byte(`{
					"object": "list",
					"data": [
						{
							"object": "server_database",
							"attributes": {
								"id": 1,
								"server": 1,
								"host": 1,
								"database": "minecraft_db",
								"username": "minecraft_user",
								"remote": "%",
								"max_connections": 10,
								"created_at": "2023-01-01T00:00:00Z",
								"updated_at": "2023-01-01T00:00:00Z"
							}
						},
						{
							"object": "server_database",
							"attributes": {
								"id": 2,
								"server": 1,
								"host": 1,
								"database": "webapp_db",
								"username": "webapp_user",
								"remote": "192.168.1.100",
								"max_connections": null,
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
			expectedPath:   "/api/application/servers/1/databases",
		},
		{
			name:     "Successful list without pagination",
			serverID: 2,
			options:  api.PaginationOptions{},
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
			expectedPath:   "/api/application/servers/2/databases",
		},
		{
			name:     "API error response",
			serverID: 3,
			options:  api.PaginationOptions{Page: 1},
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body: []byte(`{
					"errors": [{
						"code": "NotFoundHttpException",
						"status": "404",
						"detail": "The requested server could not be located."
					}]
				}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers/3/databases",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{
				Responses: []testutil.MockResponse{tc.mockResponse},
			}

			service := newDatabaseService(mock, tc.serverID)

			databases, meta, err := service.List(context.Background(), tc.options)

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
			if len(databases) != tc.expectedCount {
				t.Errorf("expected %d databases, got %d", tc.expectedCount, len(databases))
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

			// Verify meta is present for successful Responses
			if meta == nil {
				t.Error("expected meta to be non-nil")
			}

			// Verify database data if we have any
			if tc.expectedCount > 0 && len(databases) > 0 {
				db := databases[0]
				if db.ID == 0 {
					t.Error("expected database ID to be non-zero")
				}
				if db.DatabaseName == "" {
					t.Error("expected database name to be non-empty")
				}
			}
		})
	}
}

func TestDatabaseService_Get(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		serverID       int
		databaseID     int
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:       "Successful get",
			serverID:   1,
			databaseID: 123,
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body: []byte(`{
					"object": "server_database",
					"attributes": {
						"id": 123,
						"server": 1,
						"host": 1,
						"database": "minecraft_db",
						"username": "minecraft_user",
						"remote": "%",
						"max_connections": 10,
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z"
					}
				}`),
			},
			expectedError:  false,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers/1/databases/123",
		},
		{
			name:       "Database not found",
			serverID:   2,
			databaseID: 999,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body: []byte(`{
					"errors": [{
						"code": "NotFoundHttpException",
						"status": "404",
						"detail": "The requested database could not be located."
					}]
				}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers/2/databases/999",
		},
		{
			name:       "Server not found",
			serverID:   999,
			databaseID: 123,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body: []byte(`{
					"errors": [{
						"code": "NotFoundHttpException",
						"status": "404",
						"detail": "The requested server could not be located."
					}]
				}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/servers/999/databases/123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{
				Responses: []testutil.MockResponse{tc.mockResponse},
			}

			service := newDatabaseService(mock, tc.serverID)

			database, err := service.Get(context.Background(), tc.databaseID)

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

			// Verify database data
			if database == nil {
				t.Error("expected database to be non-nil")
			} else {
				if database.ID != tc.databaseID {
					t.Errorf("expected database ID %d, got %d", tc.databaseID, database.ID)
				}
				if database.DatabaseName == "" {
					t.Error("expected database name to be non-empty")
				}
			}
		})
	}
}

func TestDatabaseService_Create(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		serverID      int
		options       api.DatabaseCreateOptions
		mockResponse  testutil.MockResponse
		expectedError bool
		expectedBody  string
	}{
		{
			name:     "Successful creation with host",
			serverID: 1,
			options: api.DatabaseCreateOptions{
				DatabaseName: "new_database",
				Remote:       "192.168.1.100",
				Host:         &[]int{1}[0],
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body: []byte(`{
					"object": "server_database",
					"attributes": {
						"id": 456,
						"server": 1,
						"host": 1,
						"database": "new_database",
						"username": "s1_new_database",
						"remote": "192.168.1.100",
						"max_connections": null,
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z"
					}
				}`),
			},
			expectedError: false,
			expectedBody:  `{"database":"new_database","remote":"192.168.1.100","host":1}`,
		},
		{
			name:     "Successful creation without host",
			serverID: 2,
			options: api.DatabaseCreateOptions{
				DatabaseName: "auto_host_db",
				Remote:       "%",
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body: []byte(`{
					"object": "server_database",
					"attributes": {
						"id": 789,
						"server": 2,
						"host": 2,
						"database": "auto_host_db",
						"username": "s2_auto_host_db",
						"remote": "%",
						"max_connections": null,
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z"
					}
				}`),
			},
			expectedError: false,
			expectedBody:  `{"database":"auto_host_db","remote":"%"}`,
		},
		{
			name:     "API error response",
			serverID: 3,
			options: api.DatabaseCreateOptions{
				DatabaseName: "invalid_db",
				Remote:       "invalid_remote",
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
			expectedBody:  `{"database":"invalid_db","remote":"invalid_remote"}`,
		},
		{
			name:     "Server not found",
			serverID: 999,
			options: api.DatabaseCreateOptions{
				DatabaseName: "test_db",
				Remote:       "%",
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body: []byte(`{
					"errors": [{
						"code": "NotFoundHttpException",
						"status": "404",
						"detail": "The requested server could not be located."
					}]
				}`),
			},
			expectedError: true,
			expectedBody:  `{"database":"test_db","remote":"%"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{
				Responses: []testutil.MockResponse{tc.mockResponse},
			}

			service := newDatabaseService(mock, tc.serverID)

			database, err := service.Create(context.Background(), tc.options)

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

			expectedPath := fmt.Sprintf("/api/application/servers/%d/databases", tc.serverID)
			if req.Endpoint != expectedPath {
				t.Errorf("expected path %s, got %s", expectedPath, req.Endpoint)
			}

			// Check request Body
			bodyStr := strings.TrimSpace(string(req.Body))
			if bodyStr != tc.expectedBody {
				t.Errorf("expected Body %s, got %s", tc.expectedBody, bodyStr)
			}

			// Verify database data
			if database == nil {
				t.Error("expected database to be non-nil")
			} else {
				if database.DatabaseName != tc.options.DatabaseName {
					t.Errorf("expected database name %s, got %s", tc.options.DatabaseName, database.DatabaseName)
				}
				if database.Remote != tc.options.Remote {
					t.Errorf("expected remote %s, got %s", tc.options.Remote, database.Remote)
				}
			}
		})
	}
}

func TestDatabaseService_ResetPassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		serverID       int
		databaseID     int
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:       "Successful password reset",
			serverID:   1,
			databaseID: 123,
			mockResponse: testutil.MockResponse{
				StatusCode: 204,
				Body:       []byte(""),
			},
			expectedError:  false,
			expectedMethod: "POST",
			expectedPath:   "/api/application/servers/1/databases/123/reset-password",
		},
		{
			name:       "Database not found",
			serverID:   2,
			databaseID: 999,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body: []byte(`{
					"errors": [{
						"code": "NotFoundHttpException",
						"status": "404",
						"detail": "The requested database could not be located."
					}]
				}`),
			},
			expectedError:  true,
			expectedMethod: "POST",
			expectedPath:   "/api/application/servers/2/databases/999/reset-password",
		},
		{
			name:       "Server not found",
			serverID:   999,
			databaseID: 123,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body: []byte(`{
					"errors": [{
						"code": "NotFoundHttpException",
						"status": "404",
						"detail": "The requested server could not be located."
					}]
				}`),
			},
			expectedError:  true,
			expectedMethod: "POST",
			expectedPath:   "/api/application/servers/999/databases/123/reset-password",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{
				Responses: []testutil.MockResponse{tc.mockResponse},
			}

			service := newDatabaseService(mock, tc.serverID)

			err := service.ResetPassword(context.Background(), tc.databaseID)

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

			// Verify no Body was sent
			if len(req.Body) > 0 {
				t.Error("expected no request Body for password reset")
			}
		})
	}
}

func TestDatabaseService_Delete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		serverID       int
		databaseID     int
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name:       "Successful deletion",
			serverID:   1,
			databaseID: 123,
			mockResponse: testutil.MockResponse{
				StatusCode: 204,
				Body:       []byte(""),
			},
			expectedError:  false,
			expectedMethod: "DELETE",
			expectedPath:   "/api/application/servers/1/databases/123",
		},
		{
			name:       "Database not found",
			serverID:   2,
			databaseID: 999,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body: []byte(`{
					"errors": [{
						"code": "NotFoundHttpException",
						"status": "404",
						"detail": "The requested database could not be located."
					}]
				}`),
			},
			expectedError:  true,
			expectedMethod: "DELETE",
			expectedPath:   "/api/application/servers/2/databases/999",
		},
		{
			name:       "Database in use",
			serverID:   3,
			databaseID: 456,
			mockResponse: testutil.MockResponse{
				StatusCode: 422,
				Body: []byte(`{
					"errors": [{
						"code": "ValidationHttpException",
						"status": "422",
						"detail": "Cannot delete database that is currently in use."
					}]
				}`),
			},
			expectedError:  true,
			expectedMethod: "DELETE",
			expectedPath:   "/api/application/servers/3/databases/456",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{
				Responses: []testutil.MockResponse{tc.mockResponse},
			}

			service := newDatabaseService(mock, tc.serverID)

			err := service.Delete(context.Background(), tc.databaseID)

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

func TestDatabaseService_Integration(t *testing.T) {
	t.Parallel()

	// Test that the service correctly handles the serverID in all operations
	serverID := 42
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
			// Get response
			{
				StatusCode: 200,
				Body: []byte(`{
					"object": "server_database",
					"attributes": {
						"id": 123,
						"server": 42,
						"host": 1,
						"database": "test_db",
						"username": "test_user",
						"remote": "%",
						"max_connections": null,
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z"
					}
				}`),
			},
			// Create response
			{
				StatusCode: 200,
				Body: []byte(`{
					"object": "server_database",
					"attributes": {
						"id": 456,
						"server": 42,
						"host": 1,
						"database": "new_db",
						"username": "new_user",
						"remote": "%",
						"max_connections": null,
						"created_at": "2023-01-01T00:00:00Z",
						"updated_at": "2023-01-01T00:00:00Z"
					}
				}`),
			},
			// Reset password response
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

	service := newDatabaseService(mock, serverID)

	// Test List
	_, _, err := service.List(context.Background(), api.PaginationOptions{})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	// Test Get
	_, err = service.Get(context.Background(), 123)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	// Test Create
	createOptions := api.DatabaseCreateOptions{
		DatabaseName: "new_db",
		Remote:       "%",
	}
	_, err = service.Create(context.Background(), createOptions)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Test ResetPassword
	err = service.ResetPassword(context.Background(), 123)
	if err != nil {
		t.Fatalf("ResetPassword failed: %v", err)
	}

	// Test Delete
	err = service.Delete(context.Background(), 123)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify all Requests used the correct serverID
	expectedBasePath := fmt.Sprintf("/api/application/servers/%d/databases", serverID)
	for i, req := range mock.Requests {
		if !strings.HasPrefix(req.Endpoint, expectedBasePath) {
			t.Errorf("request %d: expected Endpoint to start with %s, got %s", i, expectedBasePath, req.Endpoint)
		}
	}

	// Verify we made exactly 5 Requests
	if len(mock.Requests) != 5 {
		t.Errorf("expected 5 Requests, got %d", len(mock.Requests))
	}
}

func TestDatabaseService_EdgeCases(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		serverID      int
		databaseID    int
		options       api.DatabaseCreateOptions
		mockResponse  testutil.MockResponse
		expectedError bool
		description   string
	}{
		{
			name:       "Zero server ID",
			serverID:   0,
			databaseID: 123,
			mockResponse: testutil.MockResponse{
				StatusCode: 204,
				Body:       []byte(""),
			},
			expectedError: false,
			description:   "Should handle zero server ID gracefully",
		},
		{
			name:       "Negative database ID",
			serverID:   1,
			databaseID: -1,
			mockResponse: testutil.MockResponse{
				StatusCode: 400,
				Body: []byte(`{
					"errors": [{
						"code": "BadRequestHttpException",
						"status": "400",
						"detail": "Invalid database ID"
					}]
				}`),
			},
			expectedError: true,
			description:   "Should handle negative database ID with error",
		},
		{
			name:     "Empty database name",
			serverID: 1,
			options: api.DatabaseCreateOptions{
				DatabaseName: "",
				Remote:       "%",
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 422,
				Body: []byte(`{
					"errors": [{
						"code": "ValidationHttpException",
						"status": "422",
						"detail": "Database name is required"
					}]
				}`),
			},
			expectedError: true,
			description:   "Should handle empty database name with validation error",
		},
		{
			name:     "Empty remote",
			serverID: 1,
			options: api.DatabaseCreateOptions{
				DatabaseName: "test_db",
				Remote:       "",
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 422,
				Body: []byte(`{
					"errors": [{
						"code": "ValidationHttpException",
						"status": "422",
						"detail": "Remote access pattern is required"
					}]
				}`),
			},
			expectedError: true,
			description:   "Should handle empty remote with validation error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{
				Responses: []testutil.MockResponse{tc.mockResponse},
			}

			service := newDatabaseService(mock, tc.serverID)

			var err error
			if tc.options.DatabaseName != "" || tc.options.Remote != "" {
				// Test Create
				_, err = service.Create(context.Background(), tc.options)
			} else {
				// Test Delete
				err = service.Delete(context.Background(), tc.databaseID)
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

func TestDatabaseService_DataValidation(t *testing.T) {
	t.Parallel()

	// Test that database data is properly parsed
	mock := &testutil.MockRequester{
		Responses: []testutil.MockResponse{
			{
				StatusCode: 200,
				Body: []byte(`{
					"object": "server_database",
					"attributes": {
						"id": 123,
						"server": 1,
						"host": 2,
						"database": "test_database",
						"username": "test_username",
						"remote": "192.168.1.0/24",
						"max_connections": 50,
						"created_at": "2023-01-01T12:00:00Z",
						"updated_at": "2023-01-02T12:00:00Z"
					}
				}`),
			},
		},
	}

	service := newDatabaseService(mock, 1)

	database, err := service.Get(context.Background(), 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify all fields are correctly parsed
	if database.ID != 123 {
		t.Errorf("expected ID 123, got %d", database.ID)
	}
	if database.ServerID != 1 {
		t.Errorf("expected ServerID 1, got %d", database.ServerID)
	}
	if database.HostID != 2 {
		t.Errorf("expected HostID 2, got %d", database.HostID)
	}
	if database.DatabaseName != "test_database" {
		t.Errorf("expected DatabaseName 'test_database', got '%s'", database.DatabaseName)
	}
	if database.Username != "test_username" {
		t.Errorf("expected Username 'test_username', got '%s'", database.Username)
	}
	if database.Remote != "192.168.1.0/24" {
		t.Errorf("expected Remote '192.168.1.0/24', got '%s'", database.Remote)
	}
	if database.MaxConnections == nil {
		t.Error("expected MaxConnections to be non-nil")
	} else if *database.MaxConnections != 50 {
		t.Errorf("expected MaxConnections 50, got %d", *database.MaxConnections)
	}

	// Verify timestamps are parsed correctly
	expectedCreatedAt, _ := time.Parse(time.RFC3339, "2023-01-01T12:00:00Z")
	if !database.CreatedAt.Equal(expectedCreatedAt) {
		t.Errorf("expected CreatedAt %v, got %v", expectedCreatedAt, database.CreatedAt)
	}

	expectedUpdatedAt, _ := time.Parse(time.RFC3339, "2023-01-02T12:00:00Z")
	if !database.UpdatedAt.Equal(expectedUpdatedAt) {
		t.Errorf("expected UpdatedAt %v, got %v", expectedUpdatedAt, database.UpdatedAt)
	}
}
