package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/internal/requester"
)

type apiKeysService struct{ client requester.Requester }

func newAPIKeysService(client requester.Requester) APIKeysService {
	return &apiKeysService{client: client}
}

func (s *apiKeysService) List(ctx context.Context, options api.PaginationOptions) ([]*api.APIKey, *api.Meta, error) {
	req, err := s.client.NewRequest(ctx, "GET", "/api/client/account/api-keys", nil, &options)
	if err != nil {
		return nil, nil, err
	}

	res := &api.PaginatedResponse[api.APIKey]{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*api.APIKey, len(res.Data))
	for i, item := range res.Data {
		results[i] = item.Attributes
	}
	return results, &res.Meta, nil
}

func (s *apiKeysService) Create(ctx context.Context, options api.APIKeyCreateOptions) (*api.APIKey, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "POST", "/api/client/account/api-keys", bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, err
	}

	res := &api.APIKeyCreateResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	apiKey := &res.Attributes

	apiKey.Token = &res.Meta.SecretToken

	return apiKey, nil
}

func (s *apiKeysService) Delete(ctx context.Context, identifier string) error {
	endpoint := fmt.Sprintf("/api/client/account/api-keys/%s", identifier)
	req, err := s.client.NewRequest(ctx, "DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(ctx, req, nil)
	return err
}
