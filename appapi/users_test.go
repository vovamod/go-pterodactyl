package appapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/vovamod/go-pterodactyl/internal/testutil"
	"strings"
	"testing"

	"github.com/vovamod/go-pterodactyl/api"
)

func TestUsersService_List(t *testing.T) {
	t.Parallel()

	mockUserListResponse := `{
		"object": "list",
		"data": [
			{
				"object": "user",
				"attributes": { "id": 1, "username": "testuser1", "email": "test1@example.com" }
			},
			{
				"object": "user",
				"attributes": { "id": 2, "username": "testuser2", "email": "test2@example.com" }
			}
		],
		"meta": {
			"pagination": { "total": 2, "count": 2, "per_page": 10, "current_page": 1, "total_pages": 1 }
		}
	}`

	testCases := []struct {
		name          string
		options       *api.PaginationOptions
		mockResponse  testutil.MockResponse
		expectedError bool
		expectedCount int
	}{
		{
			name:    "Successful list",
			options: &api.PaginationOptions{Page: 1},
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body:       []byte(mockUserListResponse),
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name:    "API error",
			options: &api.PaginationOptions{Page: 1},
			mockResponse: testutil.MockResponse{
				StatusCode: 500,
				Body:       []byte(`{"errors": [{"code": "ServerException"}]}`),
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewUsersService(mock)

			users, meta, err := service.List(context.Background(), tc.options)

			if tc.expectedError {
				if err == nil {
					t.Fatal("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(users) != tc.expectedCount {
				t.Errorf("expected %d users, got %d", tc.expectedCount, len(users))
			}
			if meta == nil {
				t.Error("expected meta to be non-nil")
			}
			if mock.Requests[0].Endpoint != "/api/application/users" {
				t.Errorf("expected Endpoint '/api/application/users', got '%s'", mock.Requests[0].Endpoint)
			}
		})
	}
}

func TestUsersService_Get(t *testing.T) {
	t.Parallel()

	mockUserResponse := `{
		"object": "user",
		"attributes": { "id": 123, "username": "getuser", "email": "get@example.com" }
	}`

	testCases := []struct {
		name          string
		userID        int
		mockResponse  testutil.MockResponse
		expectedError bool
		expectedID    int
	}{
		{
			name:   "Successful get",
			userID: 123,
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body:       []byte(mockUserResponse),
			},
			expectedError: false,
			expectedID:    123,
		},
		{
			name:   "User not found",
			userID: 999,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body:       []byte(`{"errors": [{"code": "NotFoundHttpException"}]}`),
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewUsersService(mock)

			user, err := service.Get(context.Background(), tc.userID)

			if tc.expectedError {
				if err == nil {
					t.Fatal("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if user.ID != tc.expectedID {
				t.Errorf("expected user ID %d, got %d", tc.expectedID, user.ID)
			}

			expectedPath := fmt.Sprintf("/api/application/users/%d", tc.userID)
			if mock.Requests[0].Endpoint != expectedPath {
				t.Errorf("expected Endpoint '%s', got '%s'", expectedPath, mock.Requests[0].Endpoint)
			}
		})
	}
}

func TestUsersService_GetExternalID(t *testing.T) {
	t.Parallel()

	mockUserResponse := `{
		"object": "user",
		"attributes": { "id": 456, "external_id": "ext-123", "username": "extuser" }
	}`

	testCases := []struct {
		name           string
		externalID     string
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedUserID int
	}{
		{
			name:       "Successful get by external ID",
			externalID: "ext-123",
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body:       []byte(mockUserResponse),
			},
			expectedError:  false,
			expectedUserID: 456,
		},
		{
			name:       "User not found by external ID",
			externalID: "ext-999",
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body:       []byte(`{"errors": [{"code": "NotFoundHttpException"}]}`),
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewUsersService(mock)

			user, err := service.GetExternalID(context.Background(), tc.externalID)

			if tc.expectedError {
				if err == nil {
					t.Fatal("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if user.ID != tc.expectedUserID {
				t.Errorf("expected user ID %d, got %d", tc.expectedUserID, user.ID)
			}

			expectedPath := fmt.Sprintf("/api/application/users/external/%s", tc.externalID)
			if mock.Requests[0].Endpoint != expectedPath {
				t.Errorf("expected Endpoint '%s', got '%s'", expectedPath, mock.Requests[0].Endpoint)
			}
		})
	}
}

func TestUsersService_Create(t *testing.T) {
	t.Parallel()

	options := api.UserCreateOptions{
		Username:  "newuser",
		Email:     "new@example.com",
		FirstName: "New",
		LastName:  "User",
		Password:  "supersecret",
	}

	mockCreatedUserResponse := `{
		"object": "user",
		"attributes": { "id": 1, "username": "newuser", "email": "new@example.com" }
	}`

	expectedBody, _ := json.Marshal(options)

	testCases := []struct {
		name          string
		options       api.UserCreateOptions
		mockResponse  testutil.MockResponse
		expectedError bool
		expectedBody  string
	}{
		{
			name:    "Successful creation",
			options: options,
			mockResponse: testutil.MockResponse{
				StatusCode: 201, // 201 Created is typical
				Body:       []byte(mockCreatedUserResponse),
			},
			expectedError: false,
			expectedBody:  string(expectedBody),
		},
		{
			name:    "Validation error",
			options: options,
			mockResponse: testutil.MockResponse{
				StatusCode: 422,
				Body:       []byte(`{"errors": [{"code": "ValidationException"}]}`),
			},
			expectedError: true,
			expectedBody:  string(expectedBody),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewUsersService(mock)

			user, err := service.Create(context.Background(), tc.options)

			if tc.expectedError {
				if err == nil {
					t.Fatal("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if user == nil {
				t.Fatal("expected a user object but got nil")
			}
			if user.Username != tc.options.Username {
				t.Errorf("expected username '%s', got '%s'", tc.options.Username, user.Username)
			}

			req := mock.Requests[0]
			if req.Method != "POST" {
				t.Errorf("expected Method POST, got %s", req.Method)
			}
			if req.Endpoint != "/api/application/users" {
				t.Errorf("expected Endpoint '/api/application/users', got '%s'", req.Endpoint)
			}
			if strings.TrimSpace(string(req.Body)) != tc.expectedBody {
				t.Errorf("expected Body '%s', got '%s'", tc.expectedBody, string(req.Body))
			}
		})
	}
}

func TestUsersService_Update(t *testing.T) {
	t.Parallel()

	userID := 123
	options := api.UserUpdateOptions{
		Username:  "updateduser",
		FirstName: "Updated",
	}

	mockUpdatedUserResponse := `{
		"object": "user",
		"attributes": { "id": 123, "username": "updateduser", "first_name": "Updated" }
	}`

	expectedBody, _ := json.Marshal(options)

	testCases := []struct {
		name          string
		userID        int
		options       api.UserUpdateOptions
		mockResponse  testutil.MockResponse
		expectedError bool
		expectedBody  string
	}{
		{
			name:    "Successful update",
			userID:  userID,
			options: options,
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body:       []byte(mockUpdatedUserResponse),
			},
			expectedError: false,
			expectedBody:  string(expectedBody),
		},
		{
			name:    "User not found",
			userID:  999,
			options: options,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body:       []byte(`{"errors": [{"code": "NotFoundHttpException"}]}`),
			},
			expectedError: true,
			expectedBody:  string(expectedBody),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewUsersService(mock)

			user, err := service.Update(context.Background(), tc.userID, tc.options)

			if tc.expectedError {
				if err == nil {
					t.Fatal("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if user.Username != tc.options.Username {
				t.Errorf("expected username '%s', got '%s'", tc.options.Username, user.Username)
			}

			req := mock.Requests[0]
			if req.Method != "PATCH" {
				t.Errorf("expected Method PATCH, got %s", req.Method)
			}
			expectedPath := fmt.Sprintf("/api/application/users/%d", tc.userID)
			if req.Endpoint != expectedPath {
				t.Errorf("expected Endpoint '%s', got '%s'", expectedPath, req.Endpoint)
			}
			if strings.TrimSpace(string(req.Body)) != tc.expectedBody {
				t.Errorf("expected Body '%s', got '%s'", tc.expectedBody, string(req.Body))
			}
		})
	}
}

func TestUsersService_Delete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		userID        int
		mockResponse  testutil.MockResponse
		expectedError bool
	}{
		{
			name:   "Successful deletion",
			userID: 123,
			mockResponse: testutil.MockResponse{
				StatusCode: 204, // No Content
				Body:       []byte(""),
			},
			expectedError: false,
		},
		{
			name:   "User not found",
			userID: 999,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body:       []byte(`{"errors": [{"code": "NotFoundHttpException"}]}`),
			},
			expectedError: true,
		},
		{
			name:   "Cannot delete user with servers",
			userID: 456,
			mockResponse: testutil.MockResponse{
				StatusCode: 400, // Bad Request
				Body:       []byte(`{"errors": [{"code": "HasActiveServersException"}]}`),
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewUsersService(mock)

			err := service.Delete(context.Background(), tc.userID)

			if tc.expectedError {
				if err == nil {
					t.Fatal("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			req := mock.Requests[0]
			if req.Method != "DELETE" {
				t.Errorf("expected Method DELETE, got %s", req.Method)
			}
			expectedPath := fmt.Sprintf("/api/application/users/%d", tc.userID)
			if req.Endpoint != expectedPath {
				t.Errorf("expected Endpoint '%s', got '%s'", expectedPath, req.Endpoint)
			}
		})
	}
}
