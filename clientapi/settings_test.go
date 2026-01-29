package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/errors"
	"github.com/vovamod/go-pterodactyl/internal/testutil"
)

func TestSettingsService_Rename(t *testing.T) {
	desc := "new description"
	testCases := []struct {
		name    string
		options api.RenameOptions
	}{
		{
			name:    "with name and description",
			options: api.RenameOptions{Name: "new-name", Description: &desc},
		},
		{
			name:    "with name only",
			options: api.RenameOptions{Name: "just-a-name"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tc.options)
			mock := &testutil.MockRequester{
				Responses: []testutil.MockResponse{{StatusCode: http.StatusNoContent}},
			}
			s := newSettingsService(mock, testServerIdentifier)
			err := s.Rename(context.Background(), tc.options)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			req := mock.Requests[0]
			expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/settings/rename", testServerIdentifier)
			if req.Endpoint != expectedEndpoint {
				t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
			}
			if req.Method != http.MethodPost {
				t.Errorf("expected method POST, got %s", req.Method)
			}
			if !bytes.Equal(req.Body, jsonBody) {
				t.Errorf("expected body %s, got %s", jsonBody, req.Body)
			}
		})
	}

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{Err: &errors.APIError{HTTPStatusCode: http.StatusBadRequest}}},
		}
		s := newSettingsService(mock, testServerIdentifier)
		err := s.Rename(context.Background(), api.RenameOptions{Name: "fail"})
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}

func TestSettingsService_Reinstall(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusNoContent}},
		}
		s := newSettingsService(mock, testServerIdentifier)
		err := s.Reinstall(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		req := mock.Requests[0]
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/settings/reinstall", testServerIdentifier)
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
		if req.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", req.Method)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{Err: &errors.APIError{HTTPStatusCode: http.StatusConflict}}},
		}
		s := newSettingsService(mock, testServerIdentifier)
		err := s.Reinstall(context.Background())
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}
