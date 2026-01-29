package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/internal/requester"
)

type filesService struct {
	client           requester.Requester
	serverIdentifier string
}

// newFilesService creates a new files service.
func newFilesService(client requester.Requester, serverIdentifier string) *filesService {
	return &filesService{client: client, serverIdentifier: serverIdentifier}
}

func (s *filesService) List(ctx context.Context, directory string) ([]*api.FileObject, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/files/list?directory=%s", s.serverIdentifier, url.QueryEscape(directory))
	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create list files request: %w", err)
	}

	res := &api.FileListResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	results := make([]*api.FileObject, len(res.Data))
	for i, item := range res.Data {
		results[i] = item.Attributes
	}
	return results, nil
}

func (s *filesService) GetContents(ctx context.Context, filePath string) (string, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/files/contents?file=%s", s.serverIdentifier, url.QueryEscape(filePath))
	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create get contents request: %w", err)
	}

	// This endpoint returns raw text, not JSON.
	httpRes, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return "", err
	}
	defer httpRes.Body.Close()

	bodyBytes, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	return string(bodyBytes), nil
}

func (s *filesService) Download(ctx context.Context, filePath string) (*api.SignedURL, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/files/download?file=%s", s.serverIdentifier, url.QueryEscape(filePath))
	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create download file request: %w", err)
	}

	res := &api.SignedURLResponse{}
	_, err = s.client.Do(ctx, req, res)
	return res.Attributes, err
}

func (s *filesService) Rename(ctx context.Context, options api.RenameFilesOptions) error {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return fmt.Errorf("failed to marshal rename options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/client/servers/%s/files/rename", s.serverIdentifier)
	req, err := s.client.NewRequest(ctx, "PUT", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return fmt.Errorf("failed to create rename file request: %w", err)
	}
	_, err = s.client.Do(ctx, req, nil)
	return err
}

func (s *filesService) Copy(ctx context.Context, options api.CopyFileOptions) error {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return fmt.Errorf("failed to marshal copy options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/client/servers/%s/files/copy", s.serverIdentifier)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return fmt.Errorf("failed to create copy file request: %w", err)
	}
	_, err = s.client.Do(ctx, req, nil)
	return err
}

func (s *filesService) Write(ctx context.Context, filePath string, content io.Reader) error {
	endpoint := fmt.Sprintf("/api/client/servers/%s/files/write?file=%s", s.serverIdentifier, url.QueryEscape(filePath))
	req, err := s.client.NewRequest(ctx, "POST", endpoint, content, nil)
	if err != nil {
		return fmt.Errorf("failed to create write file request: %w", err)
	}
	// The body is raw data, so set the content type appropriately.
	req.Header.Set("Content-Type", "text/plain")

	_, err = s.client.Do(ctx, req, nil)
	return err
}

func (s *filesService) Compress(ctx context.Context, options api.CompressFilesOptions) (*api.FileObject, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal compress options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/client/servers/%s/files/compress", s.serverIdentifier)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create compress files request: %w", err)
	}

	res := &api.FileObjectResponse{}
	_, err = s.client.Do(ctx, req, res)
	return res.Attributes, err
}

func (s *filesService) Decompress(ctx context.Context, options api.DecompressFileOptions) error {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return fmt.Errorf("failed to marshal decompress options: %w", err)
	}
	endpoint := fmt.Sprintf("/api/client/servers/%s/files/decompress", s.serverIdentifier)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return fmt.Errorf("failed to create decompress file request: %w", err)
	}
	_, err = s.client.Do(ctx, req, nil)
	return err
}

func (s *filesService) Delete(ctx context.Context, options api.DeleteFilesOptions) error {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return fmt.Errorf("failed to marshal delete options: %w", err)
	}
	endpoint := fmt.Sprintf("/api/client/servers/%s/files/delete", s.serverIdentifier)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete files request: %w", err)
	}
	_, err = s.client.Do(ctx, req, nil)
	return err
}

func (s *filesService) CreateFolder(ctx context.Context, options api.CreateFolderOptions) error {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return fmt.Errorf("failed to marshal create folder options: %w", err)
	}
	endpoint := fmt.Sprintf("/api/client/servers/%s/files/create-folder", s.serverIdentifier)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return fmt.Errorf("failed to create folder request: %w", err)
	}
	_, err = s.client.Do(ctx, req, nil)
	return err
}

func (s *filesService) GetUploadURL(ctx context.Context) (*api.SignedURL, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/files/upload", s.serverIdentifier)
	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create upload URL request: %w", err)
	}
	res := &api.SignedURLResponse{}
	_, err = s.client.Do(ctx, req, res)
	return res.Attributes, err
}
