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

func TestServerService_GetDetails(t *testing.T) {
	expectedServer := &api.ClientServer{
		ServerOwner: true,
		Identifier:  testServerIdentifier,
		UUID:        "server-uuid",
		Name:        "Test Server",
	}
	res := api.ListItem[api.ClientServer]{Object: "server", Attributes: expectedServer}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newServerService(mock, testServerIdentifier)
		server, err := s.GetDetails(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(server, expectedServer) {
			t.Errorf("expected server %+v, got %+v", expectedServer, server)
		}
		req := mock.Requests[0]
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s", testServerIdentifier)
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
	})
}

func TestServerService_GetWebsocket(t *testing.T) {
	expectedDetails := &api.WebsocketDetails{
		Token:     "websocket-token",
		SocketURL: "wss://example.com/socket",
	}
	res := api.WebsocketResponse{Data: *expectedDetails}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newServerService(mock, testServerIdentifier)
		details, err := s.GetWebsocket(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(details, expectedDetails) {
			t.Errorf("expected details %+v, got %+v", expectedDetails, details)
		}
	})
}

func TestServerService_GetResources(t *testing.T) {
	expectedResources := &api.Resources{
		Object: "stats",
		Attributes: api.ResourceAttributes{
			CurrentState: "running",
		},
	}
	jsonBody, _ := json.Marshal(expectedResources)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newServerService(mock, testServerIdentifier)
		resources, err := s.GetResources(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(resources, expectedResources) {
			t.Errorf("expected resources %+v, got %+v", expectedResources, resources)
		}
	})
}

func TestServerService_SendCommand(t *testing.T) {
	command := "say hello"
	options := api.SendCommandOptions{Command: command}
	jsonBody, _ := json.Marshal(options)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusNoContent}},
		}
		s := newServerService(mock, testServerIdentifier)
		err := s.SendCommand(context.Background(), command)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		req := mock.Requests[0]
		if req.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", req.Method)
		}
		if !bytes.Equal(req.Body, jsonBody) {
			t.Errorf("expected body %s, got %s", jsonBody, req.Body)
		}
	})
}

func TestServerService_SetPowerState(t *testing.T) {
	signal := "start"
	options := api.SetPowerStateOptions{Signal: signal}
	jsonBody, _ := json.Marshal(options)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusNoContent}},
		}
		s := newServerService(mock, testServerIdentifier)
		err := s.SetPowerState(context.Background(), signal)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		req := mock.Requests[0]
		if !bytes.Equal(req.Body, jsonBody) {
			t.Errorf("expected body %s, got %s", jsonBody, req.Body)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{Err: &errors.APIError{HTTPStatusCode: http.StatusConflict}}},
		}
		s := newServerService(mock, testServerIdentifier)
		err := s.SetPowerState(context.Background(), "stop")
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}

func TestServerService_SubServices(t *testing.T) {
	mock := &testutil.MockRequester{}
	s := newServerService(mock, testServerIdentifier)

	if s.Databases() == nil {
		t.Error("Databases() returned nil")
	}
	if s.Files() == nil {
		t.Error("Files() returned nil")
	}
	if s.Schedules() == nil {
		t.Error("Schedules() returned nil")
	}
	if s.Network() == nil {
		t.Error("Network() returned nil")
	}
	if s.Users() == nil {
		t.Error("Users() returned nil")
	}
	if s.Backups() == nil {
		t.Error("Backups() returned nil")
	}
	if s.Startup() == nil {
		t.Error("Startup() returned nil")
	}
	if s.Settings() == nil {
		t.Error("Settings() returned nil")
	}
}
