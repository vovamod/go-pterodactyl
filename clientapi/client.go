package clientapi

import (
	"context"
	"fmt"
	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/internal/requester"
	"io"
)

type APIKeysService interface {
	List(ctx context.Context, options api.PaginationOptions) ([]*api.APIKey, *api.Meta, error)
	Create(ctx context.Context, options api.APIKeyCreateOptions) (*api.APIKey, error)
	Delete(ctx context.Context, identifier string) error
}

type DatabasesService interface {
	List(ctx context.Context, options api.PaginationOptions) ([]*api.ClientDatabase, *api.Meta, error)
	Create(ctx context.Context, options api.ClientDatabaseCreateOptions) (*api.ClientDatabase, error)
	RotatePassword(ctx context.Context, databaseID string) (*api.ClientDatabase, error)
	Delete(ctx context.Context, databaseID string) error
}

type FileService interface {
	List(ctx context.Context, directory string) ([]*api.FileObject, error)
	GetContents(ctx context.Context, filePath string) (string, error)
	Download(ctx context.Context, filePath string) (*api.SignedURL, error)
	Rename(ctx context.Context, options api.RenameFilesOptions) error
	Copy(ctx context.Context, options api.CopyFileOptions) error
	Write(ctx context.Context, filePath string, content io.Reader) error
	Compress(ctx context.Context, options api.CompressFilesOptions) (*api.FileObject, error)
	Decompress(ctx context.Context, options api.DecompressFileOptions) error
	Delete(ctx context.Context, options api.DeleteFilesOptions) error
	CreateFolder(ctx context.Context, options api.CreateFolderOptions) error
	GetUploadURL(ctx context.Context) (*api.SignedURL, error)
}

type ScheduleService interface {
	List(ctx context.Context, options api.PaginationOptions) ([]*api.Schedule, *api.Meta, error)
	Create(ctx context.Context, options api.ScheduleCreateOptions) (*api.Schedule, error)
	Details(ctx context.Context, scheduleID int) (*api.Schedule, error)
	Update(ctx context.Context, scheduleID int, options api.ScheduleUpdateOptions) (*api.Schedule, error)
	Delete(ctx context.Context, scheduleID int) error
	CreateTask(ctx context.Context, scheduleID int, options api.TaskCreateOptions) (*api.Task, error)
	UpdateTask(ctx context.Context, scheduleID, taskID int, options api.TaskUpdateOptions) (*api.Task, error)
	DeleteTask(ctx context.Context, scheduleID, taskID int) error
}

type NetworkService interface {
	ListAllocations(ctx context.Context, options api.PaginationOptions) ([]*api.Allocation, *api.Meta, error)
	AssignAllocation(ctx context.Context) (*api.Allocation, error)
	SetAllocationNote(ctx context.Context, allocationID int, options api.AllocationNoteOptions) (*api.Allocation, error)
	SetPrimaryAllocation(ctx context.Context, allocationID int) (*api.Allocation, error)
	UnassignAllocation(ctx context.Context, allocationID int) error
}

type UsersService interface {
	List(ctx context.Context, options api.PaginationOptions) ([]*api.Subuser, *api.Meta, error)
	Create(ctx context.Context, options api.SubuserCreateOptions) (*api.Subuser, error)
	Details(ctx context.Context, uuid string) (*api.Subuser, error)
	Update(ctx context.Context, uuid string, options api.SubuserUpdateOptions) (*api.Subuser, error)
	Delete(ctx context.Context, uuid string) error
}

type BackupService interface {
	List(ctx context.Context, options api.PaginationOptions) ([]*api.Backup, *api.Meta, error)
	Create(ctx context.Context, options api.BackupCreateOptions) (*api.Backup, error)
	Details(ctx context.Context, uuid string) (*api.Backup, error)
	Download(ctx context.Context, uuid string) (*api.BackupDownload, error)
	Delete(ctx context.Context, uuid string) error
}

type StartupService interface {
	ListVariables(ctx context.Context, options api.PaginationOptions) ([]*api.StartupVariable, *api.Meta, error)
	UpdateVariable(ctx context.Context, options api.UpdateVariableOptions) (*api.StartupVariable, error)
}

type SettingsService interface {
	Rename(ctx context.Context, options api.RenameOptions) error
	Reinstall(ctx context.Context) error
}

type ServersService interface {
	GetDetails(ctx context.Context) (*api.ClientServer, error)
	GetWebsocket(ctx context.Context) (*api.WebsocketDetails, error)
	GetResources(ctx context.Context) (*api.Resources, error)
	SendCommand(ctx context.Context, command string) error
	SetPowerState(ctx context.Context, signal string) error

	Databases() DatabasesService
	Files() FileService
	Schedules() ScheduleService
	Network() NetworkService
	Users() UsersService
	Backups() BackupService
	Startup() StartupService
	Settings() SettingsService
}

type AccountService interface {
	GetDetails(ctx context.Context) (*api.Account, error)
	GetTwoFactorDetails(ctx context.Context) (*api.TwoFactorDetails, error)
	EnableTwoFactor(ctx context.Context, options api.TwoFactorEnableOptions) error
	DisableTwoFactor(ctx context.Context, options api.TwoFactorDisableOptions) error
	UpdateEmail(ctx context.Context, options api.UpdateEmailOptions) error
	UpdatePassword(ctx context.Context, options api.UpdatePasswordOptions) error
	APIKeys() APIKeysService
}

type ClientAPI interface {
	ListServers(ctx context.Context, options api.PaginationOptions) ([]*api.ClientServer, *api.Meta, error)
	ListPermissions(ctx context.Context) (*api.Permission, error)

	Servers(identifier string) ServersService
	Account() AccountService
}

type ClientAPIService struct {
	client requester.Requester
}

func NewClientAPI(client requester.Requester) *ClientAPIService {
	return &ClientAPIService{client: client}
}

func (s *ClientAPIService) ListServers(ctx context.Context, options api.PaginationOptions) ([]*api.ClientServer, *api.Meta, error) {
	req, err := s.client.NewRequest(ctx, "GET", "/api/client", nil, &options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create client server list request: %w", err)
	}

	response := &api.PaginatedResponse[api.ClientServer]{}
	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*api.ClientServer, len(response.Data))
	for i, item := range response.Data {
		results[i] = item.Attributes
	}
	return results, &response.Meta, nil
}

func (s *ClientAPIService) ListPermissions(ctx context.Context) (*api.Permission, error) {
	req, err := s.client.NewRequest(ctx, "GET", "/api/client/permissions", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create permissions request: %w", err)
	}

	response := &api.Permission{}
	_, err = s.client.Do(ctx, req, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *ClientAPIService) Servers(identifier string) ServersService {
	return newServerService(c.client, identifier)
}

func (c *ClientAPIService) Account() AccountService {
	return newAccountService(c.client)
}
