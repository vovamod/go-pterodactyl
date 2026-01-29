package appapi

import (
	"context"
	"github.com/vovamod/go-pterodactyl/api"
)

type AllocationsService interface {
	List(ctx context.Context, options *api.PaginationOptions) ([]*api.Allocation, *api.Meta, error)
	ListAll(ctx context.Context) ([]*api.Allocation, error)
	Create(ctx context.Context, options api.AllocationCreateOptions) error //TODO Include int / allocation on return?
	Delete(ctx context.Context, allocationID int) error
}

type DatabaseService interface {
	List(ctx context.Context, options api.PaginationOptions) ([]*api.Database, *api.Meta, error)
	Get(ctx context.Context, databaseID int) (*api.Database, error)
	Create(ctx context.Context, options api.DatabaseCreateOptions) (*api.Database, error)
	ResetPassword(ctx context.Context, databaseID int) error
	Delete(ctx context.Context, databaseID int) error
}

type NodesService interface {
	List(ctx context.Context, options *api.PaginationOptions) ([]*api.Node, *api.Meta, error)
	ListAll(ctx context.Context) ([]*api.Node, error)
	Get(ctx context.Context, id int) (*api.Node, error)
	GetConfiguration(ctx context.Context, nodeID int) (*api.NodeConfiguration, error)
	Create(ctx context.Context, options api.NodeCreateOptions) (*api.Node, error)
	Update(ctx context.Context, nodeID int, options api.NodeUpdateOptions) (*api.Node, error)
	Delete(ctx context.Context, nodeID int) error
	Allocations(ctx context.Context, nodeID int) AllocationsService
}

type EggsService interface {
	List(ctx context.Context, options *api.PaginationOptions) ([]*api.Egg, *api.Meta, error)
	ListAll(ctx context.Context) ([]*api.Egg, error)
	Get(ctx context.Context, eggID int) (*api.Egg, error)
}

// NestsService defines the actions for nests and provides access to egg management.
type NestsService interface {
	List(ctx context.Context, options *api.PaginationOptions) ([]*api.Nest, *api.Meta, error)
	ListAll(ctx context.Context) ([]*api.Nest, error)
	Get(ctx context.Context, id int) (*api.Nest, error)
	Eggs(nestID int) EggsService // Returns the EggsService interface
}

// UsersService defines the actions for users.
type UsersService interface {
	List(ctx context.Context, options *api.PaginationOptions) ([]*api.User, *api.Meta, error)
	ListAll(ctx context.Context) ([]*api.User, error)
	Get(ctx context.Context, id int) (*api.User, error)
	GetExternalID(ctx context.Context, externalId string) (*api.User, error)
	Create(ctx context.Context, options api.UserCreateOptions) (*api.User, error)
	Update(ctx context.Context, id int, options api.UserUpdateOptions) (*api.User, error)
	Delete(ctx context.Context, id int) error
}

type ServersService interface {
	List(ctx context.Context, options api.PaginationOptions) ([]*api.Server, *api.Meta, error)
	ListAll(ctx context.Context) ([]*api.Server, error)
	Get(ctx context.Context, id int) (*api.Server, error)
	GetExternal(ctx context.Context, externalID string) (*api.Server, error)
	Create(ctx context.Context, options api.ServerCreateOptions) (*api.Server, error)
	UpdateDetails(ctx context.Context, serverID int, options api.ServerUpdateDetailsOptions) (*api.Server, error)
	UpdateBuild(ctx context.Context, serverID int, options api.ServerUpdateBuildOptions) (*api.Server, error)
	UpdateStartup(ctx context.Context, serverID int, options api.ServerUpdateStartupOptions) (*api.Server, error)
	Suspend(ctx context.Context, serverID int) error
	Unsuspend(ctx context.Context, serverID int) error
	Reinstall(ctx context.Context, serverID int) error
	Delete(ctx context.Context, serverID int, force bool) error
	Databases(ctx context.Context, serverID int) DatabaseService
}

// LocationsService defines the actions for locations.
type LocationsService interface {
	List(ctx context.Context, options *api.PaginationOptions) ([]*api.Location, *api.Meta, error)
	ListAll(ctx context.Context) ([]*api.Location, error)
	Get(ctx context.Context, id int) (*api.Location, error)
	Create(ctx context.Context, options api.LocationCreateOptions) (*api.Location, error)
	Update(ctx context.Context, id int, options api.LocationUpdateOptions) (*api.Location, error)
	Delete(ctx context.Context, id int) error
}

// ApplicationAPIService is the container for all top-level API services.
type ApplicationAPIService struct {
	Users     UsersService
	Nodes     NodesService
	Locations LocationsService
	Servers   ServersService
	Nests     NestsService
}
