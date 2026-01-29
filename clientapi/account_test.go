package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/errors"
	"github.com/vovamod/go-pterodactyl/internal/testutil"
)

func TestAccountService_GetDetails(t *testing.T) {
	expectedAccount := &api.Account{
		ID:        1,
		Admin:     true,
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Language:  "en",
	}
	res := api.ListItem[api.Account]{
		Object:     "user",
		Attributes: expectedAccount,
	}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusOK,
				Body:       jsonBody,
			}},
		}
		s := newAccountService(mock)
		account, err := s.GetDetails(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(account, expectedAccount) {
			t.Errorf("expected account %+v, got %+v", expectedAccount, account)
		}
		if len(mock.Requests) != 1 {
			t.Fatalf("expected 1 request, got %d", len(mock.Requests))
		}
		req := mock.Requests[0]
		if req.Method != http.MethodGet {
			t.Errorf("expected method %s, got %s", http.MethodGet, req.Method)
		}
		if req.Endpoint != "/api/client/account" {
			t.Errorf("expected endpoint %s, got %s", "/api/client/account", req.Endpoint)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusNotFound,
				Err:        &errors.APIError{HTTPStatusCode: http.StatusNotFound},
			}},
		}
		s := newAccountService(mock)
		_, err := s.GetDetails(context.Background())
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestAccountService_GetTwoFactorDetails(t *testing.T) {
	expectedDetails := &api.TwoFactorDetails{
		ImageURL: "test_image_url",
		Secret:   "test_secret",
	}
	jsonBody, _ := json.Marshal(expectedDetails)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusOK,
				Body:       jsonBody,
			}},
		}
		s := newAccountService(mock)
		details, err := s.GetTwoFactorDetails(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(details, expectedDetails) {
			t.Errorf("expected details %+v, got %+v", expectedDetails, details)
		}
		req := mock.Requests[0]
		if req.Method != http.MethodGet {
			t.Errorf("expected method %s, got %s", http.MethodGet, req.Method)
		}
		if req.Endpoint != "/api/client/account/two-factor" {
			t.Errorf("expected endpoint %s, got %s", "/api/client/account/two-factor", req.Endpoint)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusInternalServerError,
				Err:        &errors.APIError{HTTPStatusCode: http.StatusInternalServerError},
			}},
		}
		s := newAccountService(mock)
		_, err := s.GetTwoFactorDetails(context.Background())
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestAccountService_EnableTwoFactor(t *testing.T) {
	options := api.TwoFactorEnableOptions{Code: "123456"}
	jsonBytes, _ := json.Marshal(options)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusNoContent,
			}},
		}
		s := newAccountService(mock)
		err := s.EnableTwoFactor(context.Background(), options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		req := mock.Requests[0]
		if req.Method != http.MethodPost {
			t.Errorf("expected method %s, got %s", http.MethodPost, req.Method)
		}
		if req.Endpoint != "/api/client/account/two-factor" {
			t.Errorf("expected endpoint %s, got %s", "/api/client/account/two-factor", req.Endpoint)
		}
		if !bytes.Equal(req.Body, jsonBytes) {
			t.Errorf("expected body %s, got %s", jsonBytes, req.Body)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusBadRequest,
				Err:        &errors.APIError{HTTPStatusCode: http.StatusBadRequest},
			}},
		}
		s := newAccountService(mock)
		err := s.EnableTwoFactor(context.Background(), options)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestAccountService_DisableTwoFactor(t *testing.T) {
	options := api.TwoFactorDisableOptions{Password: "password"}
	jsonBytes, _ := json.Marshal(options)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusNoContent,
			}},
		}
		s := newAccountService(mock)
		err := s.DisableTwoFactor(context.Background(), options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		req := mock.Requests[0]
		if req.Method != http.MethodDelete {
			t.Errorf("expected method %s, got %s", http.MethodDelete, req.Method)
		}
		if req.Endpoint != "/api/client/account/two-factor" {
			t.Errorf("expected endpoint %s, got %s", "/api/client/account/two-factor", req.Endpoint)
		}
		if !bytes.Equal(req.Body, jsonBytes) {
			t.Errorf("expected body %s, got %s", jsonBytes, req.Body)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusForbidden,
				Err:        &errors.APIError{HTTPStatusCode: http.StatusForbidden},
			}},
		}
		s := newAccountService(mock)
		err := s.DisableTwoFactor(context.Background(), options)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestAccountService_UpdateEmail(t *testing.T) {
	options := api.UpdateEmailOptions{Email: "new@example.com", Password: "password"}
	jsonBytes, _ := json.Marshal(options)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusNoContent,
			}},
		}
		s := newAccountService(mock)
		err := s.UpdateEmail(context.Background(), options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		req := mock.Requests[0]
		if req.Method != http.MethodPut {
			t.Errorf("expected method %s, got %s", http.MethodPut, req.Method)
		}
		if req.Endpoint != "/api/client/account/email" {
			t.Errorf("expected endpoint %s, got %s", "/api/client/account/email", req.Endpoint)
		}
		if !bytes.Equal(req.Body, jsonBytes) {
			t.Errorf("expected body %s, got %s", jsonBytes, req.Body)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusUnprocessableEntity,
				Err:        &errors.APIError{HTTPStatusCode: http.StatusUnprocessableEntity},
			}},
		}
		s := newAccountService(mock)
		err := s.UpdateEmail(context.Background(), options)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestAccountService_UpdatePassword(t *testing.T) {
	options := api.UpdatePasswordOptions{CurrentPassword: "current", NewPassword: "new", PasswordConfirm: "new"}
	jsonBytes, _ := json.Marshal(options)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusNoContent,
			}},
		}
		s := newAccountService(mock)
		err := s.UpdatePassword(context.Background(), options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		req := mock.Requests[0]
		if req.Method != http.MethodPut {
			t.Errorf("expected method %s, got %s", http.MethodPut, req.Method)
		}
		if req.Endpoint != "/api/client/account/password" {
			t.Errorf("expected endpoint %s, got %s", "/api/client/account/password", req.Endpoint)
		}
		if !bytes.Equal(req.Body, jsonBytes) {
			t.Errorf("expected body %s, got %s", jsonBytes, req.Body)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusUnprocessableEntity,
				Err:        &errors.APIError{HTTPStatusCode: http.StatusUnprocessableEntity},
			}},
		}
		s := newAccountService(mock)
		err := s.UpdatePassword(context.Background(), options)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestAccountService_APIKeys(t *testing.T) {
	mock := &testutil.MockRequester{}
	s := newAccountService(mock)
	apiKeysService := s.APIKeys()
	if apiKeysService == nil {
		t.Fatal("expected APIKeysService to not be nil")
	}
}
