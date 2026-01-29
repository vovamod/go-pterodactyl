package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/errors"
	"github.com/vovamod/go-pterodactyl/internal/testutil"
)

func TestFilesService_List(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	expectedFiles := []*api.FileObject{
		{
			Name:       "test.txt",
			Mode:       "rw-r--r--",
			ModeBits:   "0644",
			Size:       1024,
			IsFile:     true,
			IsSymlink:  false,
			MimeType:   "text/plain",
			CreatedAt:  now,
			ModifiedAt: now,
		},
		{
			Name:      "a-folder",
			IsFile:    false,
			CreatedAt: now,
		},
	}

	// API returns a list of FileObjectResponses
	data := make([]*api.FileObjectResponse, len(expectedFiles))
	for i, f := range expectedFiles {
		data[i] = &api.FileObjectResponse{Object: "file_object", Attributes: f}
	}
	res := api.FileListResponse{Object: "list", Data: data}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newFilesService(mock, testServerIdentifier)

		files, err := s.List(context.Background(), "/data")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// ---- Normalise CreatedAt / ModifiedAt to a common Location ----
		normaliseTimes := func(objs []*api.FileObject) {
			for _, o := range objs {
				if !o.CreatedAt.IsZero() {
					o.CreatedAt = o.CreatedAt.UTC().Truncate(time.Second)
				}
				if !o.ModifiedAt.IsZero() {
					o.ModifiedAt = o.ModifiedAt.UTC().Truncate(time.Second)
				}
			}
		}
		normaliseTimes(expectedFiles)
		normaliseTimes(files)
		// ---------------------------------------------------------------

		if !reflect.DeepEqual(files, expectedFiles) {
			t.Errorf("expected files %+v, got %+v", expectedFiles, files)
		}

		expectedEndpoint := fmt.Sprintf(
			"/api/client/servers/%s/files/list?directory=%s",
			testServerIdentifier, url.QueryEscape("/data"),
		)
		if mock.Requests[0].Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, mock.Requests[0].Endpoint)
		}
	})
}

func TestFilesService_GetContents(t *testing.T) {
	filePath := "config/app.json"
	expectedContents := `{"version": "1.0"}`

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: []byte(expectedContents)}},
		}
		s := newFilesService(mock, testServerIdentifier)
		contents, err := s.GetContents(context.Background(), filePath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if contents != expectedContents {
			t.Errorf("expected contents '%s', got '%s'", expectedContents, contents)
		}
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/files/contents?file=%s", testServerIdentifier, url.QueryEscape(filePath))
		if mock.Requests[0].Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, mock.Requests[0].Endpoint)
		}
	})
}

func TestFilesService_Download(t *testing.T) {
	filePath := "archive.zip"
	expectedURL := &api.SignedURL{URL: "https://example.com/download/123"}
	res := api.SignedURLResponse{Object: "signed_url", Attributes: expectedURL}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newFilesService(mock, testServerIdentifier)
		url, err := s.Download(context.Background(), filePath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(url, expectedURL) {
			t.Errorf("expected url %+v, got %+v", expectedURL, url)
		}
	})
}

func TestFilesService_Rename(t *testing.T) {
	options := api.RenameFilesOptions{
		Root:  "/",
		Files: []api.RenameFile{{From: "old.txt", To: "new.txt"}},
	}
	jsonBody, _ := json.Marshal(options)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusNoContent}},
		}
		s := newFilesService(mock, testServerIdentifier)
		err := s.Rename(context.Background(), options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !bytes.Equal(mock.Requests[0].Body, jsonBody) {
			t.Errorf("expected body %s, got %s", jsonBody, mock.Requests[0].Body)
		}
	})
}

func TestFilesService_Copy(t *testing.T) {
	options := api.CopyFileOptions{Location: "/new/folder"}
	jsonBody, _ := json.Marshal(options)
	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusNoContent}},
		}
		s := newFilesService(mock, testServerIdentifier)
		err := s.Copy(context.Background(), options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !bytes.Equal(mock.Requests[0].Body, jsonBody) {
			t.Errorf("expected body %s, got %s", jsonBody, mock.Requests[0].Body)
		}
	})
}

func TestFilesService_Write(t *testing.T) {
	filePath := "test.txt"
	content := "hello world!"
	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusNoContent}},
		}
		s := newFilesService(mock, testServerIdentifier)
		err := s.Write(context.Background(), filePath, strings.NewReader(content))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(mock.Requests[0].Body) != content {
			t.Errorf("expected body '%s', got '%s'", content, mock.Requests[0].Body)
		}
	})
}

func TestFilesService_Compress(t *testing.T) {
	options := api.CompressFilesOptions{Root: "/", Files: []string{"test.txt"}}
	jsonBody, _ := json.Marshal(options)
	expectedFile := &api.FileObject{Name: "archive.zip", IsFile: true}
	res := api.FileObjectResponse{Object: "file_object", Attributes: expectedFile}
	resBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: resBody}},
		}
		s := newFilesService(mock, testServerIdentifier)
		file, err := s.Compress(context.Background(), options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(file, expectedFile) {
			t.Errorf("expected file %+v, got %+v", expectedFile, file)
		}
		if !bytes.Equal(mock.Requests[0].Body, jsonBody) {
			t.Errorf("expected body %s, got %s", jsonBody, mock.Requests[0].Body)
		}
	})
}

func TestFilesService_Decompress(t *testing.T) {
	options := api.DecompressFileOptions{Root: "/", File: "archive.zip"}
	jsonBody, _ := json.Marshal(options)
	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusNoContent}},
		}
		s := newFilesService(mock, testServerIdentifier)
		err := s.Decompress(context.Background(), options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !bytes.Equal(mock.Requests[0].Body, jsonBody) {
			t.Errorf("expected body %s, got %s", jsonBody, mock.Requests[0].Body)
		}
	})
}

func TestFilesService_Delete(t *testing.T) {
	options := api.DeleteFilesOptions{Root: "/", Files: []string{"test.txt"}}
	jsonBody, _ := json.Marshal(options)
	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusNoContent}},
		}
		s := newFilesService(mock, testServerIdentifier)
		err := s.Delete(context.Background(), options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !bytes.Equal(mock.Requests[0].Body, jsonBody) {
			t.Errorf("expected body %s, got %s", jsonBody, mock.Requests[0].Body)
		}
	})
}

func TestFilesService_CreateFolder(t *testing.T) {
	options := api.CreateFolderOptions{Root: "/", Name: "new-folder"}
	jsonBody, _ := json.Marshal(options)
	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusNoContent}},
		}
		s := newFilesService(mock, testServerIdentifier)
		err := s.CreateFolder(context.Background(), options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !bytes.Equal(mock.Requests[0].Body, jsonBody) {
			t.Errorf("expected body %s, got %s", jsonBody, mock.Requests[0].Body)
		}
	})
}

func TestFilesService_GetUploadURL(t *testing.T) {
	expectedURL := &api.SignedURL{URL: "https://example.com/upload/456"}
	res := api.SignedURLResponse{Object: "signed_url", Attributes: expectedURL}
	jsonBody, _ := json.Marshal(res)
	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newFilesService(mock, testServerIdentifier)
		url, err := s.GetUploadURL(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(url, expectedURL) {
			t.Errorf("expected url %+v, got %+v", url, expectedURL)
		}
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/files/upload", testServerIdentifier)
		if mock.Requests[0].Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, mock.Requests[0].Endpoint)
		}
	})
}

func TestFilesService_Error(t *testing.T) {
	apiErr := &errors.APIError{HTTPStatusCode: http.StatusNotFound}
	mock := &testutil.MockRequester{
		Responses: []testutil.MockResponse{{Err: apiErr}},
	}
	s := newFilesService(mock, testServerIdentifier)

	_, err := s.List(context.Background(), "/")
	if err != apiErr {
		t.Errorf("expected error %v, got %v", apiErr, err)
	}
}
