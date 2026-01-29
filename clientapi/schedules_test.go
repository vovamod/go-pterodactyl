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

func TestSchedulesService_List(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	expectedSchedules := []*api.Schedule{
		{ID: 1, Name: "Daily Restart", IsActive: true, CreatedAt: now},
		{ID: 2, Name: "Weekly Backup", IsActive: false, CreatedAt: now},
	}
	data := make([]*api.ListItem[api.Schedule], len(expectedSchedules))
	for i, s := range expectedSchedules {
		data[i] = &api.ListItem[api.Schedule]{Object: "schedule", Attributes: s}
	}
	meta := api.Meta{Pagination: api.Pagination{Total: 2, PerPage: 25, CurrentPage: 1, TotalPages: 1}}
	res := api.PaginatedResponse[api.Schedule]{Object: "list", Data: data, Meta: meta}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newSchedulesService(mock, testServerIdentifier)

		schedules, m, err := s.List(context.Background(), api.PaginationOptions{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// ---- Normalise CreatedAt so zone pointers are identical ----
		normaliseTimes := func(list []*api.Schedule) {
			for _, sch := range list {
				sch.CreatedAt = sch.CreatedAt.UTC().Truncate(time.Second)
			}
		}
		normaliseTimes(expectedSchedules)
		normaliseTimes(schedules)

		if !reflect.DeepEqual(schedules, expectedSchedules) {
			t.Errorf("expected schedules %+v, got %+v", expectedSchedules, schedules)
		}
		if !reflect.DeepEqual(m, &meta) {
			t.Errorf("expected meta %+v, got %+v", &meta, m)
		}
	})
}

func TestSchedulesService_Create(t *testing.T) {
	options := api.ScheduleCreateOptions{Name: "New Schedule", Minute: "*/5"}
	jsonOptions, _ := json.Marshal(options)
	expectedSchedule := &api.Schedule{ID: 3, Name: options.Name}
	res := api.ScheduleResponse{Object: "schedule", Attributes: expectedSchedule}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newSchedulesService(mock, testServerIdentifier)
		schedule, err := s.Create(context.Background(), options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(schedule, expectedSchedule) {
			t.Errorf("expected schedule %+v, got %+v", expectedSchedule, schedule)
		}
		if !bytes.Equal(mock.Requests[0].Body, jsonOptions) {
			t.Errorf("expected body %s, got %s", jsonOptions, mock.Requests[0].Body)
		}
	})
}

func TestSchedulesService_Details(t *testing.T) {
	scheduleID := 1
	baseSchedule := &api.Schedule{ID: scheduleID, Name: "Test Schedule"}
	tasks := []*api.Task{
		{ID: 1, Action: "command", Payload: "say hello"},
		{ID: 2, Action: "power", Payload: "restart"},
	}
	// Build the complex response struct
	taskData := make([]*api.ListItem[api.Task], len(tasks))
	for i, task := range tasks {
		taskData[i] = &api.ListItem[api.Task]{Object: "schedule_task", Attributes: task}
	}
	res := api.ScheduleDetailResponse{
		Object:     "schedule",
		Attributes: baseSchedule,
		Relationships: &struct {
			Tasks *api.PaginatedResponse[api.Task] `json:"tasks"`
		}{
			Tasks: &api.PaginatedResponse[api.Task]{
				Object: "list",
				Data:   taskData,
			},
		},
	}
	jsonBody, _ := json.Marshal(res)

	// The final schedule should have the tasks populated
	expectedSchedule := &api.Schedule{
		ID:    baseSchedule.ID,
		Name:  baseSchedule.Name,
		Tasks: tasks,
	}

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newSchedulesService(mock, testServerIdentifier)
		schedule, err := s.Details(context.Background(), scheduleID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(schedule, expectedSchedule) {
			t.Errorf("expected schedule %+v, got %+v", expectedSchedule, schedule)
		}
	})
}

func TestSchedulesService_Update(t *testing.T) {
	scheduleID := 1
	options := api.ScheduleUpdateOptions{Name: "Updated Name"}
	jsonOptions, _ := json.Marshal(options)
	expectedSchedule := &api.Schedule{ID: scheduleID, Name: options.Name}
	res := api.ScheduleResponse{Object: "schedule", Attributes: expectedSchedule}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newSchedulesService(mock, testServerIdentifier)
		schedule, err := s.Update(context.Background(), scheduleID, options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(schedule, expectedSchedule) {
			t.Errorf("expected schedule %+v, got %+v", expectedSchedule, schedule)
		}
		req := mock.Requests[0]
		if !bytes.Equal(req.Body, jsonOptions) {
			t.Errorf("expected body %s, got %s", jsonOptions, req.Body)
		}
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/schedules/%d", testServerIdentifier, scheduleID)
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
	})
}

func TestSchedulesService_Delete(t *testing.T) {
	scheduleID := 1
	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusNoContent}},
		}
		s := newSchedulesService(mock, testServerIdentifier)
		err := s.Delete(context.Background(), scheduleID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestSchedulesService_CreateTask(t *testing.T) {
	scheduleID := 1
	options := api.TaskCreateOptions{Action: "command", Payload: "say test"}
	jsonOptions, _ := json.Marshal(options)
	expectedTask := &api.Task{ID: 3, Action: options.Action, Payload: options.Payload}
	res := api.TaskResponse{Object: "schedule_task", Attributes: expectedTask}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newSchedulesService(mock, testServerIdentifier)
		task, err := s.CreateTask(context.Background(), scheduleID, options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(task, expectedTask) {
			t.Errorf("expected task %+v, got %+v", expectedTask, task)
		}
		if !bytes.Equal(mock.Requests[0].Body, jsonOptions) {
			t.Errorf("expected body %s, got %s", jsonOptions, mock.Requests[0].Body)
		}
	})
}

func TestSchedulesService_UpdateTask(t *testing.T) {
	scheduleID, taskID := 1, 1
	options := api.TaskUpdateOptions{Payload: "say updated"}
	jsonOptions, _ := json.Marshal(options)
	expectedTask := &api.Task{ID: taskID, Payload: options.Payload}
	res := api.TaskResponse{Object: "schedule_task", Attributes: expectedTask}
	jsonBody, _ := json.Marshal(res)

	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusOK, Body: jsonBody}},
		}
		s := newSchedulesService(mock, testServerIdentifier)
		task, err := s.UpdateTask(context.Background(), scheduleID, taskID, options)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(task, expectedTask) {
			t.Errorf("expected task %+v, got %+v", expectedTask, task)
		}
		if !bytes.Equal(mock.Requests[0].Body, jsonOptions) {
			t.Errorf("expected body %s, got %s", jsonOptions, mock.Requests[0].Body)
		}
	})
}

func TestSchedulesService_DeleteTask(t *testing.T) {
	scheduleID, taskID := 1, 1
	t.Run("success", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{StatusCode: http.StatusNoContent}},
		}
		s := newSchedulesService(mock, testServerIdentifier)
		err := s.DeleteTask(context.Background(), scheduleID, taskID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		req := mock.Requests[0]
		expectedEndpoint := fmt.Sprintf("/api/client/servers/%s/schedules/%d/tasks/%d", testServerIdentifier, scheduleID, taskID)
		if req.Endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got %s", expectedEndpoint, req.Endpoint)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &testutil.MockRequester{
			Responses: []testutil.MockResponse{{Err: &errors.APIError{HTTPStatusCode: http.StatusNotFound}}},
		}
		s := newSchedulesService(mock, testServerIdentifier)
		err := s.DeleteTask(context.Background(), scheduleID, taskID)
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}
