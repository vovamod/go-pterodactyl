package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/errors"
	"github.com/vovamod/go-pterodactyl/internal/testutil"
)

func TestProcessCreateOrRotateResponse(t *testing.T) {
	baseDB := &api.ClientDatabase{ID: "test-id", Name: "test_db"}
	password := "super-secret-password"

	t.Run("with password relationship", func(t *testing.T) {
		res := &api.ClientDatabaseCreateResponse{
			Object:     "server_database",
			Attributes: baseDB,
			Relationships: &struct {
				Password *struct {
					Object     string "json:\"object\""
					Attributes *struct {
						Password string "json:\"password\""
					} "json:\"attributes\""
				} "json:\"password\""
			}{
				Password: &struct {
					Object     string "json:\"object\""
					Attributes *struct {
						Password string "json:\"password\""
					} "json:\"attributes\""
				}{
					Object: "password",
					Attributes: &struct {
						Password string "json:\"password\""
					}{Password: password},
				},
			},
		}
		db := processCreateOrRotateResponse(res)
		if db.Password != password {
			t.Errorf("expected password to be '%s', got '%s'", password, db.Password)
		}
	})
}

func TestDatabasesService_List(t *testing.T) {
	expectedDBs := []*api.ClientDatabase{
		{
			ID:              "db-id-1",
			Name:            "db_one",
			Username:        "user_one",
			ConnectionsFrom: "%",
			Host: struct {
				Address string `json:"address"`
				Port    int    `json:"port"`
			}{Address: "127.0.0.1", Port: 3306},
		},
		{
			ID:              "db-id-2",
			Name:            "db_two",
			Username:        "user_two",
			ConnectionsFrom: "192.168.1.1",
			Host: struct {
				Address string `json:"address"`
				Port    int    `json:"port"`
			}{Address: "127.0.0.1", Port: 3306},
		},
	}

	data := make([]*api.ListItem[api.ClientDatabase], len(expectedDBs))
	for i, db := range expectedDBs {
		data[i] = &api.ListItem[api.ClientDatabase]{Object: "server_database", Attributes: db}
	}
	meta := api.Meta{Pagination: api.Pagination{Total: 2, PerPage: 25, CurrentPage: 1, TotalPages: 1}}
	res := api.PaginatedResponse[api.ClientDatabase]{
		Object: "list",
		Data:   data,
		Meta:   meta,
	}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusOK,
				Body:       jsonBody,
			}},
		}
		s := newDatabasesService(mock, testServerIdentifier)
		dbs, m, err := s.List(context.Background(), api.PaginationOptions{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(dbs, expectedDBs) {
			t.Errorf("expected dbs %+v, got %+v", expectedDBs, dbs)
		}
		if !reflect.DeepEqual(m, &meta) {
			t.Errorf("expected meta %+v, got %+v", &meta, m)
		}
	})
}

func TestDatabasesService_Create(t *testing.T) {
	options := api.ClientDatabaseCreateOptions{DatabaseName: "new_db", Remote: "%"}
	jsonOptions, _ := json.Marshal(options)
	password := "new-secret-password"
	// The password field from the response gets processed and added to the attributes,
	// so our expected struct should have it.
	expectedDB := &api.ClientDatabase{
		ID:       "new-db-id",
		Name:     options.DatabaseName,
		Username: "user_new",
		Password: password,
	}

	// However, the raw JSON response has the password in the relationships.
	responseDB := &api.ClientDatabase{
		ID:       "new-db-id",
		Name:     options.DatabaseName,
		Username: "user_new",
	}

	res := api.ClientDatabaseCreateResponse{
		Object:     "server_database",
		Attributes: responseDB,
		Relationships: &struct {
			Password *struct {
				Object     string "json:\"object\""
				Attributes *struct {
					Password string "json:\"password\""
				} "json:\"attributes\""
			} "json:\"password\""
		}{
			Password: &struct {
				Object     string "json:\"object\""
				Attributes *struct {
					Password string "json:\"password\""
				} "json:\"attributes\""
			}{
				Object: "password",
				Attributes: &struct {
					Password string "json:\"password\""
				}{Password: password},
			},
		},
	}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusOK,
				Body:       jsonBody,
			}},
		}
		s := newDatabasesService(mock, testServerIdentifier)
		db, err := s.Create(context.Background(), options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(db, expectedDB) {
			t.Errorf("expected db %+v, got %+v", expectedDB, db)
		}
		req := mock.Requests[0]
		if !bytes.Equal(req.Body, jsonOptions) {
			t.Errorf("expected body %s, got %s", jsonOptions, req.Body)
		}
	})
}

func TestDatabasesService_RotatePassword(t *testing.T) {
	dbID := "db-to-rotate"
	newPassword := "rotated-secret-password"
	expectedDB := &api.ClientDatabase{
		ID:       dbID,
		Name:     "rotated_db",
		Password: newPassword,
	}
	// The raw response does not contain the password in the attributes.
	responseDB := &api.ClientDatabase{
		ID:   dbID,
		Name: "rotated_db",
	}

	res := api.ClientDatabaseCreateResponse{
		Object:     "server_database",
		Attributes: responseDB,
		Relationships: &struct {
			Password *struct {
				Object     string "json:\"object\""
				Attributes *struct {
					Password string "json:\"password\""
				} "json:\"attributes\""
			} "json:\"password\""
		}{
			Password: &struct {
				Object     string "json:\"object\""
				Attributes *struct {
					Password string "json:\"password\""
				} "json:\"attributes\""
			}{
				Object: "password",
				Attributes: &struct {
					Password string "json:\"password\""
				}{Password: newPassword},
			},
		},
	}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusOK,
				Body:       jsonBody,
			}},
		}
		s := newDatabasesService(mock, testServerIdentifier)
		db, err := s.RotatePassword(context.Background(), dbID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(db, expectedDB) {
			t.Errorf("expected db %+v, got %+v", expectedDB, db)
		}
		req := mock.Requests[0]
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/databases/%s/rotate-password", testServerIdentifier, dbID)
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
	})
}

func TestDatabasesService_Delete(t *testing.T) {
	dbID := "db-to-delete"
	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{
				StatusCode: http.StatusNoContent,
			}},
		}
		s := newDatabasesService(mock, testServerIdentifier)
		err := s.Delete(context.Background(), dbID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		req := mock.Requests[0]
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/databases/%s", testServerIdentifier, dbID)
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
		s := newDatabasesService(mock, testServerIdentifier)
		err := s.Delete(context.Background(), dbID)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}
