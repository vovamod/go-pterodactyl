package appapi

import (
	"context"
	"fmt"

	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/internal/crud"
	"github.com/vovamod/go-pterodactyl/internal/requester"
)

type eggsService struct {
	client requester.Requester
	nestID int
}

func NewEggsService(client requester.Requester, nestID int) *eggsService {
	return &eggsService{client: client, nestID: nestID}
}

func (s *eggsService) List(ctx context.Context, options *api.PaginationOptions) ([]*api.Egg, *api.Meta, error) {
	endpoint := fmt.Sprintf("/api/application/nests/%d/eggs", s.nestID)
	return crud.List[api.Egg](ctx, s.client, endpoint, options)
}

func (s *eggsService) ListAll(ctx context.Context) ([]*api.Egg, error) {
	endpoint := fmt.Sprintf("/api/application/nests/%d/eggs", s.nestID)
	return crud.ListAll[api.Egg](ctx, s.client, endpoint, 100)
}

func (s *eggsService) Get(ctx context.Context, eggID int) (*api.Egg, error) {
	endpoint := fmt.Sprintf("/api/application/nests/%d/eggs", s.nestID)
	return crud.Get[api.Egg](ctx, s.client, endpoint, eggID)
}

func (s *eggsService) GetWithVariables(ctx context.Context, eggID int) (*api.Egg, []*api.EggVariable, error) {
	endpoint := fmt.Sprintf("/api/application/nests/%d/eggs/%d", s.nestID, eggID)

	opts := &api.PaginationOptions{Include: []string{"variables"}}
	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil, opts)
	if err != nil {
		return nil, nil, err
	}

	resp := &api.EggWithVariablesResponse{}
	if _, err = s.client.Do(ctx, req, resp); err != nil {
		return nil, nil, err
	}

	vars := make([]*api.EggVariable, len(resp.Relationships.Variables.Data))
	for i, item := range resp.Relationships.Variables.Data {
		vars[i] = item.Attributes
	}

	return resp.Attributes, vars, nil
}
