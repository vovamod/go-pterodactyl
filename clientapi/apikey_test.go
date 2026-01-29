package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/errors"
	"github.com/vovamod/go-pterodactyl/internal/testutil"
)

func normaliseKeyTimes(k *api.APIKey) {
	if !k.CreatedAt.IsZero() {
		k.CreatedAt = k.CreatedAt.UTC().Truncate(time.Second)
	}
	if k.LastUsedAt != nil && !k.LastUsedAt.IsZero() {
		t := k.LastUsedAt.UTC().Truncate(time.Second)
		k.LastUsedAt = &t
	}
}

func normaliseKeySlice(list []*api.APIKey) {
	for _, k := range list {
		normaliseKeyTimes(k)
	}
}

func TestAPIKeysService_List(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	expectedKeys := []*api.APIKey{
		{
			Identifier:  "key1",
			Description: "Test Key 1",
			AllowedIPs:  []string{},
			LastUsedAt:  nil,
			CreatedAt:   now,
		},
		{
			Identifier:  "key2",
			Description: "Test Key 2",
			AllowedIPs:  []string{"127.0.0.1"},
			LastUsedAt:  &now,
			CreatedAt:   now,
		},
	}

	data := make([]*api.ListItem[api.APIKey], len(expectedKeys))
	for i, k := range expectedKeys {
		data[i] = &api.ListItem[api.APIKey]{Object: "api_key", Attributes: k}
	}
	meta := api.Meta{Pagination: api.Pagination{Total: 2, PerPage: 25, CurrentPage: 1, TotalPages: 1}}
	res := api.PaginatedResponse[api.APIKey]{Object: "list", Data: data, Meta: meta}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newAPIKeysService(mock)

		keys, m, err := s.List(context.Background(), api.PaginationOptions{Page: 1})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		normaliseKeySlice(expectedKeys)
		normaliseKeySlice(keys)

		if !reflect.DeepEqual(keys, expectedKeys) {
			t.Errorf("expected keys %+v, got %+v", expectedKeys, keys)
		}
		if !reflect.DeepEqual(m, &meta) {
			t.Errorf("expected meta %+v, got %+v", &meta, m)
		}

		req := mock.Requests[0]
		if req.Method != http.MethodGet {
			t.Errorf("expected method %s, got %s", http.MethodGet, req.Method)
		}
		if req.Endpoint != "/api/client/account/api-keys" {
			t.Errorf("expected endpoint %s, got %s", "/api/client/account/api-keys", req.Endpoint)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusInternalServerError,
				Err:        &errors.APIError{HTTPStatusCode: http.StatusInternalServerError},
			}},
		}
		s := newAPIKeysService(mock)
		if _, _, err := s.List(context.Background(), api.PaginationOptions{}); err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestAPIKeysService_Create(t *testing.T) {
	options := api.APIKeyCreateOptions{
		Description: "New Key",
		AllowedIPs:  []string{"192.168.1.1"},
	}
	jsonOptions, _ := json.Marshal(options)

	token := "secret-token-shh"
	now := time.Now().UTC().Truncate(time.Second)

	expectedKey := &api.APIKey{
		Identifier:  "new_key",
		Description: options.Description,
		AllowedIPs:  options.AllowedIPs,
		CreatedAt:   now,
		Token:       &token,
	}

	res := api.APIKeyCreateResponse{
		Object: "api_key",
		Attributes: api.APIKey{
			Identifier:  expectedKey.Identifier,
			Description: expectedKey.Description,
			AllowedIPs:  expectedKey.AllowedIPs,
			CreatedAt:   expectedKey.CreatedAt,
		},
		Meta: struct {
			SecretToken string `json:"secret_token"`
		}{SecretToken: token},
	}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusCreated, Body: jsonBody}},
		}
		s := newAPIKeysService(mock)

		key, err := s.Create(context.Background(), options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		normaliseKeyTimes(expectedKey)
		normaliseKeyTimes(key)

		if !reflect.DeepEqual(key, expectedKey) {
			t.Errorf("expected key %+v, got %+v", expectedKey, key)
		}

		req := mock.Requests[0]
		if req.Method != http.MethodPost {
			t.Errorf("expected method %s, got %s", http.MethodPost, req.Method)
		}
		if !bytes.Equal(req.Body, jsonOptions) {
			t.Errorf("expected body %s, got %s", string(jsonOptions), string(req.Body))
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusBadRequest,
				Err:        &errors.APIError{HTTPStatusCode: http.StatusBadRequest},
			}},
		}
		s := newAPIKeysService(mock)
		if _, err := s.Create(context.Background(), options); err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestAPIKeysService_Delete(t *testing.T) {
	identifier := "key_to_delete"

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusNoContent,
			}},
		}
		s := newAPIKeysService(mock)
		err := s.Delete(context.Background(), identifier)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		req := mock.Requests[0]
		if req.Method != http.MethodDelete {
			t.Errorf("expected method %s, got %s", http.MethodDelete, req.Method)
		}
		expectedEndpoint := "/api/client/account/api-keys/" + identifier
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusNotFound,
				Err:        &errors.APIError{HTTPStatusCode: http.StatusNotFound},
			}},
		}
		s := newAPIKeysService(mock)
		err := s.Delete(context.Background(), identifier)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}
