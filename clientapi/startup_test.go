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

func TestStartupService_ListVariables(t *testing.T) {
	expectedVars := []*api.StartupVariable{
		{
			Name:         "Server Jar File",
			Description:  "The JAR file to use for the server.",
			EnvVariable:  "SERVER_JARFILE",
			DefaultValue: "server.jar",
			ServerValue:  "paper.jar",
			IsEditable:   true,
			Rules:        "required|string",
		},
		{
			Name:         "Server Memory",
			Description:  "The amount of memory to allocate.",
			EnvVariable:  "SERVER_MEMORY",
			DefaultValue: "1024",
			ServerValue:  "4096",
			IsEditable:   true,
			Rules:        "required|numeric",
		},
	}
	data := make([]*api.ListItem[api.StartupVariable], len(expectedVars))
	for i, v := range expectedVars {
		data[i] = &api.ListItem[api.StartupVariable]{Object: "startup_variable", Attributes: v}
	}
	meta := api.Meta{Pagination: api.Pagination{Total: 2, PerPage: 25, CurrentPage: 1, TotalPages: 1}}
	res := api.PaginatedResponse[api.StartupVariable]{
		Object: "list",
		Data:   data,
		Meta:   meta,
	}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newStartupService(mock, testServerIdentifier)
		vars, m, err := s.ListVariables(context.Background(), api.PaginationOptions{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(vars, expectedVars) {
			t.Errorf("expected variables %+v, got %+v", expectedVars, vars)
		}
		if !reflect.DeepEqual(m, &meta) {
			t.Errorf("expected meta %+v, got %+v", &meta, m)
		}
		req := mock.Requests[0]
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/startup", testServerIdentifier)
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{Err: &errors.APIError{HTTPStatusCode: http.StatusInternalServerError}}},
		}
		s := newStartupService(mock, testServerIdentifier)
		_, _, err := s.ListVariables(context.Background(), api.PaginationOptions{})
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}

func TestStartupService_UpdateVariable(t *testing.T) {
	options := api.UpdateVariableOptions{Key: "SERVER_MEMORY", Value: "8192"}
	jsonOptions, _ := json.Marshal(options)
	expectedVar := &api.StartupVariable{
		Name:        "Server Memory",
		EnvVariable: "SERVER_MEMORY",
		ServerValue: options.Value,
	}
	res := api.UpdateVariableResponse{Object: "startup_variable", Attributes: expectedVar}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newStartupService(mock, testServerIdentifier)
		variable, err := s.UpdateVariable(context.Background(), options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(variable, expectedVar) {
			t.Errorf("expected variable %+v, got %+v", expectedVar, variable)
		}
		req := mock.Requests[0]
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/startup/variable", testServerIdentifier)
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
		if req.Method != http.MethodPut {
			t.Errorf("expected method PUT, got %s", req.Method)
		}
		if !bytes.Equal(req.Body, jsonOptions) {
			t.Errorf("expected body %s, got %s", jsonOptions, req.Body)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{Err: &errors.APIError{HTTPStatusCode: http.StatusBadRequest}}},
		}
		s := newStartupService(mock, testServerIdentifier)
		_, err := s.UpdateVariable(context.Background(), options)
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}
