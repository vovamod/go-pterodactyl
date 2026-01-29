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

func TestNodesService_List(t *testing.T) {
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
				Body: []byte(`{"object": "list", "data": [
					{"object": "node", "attributes": {"id": 1, "uuid": "uuid-1", "public": true, "name": "Node1", "location_id": 1, "fqdn": "node1.example.com", "scheme": "https", "memory": 2048, "memory_overallocate": 0, "disk": 10000, "disk_overallocate": 0, "daemon_listen": 8080, "daemon_sftp": 2022, "daemon_base": "/srv/daemon", "maintenance_mode": false, "upload_size": 100, "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-01T00:00:00Z"}}], "meta": {"pagination": {"total": 1, "count": 1, "per_page": 10, "current_page": 1, "total_pages": 1}}}`),
			},
			expectedError:  false,
			expectedCount:  1,
			expectedMethod: "GET",
			expectedPath:   "/api/application/nodes",
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
			expectedPath:   "/api/application/nodes",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewNodesService(mock)
			nodes, meta, err := service.List(context.Background(), tc.options)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(nodes) != tc.expectedCount {
				t.Errorf("expected %d nodes, got %d", tc.expectedCount, len(nodes))
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
			if tc.expectedCount > 0 && len(nodes) > 0 {
				node := nodes[0]
				if node.ID == 0 {
					t.Error("expected node ID to be non-zero")
				}
				if node.Name == "" {
					t.Error("expected node name to be non-empty")
				}
			}
		})
	}
}

func TestNodesService_ListAll(t *testing.T) {
	t.Parallel()
	mock := &testutil.MockRequester{Responses: []testutil.MockResponse{{StatusCode: 200, Body: []byte(`{"object": "list", "data": [{"object": "node", "attributes": {"id": 1, "uuid": "uuid-1", "public": true, "name": "Node1", "location_id": 1, "fqdn": "node1.example.com", "scheme": "https", "memory": 2048, "memory_overallocate": 0, "disk": 10000, "disk_overallocate": 0, "daemon_listen": 8080, "daemon_sftp": 2022, "daemon_base": "/srv/daemon", "maintenance_mode": false, "upload_size": 100, "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-01T00:00:00Z"}}], "meta": {"pagination": {"total": 1, "count": 1, "per_page": 100, "current_page": 1, "total_pages": 1}}}}`)}}}
	service := NewNodesService(mock)
	nodes, err := service.ListAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(nodes))
	}
	if nodes[0].Name != "Node1" {
		t.Errorf("expected name 'Node1', got '%s'", nodes[0].Name)
	}
}

func TestNodesService_Get(t *testing.T) {
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
				Body:       []byte(`{"object": "node", "attributes": {"id": 1, "uuid": "uuid-1", "public": true, "name": "Node1", "location_id": 1, "fqdn": "node1.example.com", "scheme": "https", "memory": 2048, "memory_overallocate": 0, "disk": 10000, "disk_overallocate": 0, "daemon_listen": 8080, "daemon_sftp": 2022, "daemon_base": "/srv/daemon", "maintenance_mode": false, "upload_size": 100, "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-01T00:00:00Z"}}`),
			},
			expectedError:  false,
			expectedMethod: "GET",
			expectedPath:   "/api/application/nodes/1",
		},
		{
			name: "Node not found",
			id:   999,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "GET",
			expectedPath:   "/api/application/nodes/999",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewNodesService(mock)
			node, err := service.Get(context.Background(), tc.id)
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
			if node == nil {
				t.Error("expected node to be non-nil")
			} else {
				if node.ID != tc.id {
					t.Errorf("expected node ID %d, got %d", tc.id, node.ID)
				}
				if node.Name == "" {
					t.Error("expected node name to be non-empty")
				}
			}
		})
	}
}

func TestNodesService_GetConfiguration(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		nodeID        int
		mockResponse  testutil.MockResponse
		expectedError bool
		expectedPath  string
	}{
		{
			name:   "Successful get configuration",
			nodeID: 1,
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body:       []byte(`{"debug":true,"uuid":"uuid-1","token_id":"tid","token":"tok","api":{"host":"127.0.0.1","port":8080,"ssl":{"enabled":true,"cert":"/cert","key":"/key"},"upload_limit":100},"system":{"data":"/srv/daemon","sftp":{"bind_port":2022}},"remote":"https://remote"}`),
			},
			expectedError: false,
			expectedPath:  "/api/application/nodes/1/configuration",
		},
		{
			name:   "Node not found",
			nodeID: 999,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Not found."}]}`),
			},
			expectedError: true,
			expectedPath:  "/api/application/nodes/999/configuration",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewNodesService(mock)
			config, err := service.GetConfiguration(context.Background(), tc.nodeID)
			if tc.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if config == nil {
				t.Error("expected config to be non-nil")
			}
			if len(mock.Requests) != 1 {
				t.Fatalf("expected 1 request, got %d", len(mock.Requests))
			}
			req := mock.Requests[0]
			if req.Endpoint != tc.expectedPath {
				t.Errorf("expected path %s, got %s", tc.expectedPath, req.Endpoint)
			}
		})
	}
}

func TestNodesService_Create(t *testing.T) {
	t.Parallel()
	desc := "desc"
	behindProxy := true
	maint := false
	testCases := []struct {
		name          string
		options       api.NodeCreateOptions
		mockResponse  testutil.MockResponse
		expectedError bool
		expectedBody  string
	}{
		{
			name: "Successful creation",
			options: api.NodeCreateOptions{
				Name: "Node1", LocationID: 1, FQDN: "node1.example.com", Scheme: "https", Memory: 2048, MemoryOverallocate: 0, Disk: 10000, DiskOverallocate: 0, DaemonSFTP: 2022, DaemonListen: 8080, Description: &desc, BehindProxy: &behindProxy, MaintenanceMode: &maint, UploadSize: nil,
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body:       []byte(`{"object": "node", "attributes": {"id": 1, "uuid": "uuid-1", "public": true, "name": "Node1", "location_id": 1, "fqdn": "node1.example.com", "scheme": "https", "memory": 2048, "memory_overallocate": 0, "disk": 10000, "disk_overallocate": 0, "daemon_listen": 8080, "daemon_sftp": 2022, "daemon_base": "/srv/daemon", "maintenance_mode": false, "upload_size": 100, "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-01T00:00:00Z"}}`),
			},
			expectedError: false,
			expectedBody:  `{"name":"Node1","location_id":1,"fqdn":"node1.example.com","scheme":"https","memory":2048,"memory_overallocate":0,"disk":10000,"disk_overallocate":0,"daemon_sftp":2022,"daemon_listen":8080,"description":"desc","behind_proxy":true,"maintenance_mode":false}`,
		},
		{
			name:    "API error response",
			options: api.NodeCreateOptions{Name: "", LocationID: 0, FQDN: "", Scheme: "", Memory: 0, MemoryOverallocate: 0, Disk: 0, DiskOverallocate: 0, DaemonSFTP: 0, DaemonListen: 0},
			mockResponse: testutil.MockResponse{
				StatusCode: 422,
				Body:       []byte(`{"errors": [{"code": "ValidationHttpException", "status": "422", "detail": "Invalid data."}]}`),
			},
			expectedError: true,
			expectedBody:  `{"name":"","location_id":0,"fqdn":"","scheme":"","memory":0,"memory_overallocate":0,"disk":0,"disk_overallocate":0,"daemon_sftp":0,"daemon_listen":0}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewNodesService(mock)
			node, err := service.Create(context.Background(), tc.options)
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
			if req.Endpoint != "/api/application/nodes" {
				t.Errorf("expected path /api/application/nodes, got %s", req.Endpoint)
			}
			bodyStr := strings.TrimSpace(string(req.Body))
			if bodyStr != tc.expectedBody {
				t.Errorf("expected Body %s, got %s", tc.expectedBody, bodyStr)
			}
			if node == nil {
				t.Error("expected node to be non-nil")
			}
		})
	}
}

func TestNodesService_Update(t *testing.T) {
	t.Parallel()
	desc := "desc2"
	maint := true
	testCases := []struct {
		name          string
		id            int
		options       api.NodeUpdateOptions
		mockResponse  testutil.MockResponse
		expectedError bool
		expectedBody  string
	}{
		{
			name:    "Successful update",
			id:      1,
			options: api.NodeUpdateOptions{Name: "Node2", Description: &desc, MaintenanceMode: &maint},
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body:       []byte(`{"object": "node", "attributes": {"id": 1, "uuid": "uuid-1", "public": true, "name": "Node2", "location_id": 1, "fqdn": "node2.example.com", "scheme": "https", "memory": 4096, "memory_overallocate": 0, "disk": 20000, "disk_overallocate": 0, "daemon_listen": 8081, "daemon_sftp": 2023, "daemon_base": "/srv/daemon", "maintenance_mode": true, "upload_size": 200, "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-02T00:00:00Z"}}`),
			},
			expectedError: false,
			expectedBody:  `{"name":"Node2","description":"desc2","maintenance_mode":true}`,
		},
		{
			name:    "API error response",
			id:      2,
			options: api.NodeUpdateOptions{Name: ""},
			mockResponse: testutil.MockResponse{
				StatusCode: 422,
				Body:       []byte(`{"errors": [{"code": "ValidationHttpException", "status": "422", "detail": "Invalid data."}]}`),
			},
			expectedError: true,
			expectedBody:  `{"name":""}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewNodesService(mock)
			node, err := service.Update(context.Background(), tc.id, tc.options)
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
			expectedPath := fmt.Sprintf("/api/application/nodes/%d", tc.id)
			if req.Endpoint != expectedPath {
				t.Errorf("expected path %s, got %s", expectedPath, req.Endpoint)
			}
			bodyStr := strings.TrimSpace(string(req.Body))
			if bodyStr != tc.expectedBody {
				t.Errorf("expected Body %s, got %s", tc.expectedBody, bodyStr)
			}
			if node == nil {
				t.Error("expected node to be non-nil")
			}
		})
	}
}

func TestNodesService_Delete(t *testing.T) {
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
			expectedPath:   "/api/application/nodes/1",
		},
		{
			name: "Node not found",
			id:   999,
			mockResponse: testutil.MockResponse{
				StatusCode: 404,
				Body:       []byte(`{"errors": [{"code": "NotFoundHttpException", "status": "404", "detail": "Not found."}]}`),
			},
			expectedError:  true,
			expectedMethod: "DELETE",
			expectedPath:   "/api/application/nodes/999",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &testutil.MockRequester{Responses: []testutil.MockResponse{tc.mockResponse}}
			service := NewNodesService(mock)
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
			expectedPath := fmt.Sprintf("/api/application/nodes/%d", tc.id)
			if req.Endpoint != expectedPath {
				t.Errorf("expected path %s, got %s", expectedPath, req.Endpoint)
			}
		})
	}
}

func TestNodesService_Allocations(t *testing.T) {
	t.Parallel()
	mock := &testutil.MockRequester{}
	service := NewNodesService(mock)
	allocService := service.Allocations(context.Background(), 42)
	if allocService == nil {
		t.Fatal("expected allocations service to be non-nil")
	}
}

func TestNodesService_DataValidation(t *testing.T) {
	t.Parallel()
	mock := &testutil.MockRequester{Responses: []testutil.MockResponse{{StatusCode: 200, Body: []byte(`{"object": "node", "attributes": {"id": 1, "uuid": "uuid-1", "public": true, "name": "Node1", "location_id": 1, "fqdn": "node1.example.com", "scheme": "https", "memory": 2048, "memory_overallocate": 0, "disk": 10000, "disk_overallocate": 0, "daemon_listen": 8080, "daemon_sftp": 2022, "daemon_base": "/srv/daemon", "maintenance_mode": false, "upload_size": 100, "created_at": "2023-01-01T12:00:00Z", "updated_at": "2023-01-02T12:00:00Z"}}`)}}}
	service := NewNodesService(mock)
	node, err := service.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.ID != 1 {
		t.Errorf("expected ID 1, got %d", node.ID)
	}
	if node.UUID != "uuid-1" {
		t.Errorf("expected UUID 'uuid-1', got '%s'", node.UUID)
	}
	if node.Name != "Node1" {
		t.Errorf("expected Name 'Node1', got '%s'", node.Name)
	}
	if node.FQDN != "node1.example.com" {
		t.Errorf("expected FQDN 'node1.example.com', got '%s'", node.FQDN)
	}
	expectedCreatedAt, _ := time.Parse(time.RFC3339, "2023-01-01T12:00:00Z")
	if !node.CreatedAt.Equal(expectedCreatedAt) {
		t.Errorf("expected CreatedAt %v, got %v", expectedCreatedAt, node.CreatedAt)
	}
	expectedUpdatedAt, _ := time.Parse(time.RFC3339, "2023-01-02T12:00:00Z")
	if !node.UpdatedAt.Equal(expectedUpdatedAt) {
		t.Errorf("expected UpdatedAt %v, got %v", expectedUpdatedAt, node.UpdatedAt)
	}
}
