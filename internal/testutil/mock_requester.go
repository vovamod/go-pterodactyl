package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/errors"
	"io"
	"net/http"
)

// mockRequester implements the requester.Requester interface for testing
type MockRequester struct {
	Requests     []MockRequest
	Responses    []MockResponse
	CurrentIndex int
}

type MockRequest struct {
	Method   string
	Endpoint string
	Body     []byte
	Options  *api.PaginationOptions
}

type MockResponse struct {
	StatusCode int
	Body       []byte
	Err        error
}

func (m *MockRequester) NewRequest(ctx context.Context, method, endpoint string, body io.Reader, options *api.PaginationOptions) (*http.Request, error) {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = io.ReadAll(body)
	}

	m.Requests = append(m.Requests, MockRequest{
		Method:   method,
		Endpoint: endpoint,
		Body:     bodyBytes,
		Options:  options,
	})

	// Create a minimal request for testing
	req, _ := http.NewRequestWithContext(ctx, method, "http://test.com"+endpoint, bytes.NewReader(bodyBytes))
	return req, nil
}

func (m *MockRequester) Do(ctx context.Context, req *http.Request, v any) (*http.Response, error) {
	if m.CurrentIndex >= len(m.Responses) {
		return nil, fmt.Errorf("no more mock Responses available")
	}

	response := m.Responses[m.CurrentIndex]
	m.CurrentIndex++

	if response.Err != nil {
		return nil, response.Err
	}

	// Create a mock response
	resp := &http.Response{
		StatusCode: response.StatusCode,
		Body:       io.NopCloser(bytes.NewReader(response.Body)),
	}

	// If we have a target to decode into and the response is successful
	if v != nil && response.StatusCode >= 200 && response.StatusCode < 300 {
		if err := json.NewDecoder(bytes.NewReader(response.Body)).Decode(v); err != nil {
			return nil, err
		}
	}

	// If it's an error response, create an APIError
	if response.StatusCode >= 400 {
		apiErr := &errors.APIError{HTTPStatusCode: response.StatusCode}
		if len(response.Body) > 0 {
			err := json.NewDecoder(bytes.NewReader(response.Body)).Decode(apiErr)
			if err != nil {
				return nil, err
			}
		}
		return nil, apiErr
	}

	return resp, nil
}
