package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/errors"
	"github.com/vovamod/go-pterodactyl/internal/testutil"
)

func TestUsersService_List(t *testing.T) {
	expectedUsers := []*api.Subuser{
		{UUID: "uuid-1", Username: "user1"},
		{UUID: "uuid-2", Username: "user2"},
	}
	data := make([]*api.ListItem[api.Subuser], len(expectedUsers))
	for i, u := range expectedUsers {
		data[i] = &api.ListItem[api.Subuser]{Object: "subuser", Attributes: u}
	}
	meta := api.Meta{Pagination: api.Pagination{Total: 2, PerPage: 25}}
	res := api.PaginatedResponse[api.Subuser]{Object: "list", Data: data, Meta: meta}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newUsersService(mock, testServerIdentifier)
		users, m, err := s.List(context.Background(), api.PaginationOptions{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(users, expectedUsers) {
			t.Errorf("expected users %+v, got %+v", expectedUsers, users)
		}
		if !reflect.DeepEqual(m, &meta) {
			t.Errorf("expected meta %+v, got %+v", &meta, m)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{Err: &errors.APIError{HTTPStatusCode: http.StatusInternalServerError}}},
		}
		s := newUsersService(mock, testServerIdentifier)
		_, _, err := s.List(context.Background(), api.PaginationOptions{})
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}

func TestUsersService_Create(t *testing.T) {
	options := api.SubuserCreateOptions{Email: "new@user.com", Permissions: []string{"control.console"}}
	jsonOptions, _ := json.Marshal(options)
	expectedUser := &api.Subuser{UUID: "new-uuid", Email: options.Email, Permissions: options.Permissions}
	res := api.SubuserResponse{Object: "subuser", Attributes: expectedUser}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newUsersService(mock, testServerIdentifier)
		user, err := s.Create(context.Background(), options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(user, expectedUser) {
			t.Errorf("expected user %+v, got %+v", expectedUser, user)
		}
		req := mock.Requests[0]
		if req.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", req.Method)
		}
		if !bytes.Equal(req.Body, jsonOptions) {
			t.Errorf("expected body %s, got %s", jsonOptions, req.Body)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{Err: &errors.APIError{HTTPStatusCode: http.StatusBadRequest}}},
		}
		s := newUsersService(mock, testServerIdentifier)
		_, err := s.Create(context.Background(), options)
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}

func TestUsersService_Details(t *testing.T) {
	uuid := "test-uuid"
	expectedUser := &api.Subuser{UUID: uuid, CreatedAt: time.Now()}
	res := api.SubuserResponse{Object: "subuser", Attributes: expectedUser}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newUsersService(mock, testServerIdentifier)
		user, err := s.Details(context.Background(), uuid)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Truncate time for comparison
		expectedUser.CreatedAt = expectedUser.CreatedAt.UTC().Truncate(time.Second)
		user.CreatedAt = user.CreatedAt.UTC().Truncate(time.Second)

		if !reflect.DeepEqual(user, expectedUser) {
			t.Errorf("expected user %+v, got %+v", expectedUser, user)
		}
		req := mock.Requests[0]
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/users/%s", testServerIdentifier, uuid)
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{Err: &errors.APIError{HTTPStatusCode: http.StatusNotFound}}},
		}
		s := newUsersService(mock, testServerIdentifier)
		_, err := s.Details(context.Background(), uuid)
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}

func TestUsersService_Update(t *testing.T) {
	uuid := "test-uuid"
	options := api.SubuserUpdateOptions{Permissions: []string{"control.start", "control.stop"}}
	jsonOptions, _ := json.Marshal(options)
	expectedUser := &api.Subuser{UUID: uuid, Permissions: options.Permissions}
	res := api.SubuserResponse{Object: "subuser", Attributes: expectedUser}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newUsersService(mock, testServerIdentifier)
		user, err := s.Update(context.Background(), uuid, options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(user, expectedUser) {
			t.Errorf("expected user %+v, got %+v", expectedUser, user)
		}
		req := mock.Requests[0]
		if req.Method != http.MethodPost { // Note: Update uses POST
			t.Errorf("expected method POST, got %s", req.Method)
		}
		if !bytes.Equal(req.Body, jsonOptions) {
			t.Errorf("expected body %s, got %s", jsonOptions, req.Body)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{Err: &errors.APIError{HTTPStatusCode: http.StatusUnprocessableEntity}}},
		}
		s := newUsersService(mock, testServerIdentifier)
		_, err := s.Update(context.Background(), uuid, options)
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}

func TestUsersService_Delete(t *testing.T) {
	uuid := "test-uuid"

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusNoContent}},
		}
		s := newUsersService(mock, testServerIdentifier)
		err := s.Delete(context.Background(), uuid)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		req := mock.Requests[0]
		if req.Method != http.MethodDelete {
			t.Errorf("expected method DELETE, got %s", req.Method)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{Err: &errors.APIError{HTTPStatusCode: http.StatusNotFound}}},
		}
		s := newUsersService(mock, testServerIdentifier)
		err := s.Delete(context.Background(), uuid)
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}
