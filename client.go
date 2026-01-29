package pterodactyl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/appapi"
	"github.com/vovamod/go-pterodactyl/clientapi"
	"github.com/vovamod/go-pterodactyl/errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Public types
// ---------------------------------------------------------------------------

// KeyType distinguishes admin (Application) and user (Client) tokens.
// Prefix verification in NewClient guards accidental mix‑ups.
type KeyType int

const (
	ApplicationKey KeyType = iota
	ClientKey
)

// Client is the root value of the SDK.  It is cheap to create but holds a
// connection‑pooled *http.Client that ***must be reused*** instead of creating a
// new one per request.
//
//   sdk, _ := pterodactyl.NewClient(baseURL, token, pterodactyl.ApplicationKey)
//   nodes, _ := sdk.ApplicationAPI.Nodes.ListAll(ctx, 0)
//
// Client is ***safe for concurrent use*** by multiple goroutines.
//
// Functional options allow callers to insert their own *http.Client or
// http.RoundTripper when tighter time‑outs, custom proxies, or tracing is
// required.
//   retryTr := retrace.NewRoundTripper(http.DefaultTransport)
//   sdk, _ := pterodactyl.NewClient(baseURL, token, key,
//       pterodactyl.WithTransport(retryTr))
//
// Any zero‑value field not set by an option receives a sensible default.
//
// Note: xAPIService fields are exported so that external code can embed or
// stub them in tests.
// ---------------------------------------------------------------------------

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client

	ApplicationAPI *appapi.ApplicationAPIService
	ClientAPI      *clientapi.ClientAPIService
}

// Option configures a Client instance lazily.
// Implemented via functional‑options pattern.
// "With…" helpers live below.
//
// Options are **applied in the order passed** to NewClient; later options may
// overwrite earlier ones.
// ---------------------------------------------------------------------------

type Option func(*Client)

// WithHTTPClient replaces the default http.Client entirely.
//
//	c, _ := pterodactyl.NewClient(baseURL, key, t, pterodactyl.WithHTTPClient(my))
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithTransport swaps only the underlying RoundTripper, keeping other settings
// (Timeout, Jar…) untouched.
func WithTransport(rt http.RoundTripper) Option {
	return func(c *Client) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{Transport: rt}
			return
		}
		c.httpClient.Transport = rt
	}
}

// WithTimeout changes the http.Client.Timeout – handy for short‑lived CLI
// programs.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{Timeout: d}
			return
		}
		c.httpClient.Timeout = d
	}
}

// NewClient validates the arguments, builds a reusable SDK instance and applies
// any functional options.
// baseURL must be scheme+host, apiKey must start with ptla_ or ptlc_.
//
// Recommended default http.Client timeout is 10s – callers can override via
// WithTimeout.
// ---------------------------------------------------------------------------

func NewClient(baseURL, apiKey string, keyType KeyType, opts ...Option) (*Client, error) {
	// ----- token sanity check -------------------------------------------------
	if keyType == ApplicationKey && !strings.HasPrefix(apiKey, "ptla_") {
		return nil, fmt.Errorf("invalid application key: must start with 'ptla_'")
	}
	if keyType == ClientKey && !strings.HasPrefix(apiKey, "ptlc_") {
		return nil, fmt.Errorf("invalid client key: must start with 'ptlc_'")
	}

	// ----- URL sanity check ---------------------------------------------------
	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, fmt.Errorf("invalid baseURL: %w", err)
	}

	// ----- construct default instance ----------------------------------------
	c := &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}

	// Apply caller‑supplied options
	for _, o := range opts {
		o(c)
	}

	// ----- wire sub‑services --------------------------------------------------
	c.ApplicationAPI = &appapi.ApplicationAPIService{}
	c.ApplicationAPI.Users = appapi.NewUsersService(c)
	c.ApplicationAPI.Nodes = appapi.NewNodesService(c)
	c.ApplicationAPI.Locations = appapi.NewLocationService(c)
	c.ApplicationAPI.Servers = appapi.NewServersService(c)
	c.ApplicationAPI.Nests = appapi.NewNestsService(c)

	c.ClientAPI = clientapi.NewClientAPI(c)

	return c, nil
}

func (c *Client) NewRequest(ctx context.Context, method, endpoint string, body io.Reader, options *api.PaginationOptions) (*http.Request, error) {

	rel, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	if options != nil {
		q := rel.Query()
		if options.Page > 0 {
			q.Set("page", strconv.Itoa(options.Page))
		}
		if options.PerPage > 0 {
			q.Set("per_page", strconv.Itoa(options.PerPage))
		}
		if len(options.Include) > 0 {
			q.Set("include", strings.Join(options.Include, ","))
		}
		rel.RawQuery = q.Encode()
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}
	fullURL := u.ResolveReference(rel)

	req, err := http.NewRequestWithContext(ctx, method, fullURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Accept", "application/json")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v any) (*http.Response, error) {
	req = req.WithContext(ctx)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer res.Body.Close() // ignore error

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		// Error handling logic
		apiErr := &errors.APIError{HTTPStatusCode: res.StatusCode}
		if err = json.NewDecoder(res.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("pterodactyl: API error (status %d), failed to parse error response: %w", res.StatusCode, err)
		}
		return nil, apiErr
	}

	// If v is not nil, decode the successful response body into it.
	if v != nil {
		if err = json.NewDecoder(res.Body).Decode(v); err != nil {
			return nil, fmt.Errorf("failed to decode successful response: %w", err)
		}
	}

	return res, nil
}

// unmarshalList is an internal helper that decodes a paginated list response
// from the Pterodactyl API and flattens it into a simple slice of models.
// It uses generics to work with any model type (api.User, api.Server, etc.).
func unmarshalList[T any](body io.Reader) ([]*T, *api.Meta, error) {
	// Create an instance of our generic response wrapper.
	// We pass the type T to it.
	response := &api.PaginatedResponse[T]{}

	// Decode the entire JSON response into our struct.
	if err := json.NewDecoder(body).Decode(response); err != nil {
		return nil, nil, fmt.Errorf("failed to decode api list response: %w", err)
	}

	// Flatten the nested structure into a simple slice of models.
	// This is the logic you wanted to avoid repeating!
	results := make([]*T, len(response.Data))
	for i, item := range response.Data {
		results[i] = item.Attributes
	}

	// Return the flattened list, the pagination metadata, and no error.
	return results, &response.Meta, nil
}
