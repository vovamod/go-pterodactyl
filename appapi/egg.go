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
