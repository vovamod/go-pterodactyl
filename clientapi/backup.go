package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/internal/crud"
	"github.com/vovamod/go-pterodactyl/internal/requester"
)

type backupsService struct {
	client           requester.Requester
	serverIdentifier string
}

// newBackupsService creates a new backups service.
func newBackupsService(client requester.Requester, serverIdentifier string) *backupsService {
	return &backupsService{client: client, serverIdentifier: serverIdentifier}
}

// List retrieves all backups for the server.
func (s *backupsService) List(ctx context.Context, options api.PaginationOptions) ([]*api.Backup, *api.Meta, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/backups", s.serverIdentifier)
	return crud.List[api.Backup](ctx, s.client, endpoint, &options)
}

// Create sends a request to begin a new backup creation process.
func (s *backupsService) Create(ctx context.Context, options api.BackupCreateOptions) (*api.Backup, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create backup options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/client/servers/%s/backups", s.serverIdentifier)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new backup request: %w", err)
	}

	res := &api.BackupResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return res.Attributes, nil
}

// Details retrieves the details of a specific backup by its UUID.
func (s *backupsService) Details(ctx context.Context, uuid string) (*api.Backup, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/backups/%s", s.serverIdentifier, uuid)
	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup details request: %w", err)
	}

	res := &api.BackupResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return res.Attributes, nil
}

func (s *backupsService) Download(ctx context.Context, uuid string) (*api.BackupDownload, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/backups/%s/download", s.serverIdentifier, uuid)
	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup download request: %w", err)
	}

	res := &api.BackupDownloadResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	return res.Attributes, nil
}

func (s *backupsService) Delete(ctx context.Context, uuid string) error {
	endpoint := fmt.Sprintf("/api/client/servers/%s/backups/%s", s.serverIdentifier, uuid)
	req, err := s.client.NewRequest(ctx, "DELETE", endpoint, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete backup request: %w", err)
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}
