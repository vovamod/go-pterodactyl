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

func TestLocationService_List(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		options        *api.PaginationOptions
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedCount  int
		expectedMethod string
		expectedPath   string
	}{
		{
			name:    "Successful list with pagination",
			options: &api.PaginationOptions{Page: 1, PerPage: 10},
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body: []byte(`{
					"object": "list",
					"data": [
						{"object": "location", "attributes": {"id": 1, "short": "us", "long": "United States", "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-01T00:00:00Z"}},
						{"object": "location", "attributes": {"id": 2, "short": "eu", "long": "Europe", "created_at": "2023-01-02T00:00:00Z", "updated_at": "2023-01-02T00:00:00Z"}}
					],
					"meta": {"pagination": {"total": 2, "count": 2, "per_page": 10, "current_page": 1, "total_pages": 1}}
				}`),
			},
			expectedError:  false,
			expectedCount:  2,
			expectedMethod: "GET",
			expectedPath:   "/api/application/locations",
		},
		{
			name:    "Successful list without pagination",
			options: nil,
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body:       []byte(`{"object": "list", "data": [], "meta": {"pagination": {"total": 0, "count": 0, "per_page": 100, "current_page": 1, "total_pages": 0}}}`),
			},
			expectedError:  false,
			expectedCount:  0,
			expectedMethod: "GET",
			expectedPath:   "/api/application/locations",
		},
		{
			name:    "API error response",
			options: &api.PaginationOptions{Page: 1},
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/locations",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewLocationService(mock)
			locations, meta, err := service.List(context.Background(), tc.options)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(locations) != tc.expectedCount {
				t.Errorf("expected %d locations, got %d", tc.expectedCount, len(locations))
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
			if tc.expectedCount > 0 && len(locations) > 0 {
				loc := locations[0]
				if loc.ID == 0 {
					t.Error("expected location ID to be non-zero")
				}
				if loc.ShortCode == "" {
					t.Error("expected short code to be non-empty")
				}
			}
		})
	}
}

func TestLocationService_ListAll(t *testing.T) {
	t.Parallel()
	mock := &testutil.MockRequester{Responses: []testutil.MockResponse{{StatusCode: 200, Body: []byte(`{"object": "list", "data": [{"object": "location", "attributes": {"id": 1, "short": "us", "long": "United States", "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-01T00:00:00Z"}}], "meta": {"pagination": {"total": 1, "count": 1, "per_page": 100, "current_page": 1, "total_pages": 1}}}}`)}}}
	service := NewLocationService(mock)
	locations, err := service.ListAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(locations) != 1 {
		t.Errorf("expected 1 location, got %d", len(locations))
	}
	if locations[0].ShortCode != "us" {
		t.Errorf("expected short code 'us', got '%s'", locations[0].ShortCode)
	}
}

func TestLocationService_Get(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		id             int
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name: "Successful get",
			id:   1,
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body:       []byte(`{"object": "location", "attributes": {"id": 1, "short": "us", "long": "United States", "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-01T00:00:00Z"}}`),
			},
			expectedError:  false,
			expectedMethod: "GET",
			expectedPath:   "/api/application/locations/1",
		},
		{
			name: "Location not found",
			id:   999,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/locations/999",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewLocationService(mock)
			location, err := service.Get(context.Background(), tc.id)
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
			if location == nil {
				t.Error("expected location to be non-nil")
			} else {
				if location.ID != tc.id {
					t.Errorf("expected location ID %d, got %d", tc.id, location.ID)
				}
				if location.ShortCode == "" {
					t.Error("expected short code to be non-empty")
				}
			}
		})
	}
}

func TestLocationService_Create(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		options       api.LocationCreateOptions
		mockResponse  testutil.MockResponse
		expectedError bool
		expectedBody  string
	}{
		{
			name:    "Successful creation",
			options: api.LocationCreateOptions{ShortCode: "us", Description: "United States"},
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body:       []byte(`{"object": "location", "attributes": {"id": 1, "short": "us", "long": "United States", "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-01T00:00:00Z"}}`),
			},
			expectedError: false,
			expectedBody:  `{"short":"us","long":"United States"}`,
		},
		{
			name:    "API error response",
			options: api.LocationCreateOptions{ShortCode: "", Description: ""},
			mockResponse: testutil.MockResponse{
				StatusCode: 422,
				Body:       []byte(`{"errors": [{"code": "ValidationHttpException", "status": "422", "detail": "Invalid data."}]}`),
			},
			expectedError: true,
			expectedBody:  `{"short":"","long":""}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewLocationService(mock)
			location, err := service.Create(context.Background(), tc.options)
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
			if req.Method != "POST" {
				t.Errorf("expected Method POST, got %s", req.Method)
			}
			if req.Endpoint != "/api/application/locations" {
				t.Errorf("expected path /api/application/locations, got %s", req.Endpoint)
			}
			bodyStr := strings.TrimSpace(string(req.Body))
			if bodyStr != tc.expectedBody {
				t.Errorf("expected Body %s, got %s", tc.expectedBody, bodyStr)
			}
			if location == nil {
				t.Error("expected location to be non-nil")
			} else {
				if location.ShortCode != tc.options.ShortCode {
					t.Errorf("expected short code %s, got %s", tc.options.ShortCode, location.ShortCode)
				}
				if location.Description != tc.options.Description {
					t.Errorf("expected description %s, got %s", tc.options.Description, location.Description)
				}
			}
		})
	}
}

func TestLocationService_Update(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		id            int
		options       api.LocationUpdateOptions
		mockResponse  testutil.MockResponse
		expectedError bool
		expectedBody  string
	}{
		{
			name:    "Successful update",
			id:      1,
			options: api.LocationUpdateOptions{ShortCode: "us", Description: "USA"},
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body:       []byte(`{"object": "location", "attributes": {"id": 1, "short": "us", "long": "USA", "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-02T00:00:00Z"}}`),
			},
			expectedError: false,
			expectedBody:  `{"short":"us","long":"USA"}`,
		},
		{
			name:    "API error response",
			id:      2,
			options: api.LocationUpdateOptions{ShortCode: "", Description: ""},
			mockResponse: testutil.MockResponse{
				StatusCode: 422,
				Body:       []byte(`{"errors": [{"code": "ValidationHttpException", "status": "422", "detail": "Invalid data."}]}`),
			},
			expectedError: true,
			expectedBody:  `{"short":"","long":""}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewLocationService(mock)
			location, err := service.Update(context.Background(), tc.id, tc.options)
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
			if req.Method != "PATCH" {
				t.Errorf("expected Method PATCH, got %s", req.Method)
			}
			expectedPath := fmt.Sprintf("/api/application/locations/%d", tc.id)
			if req.Endpoint != expectedPath {
				t.Errorf("expected path %s, got %s", expectedPath, req.Endpoint)
			}
			bodyStr := strings.TrimSpace(string(req.Body))
			if bodyStr != tc.expectedBody {
				t.Errorf("expected Body %s, got %s", tc.expectedBody, bodyStr)
			}
			if location == nil {
				t.Error("expected location to be non-nil")
			} else {
				if location.ShortCode != tc.options.ShortCode {
					t.Errorf("expected short code %s, got %s", tc.options.ShortCode, location.ShortCode)
				}
				if location.Description != tc.options.Description {
					t.Errorf("expected description %s, got %s", tc.options.Description, location.Description)
				}
			}
		})
	}
}

func TestLocationService_Delete(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		id             int
		mockResponse   testutil.MockResponse
		expectedError  bool
		expectedMethod string
		expectedPath   string
	}{
		{
			name: "Successful deletion",
			id:   1,
			mockResponse: testutil.MockResponse{
				StatusCode: 204,
				Body:       []byte(""),
			},
			expectedError:  false,
			expectedMethod: "DELETE",
			expectedPath:   "/api/application/locations/1",
		},
		{
			name: "Location not found",
			id:   999,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "DELETE",
			expectedPath:   "/api/application/locations/999",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewLocationService(mock)
			err := service.Delete(context.Background(), tc.id)
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
			expectedPath := fmt.Sprintf("/api/application/locations/%d", tc.id)
			if req.Endpoint != expectedPath {
				t.Errorf("expected path %s, got %s", expectedPath, req.Endpoint)
			}
		})
	}
}

func TestLocationService_DataValidation(t *testing.T) {
	t.Parallel()
	mock := &testutil.MockRequester{Responses: []testutil.MockResponse{{StatusCode: 200, Body: []byte(`{"object": "location", "attributes": {"id": 1, "short": "us", "long": "United States", "created_at": "2023-01-01T12:00:00Z", "updated_at": "2023-01-02T12:00:00Z"}}`)}}}
	service := NewLocationService(mock)
	location, err := service.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if location.ID != 1 {
		t.Errorf("expected ID 1, got %d", location.ID)
	}
	if location.ShortCode != "us" {
		t.Errorf("expected ShortCode 'us', got '%s'", location.ShortCode)
	}
	if location.Description != "United States" {
		t.Errorf("expected Description 'United States', got '%s'", location.Description)
	}
	expectedCreatedAt, _ := time.Parse(time.RFC3339, "2023-01-01T12:00:00Z")
	if !location.CreatedAt.Equal(expectedCreatedAt) {
		t.Errorf("expected CreatedAt %v, got %v", expectedCreatedAt, location.CreatedAt)
	}
	expectedUpdatedAt, _ := time.Parse(time.RFC3339, "2023-01-02T12:00:00Z")
	if !location.UpdatedAt.Equal(expectedUpdatedAt) {
		t.Errorf("expected UpdatedAt %v, got %v", expectedUpdatedAt, location.UpdatedAt)
	}
}
