package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/errors"
	"github.com/vovamod/go-pterodactyl/internal/testutil"
)

const testServerIdentifier = "test-server"

func normaliseBackupTimes(b *api.Backup) {
	if !b.CreatedAt.IsZero() {
		b.CreatedAt = b.CreatedAt.UTC().Truncate(time.Second)
	}
	if b.CompletedAt != nil && !b.CompletedAt.IsZero() {
		t := b.CompletedAt.UTC().Truncate(time.Second)
		b.CompletedAt = &t
	}
}

func normaliseBackupSlice(list []*api.Backup) {
	for _, b := range list {
		normaliseBackupTimes(b)
	}
}

func TestBackupsService_List(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	checksum := "test-checksum"

	expectedBackups := []*api.Backup{
		{
			UUID:         "uuid1",
			Name:         "backup1",
			IsSuccessful: true,
			IsLocked:     false,
			Bytes:        1024,
			CreatedAt:    now,
			CompletedAt:  &now,
		},
		{
			UUID:         "uuid2",
			Name:         "backup2",
			IsSuccessful: true,
			IsLocked:     true,
			IgnoredFiles: []string{"/ignore.txt"},
			Checksum:     &checksum,
			Bytes:        2048,
			CreatedAt:    now,
			CompletedAt:  &now,
		},
	}

	data := make([]*api.ListItem[api.Backup], len(expectedBackups))
	for i, b := range expectedBackups {
		data[i] = &api.ListItem[api.Backup]{Object: "backup", Attributes: b}
	}
	meta := api.Meta{Pagination: api.Pagination{Total: 2, PerPage: 25, CurrentPage: 1, TotalPages: 1}}
	res := api.PaginatedResponse[api.Backup]{Object: "list", Data: data, Meta: meta}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newBackupsService(mock, testServerIdentifier)

		backups, m, err := s.List(context.Background(), api.PaginationOptions{Page: 1})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// normalise Location pointers before DeepEqual
		normaliseBackupSlice(expectedBackups)
		normaliseBackupSlice(backups)

		if !reflect.DeepEqual(backups, expectedBackups) {
			t.Errorf("expected backups %+v, got %+v", expectedBackups, backups)
		}
		if !reflect.DeepEqual(m, &meta) {
			t.Errorf("expected meta %+v, got %+v", &meta, m)
		}

		req := mock.Requests[0]
		if req.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", req.Method)
		}
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/backups", testServerIdentifier)
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusInternalServerError,
				Err:        &errors.APIError{HTTPStatusCode: http.StatusInternalServerError},
			}},
		}
		s := newBackupsService(mock, testServerIdentifier)
		if _, _, err := s.List(context.Background(), api.PaginationOptions{}); err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestBackupsService_Create(t *testing.T) {
	name := "new-backup"
	options := api.BackupCreateOptions{Name: &name}
	jsonOptions, _ := json.Marshal(options)

	expectedBackup := &api.Backup{
		UUID:         "new-uuid",
		IsSuccessful: false, // still running
		Name:         name,
		CreatedAt:    time.Now().UTC().Truncate(time.Second),
	}
	res := api.BackupResponse{Object: "backup", Attributes: expectedBackup}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newBackupsService(mock, testServerIdentifier)

		backup, err := s.Create(context.Background(), options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		normaliseBackupTimes(expectedBackup)
		normaliseBackupTimes(backup)

		if !reflect.DeepEqual(backup, expectedBackup) {
			t.Errorf("expected backup %+v, got %+v", expectedBackup, backup)
		}

		req := mock.Requests[0]
		if req.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", req.Method)
		}
		if !bytes.Equal(req.Body, jsonOptions) {
			t.Errorf("expected body %s, got %s", jsonOptions, req.Body)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusLocked,
				Err:        &errors.APIError{HTTPStatusCode: http.StatusLocked},
			}},
		}
		s := newBackupsService(mock, testServerIdentifier)
		if _, err := s.Create(context.Background(), options); err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestBackupsService_Details(t *testing.T) {
	uuid := "test-uuid"
	expectedBackup := &api.Backup{
		UUID:         uuid,
		Name:         "details-backup",
		IsSuccessful: true,
		Bytes:        4096,
		CreatedAt:    time.Now().UTC().Truncate(time.Second),
	}
	res := api.BackupResponse{Object: "backup", Attributes: expectedBackup}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newBackupsService(mock, testServerIdentifier)

		backup, err := s.Details(context.Background(), uuid)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		normaliseBackupTimes(expectedBackup)
		normaliseBackupTimes(backup)

		if !reflect.DeepEqual(backup, expectedBackup) {
			t.Errorf("expected backup %+v, got %+v", expectedBackup, backup)
		}

		req := mock.Requests[0]
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/backups/%s", testServerIdentifier, uuid)
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusNotFound,
				Err:        &errors.APIError{HTTPStatusCode: http.StatusNotFound},
			}},
		}
		s := newBackupsService(mock, testServerIdentifier)
		if _, err := s.Details(context.Background(), uuid); err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestBackupsService_Download(t *testing.T) {
	uuid := "test-uuid"
	expectedDownload := &api.BackupDownload{URL: "https://example.com/download/backup.zip"}
	res := api.BackupDownloadResponse{
		Object:     "backup_download",
		Attributes: expectedDownload,
	}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusOK,
				Body:       jsonBody,
			}},
		}
		s := newBackupsService(mock, testServerIdentifier)
		download, err := s.Download(context.Background(), uuid)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(download, expectedDownload) {
			t.Errorf("expected download %+v, got %+v", expectedDownload, download)
		}
		req := mock.Requests[0]
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/backups/%s/download", testServerIdentifier, uuid)
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusNotFound,
				Err:        &errors.APIError{HTTPStatusCode: http.StatusNotFound},
			}},
		}
		s := newBackupsService(mock, testServerIdentifier)
		_, err := s.Download(context.Background(), uuid)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}

func TestBackupsService_Delete(t *testing.T) {
	uuid := "test-uuid-delete"

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusNoContent,
			}},
		}
		s := newBackupsService(mock, testServerIdentifier)
		err := s.Delete(context.Background(), uuid)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		req := mock.Requests[0]
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/backups/%s", testServerIdentifier, uuid)
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusNotFound,
				Err:        &errors.APIError{HTTPStatusCode: http.StatusNotFound},
			}},
		}
		s := newBackupsService(mock, testServerIdentifier)
		err := s.Delete(context.Background(), uuid)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}
