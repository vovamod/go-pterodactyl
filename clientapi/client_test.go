package clientapi

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/errors"
	"github.com/vovamod/go-pterodactyl/internal/testutil"
)

func TestClientAPIService_ListServers(t *testing.T) {
	expectedServers := []*api.ClientServer{
		{Identifier: "id-1", Name: "Server 1"},
		{Identifier: "id-2", Name: "Server 2"},
	}
	data := make([]*api.ListItem[api.ClientServer], len(expectedServers))
	for i, s := range expectedServers {
		data[i] = &api.ListItem[api.ClientServer]{Object: "server", Attributes: s}
	}
	meta := api.Meta{Pagination: api.Pagination{Total: 2, PerPage: 25}}
	res := api.PaginatedResponse[api.ClientServer]{Object: "list", Data: data, Meta: meta}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		c := NewClientAPI(mock)
		servers, m, err := c.ListServers(context.Background(), api.PaginationOptions{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(servers, expectedServers) {
			t.Errorf("expected servers %+v, got %+v", expectedServers, servers)
		}
		if !reflect.DeepEqual(m, &meta) {
			t.Errorf("expected meta %+v, got %+v", &meta, m)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{Err: &errors.APIError{HTTPStatusCode: http.StatusInternalServerError}}},
		}
		c := NewClientAPI(mock)
		_, _, err := c.ListServers(context.Background(), api.PaginationOptions{})
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}

func TestClientAPIService_ListPermissions(t *testing.T) {
	expectedPermissions := &api.Permission{
		Object: "list",
		Attributes: map[string]api.PermissionDescriptor{
			"control.console": {Description: "Can send commands", Keys: map[string]string{"o": "all"}},
		},
	}
	jsonBody, _ := json.Marshal(expectedPermissions)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		c := NewClientAPI(mock)
		permissions, err := c.ListPermissions(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(permissions, expectedPermissions) {
			t.Errorf("expected permissions %+v, got %+v", expectedPermissions, permissions)
		}
	})
}

func TestClientAPIService_Servers(t *testing.T) {
	c := NewClientAPI(&testutil.MockRequester{})
	if c.Servers("some-id") == nil {
		t.Error("expected a ServersService instance, got nil")
	}
}

func TestClientAPIService_Account(t *testing.T) {
	c := NewClientAPI(&testutil.MockRequester{})
	if c.Account() == nil {
		t.Error("expected an AccountService instance, got nil")
	}
}
