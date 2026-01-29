package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/internal/requester"
)

type schedulesService struct {
	client           requester.Requester
	serverIdentifier string
}

// newSchedulesService creates a new schedules service.
func newSchedulesService(client requester.Requester, serverIdentifier string) *schedulesService {
	return &schedulesService{client: client, serverIdentifier: serverIdentifier}
}

// List retrieves all schedules for the server. Note: tasks are not included in this list.
func (s *schedulesService) List(ctx context.Context, options api.PaginationOptions) ([]*api.Schedule, *api.Meta, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/schedules", s.serverIdentifier)
	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil, &options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create list schedules request: %w", err)
	}

	res := &api.PaginatedResponse[api.Schedule]{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*api.Schedule, len(res.Data))
	for i, item := range res.Data {
		results[i] = item.Attributes
	}
	return results, &res.Meta, nil
}

// Create sends a request to create a new schedule.
func (s *schedulesService) Create(ctx context.Context, options api.ScheduleCreateOptions) (*api.Schedule, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create schedule options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/client/servers/%s/schedules", s.serverIdentifier)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new schedule request: %w", err)
	}

	res := &api.ScheduleResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}
	return res.Attributes, nil
}

// Details retrieves details for a single schedule, including its tasks.
func (s *schedulesService) Details(ctx context.Context, scheduleID int) (*api.Schedule, error) {
	endpoint := fmt.Sprintf("/api/client/servers/%s/schedules/%d", s.serverIdentifier, scheduleID)
	req, err := s.client.NewRequest(ctx, "GET", endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create schedule details request: %w", err)
	}

	res := &api.ScheduleDetailResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	// Manually populate the tasks from the relationships object into the main schedule object
	schedule := res.Attributes
	if res.Relationships != nil && res.Relationships.Tasks != nil {
		tasks := make([]*api.Task, len(res.Relationships.Tasks.Data))
		for i, item := range res.Relationships.Tasks.Data {
			tasks[i] = item.Attributes
		}
		schedule.Tasks = tasks
	}
	return schedule, nil
}

// Update sends a request to modify an existing schedule.
func (s *schedulesService) Update(ctx context.Context, scheduleID int, options api.ScheduleUpdateOptions) (*api.Schedule, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update schedule options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/client/servers/%s/schedules/%d", s.serverIdentifier, scheduleID)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create update schedule request: %w", err)
	}

	res := &api.ScheduleResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}
	return res.Attributes, nil
}

// Delete removes a schedule from the server.
func (s *schedulesService) Delete(ctx context.Context, scheduleID int) error {
	endpoint := fmt.Sprintf("/api/client/servers/%s/schedules/%d", s.serverIdentifier, scheduleID)
	req, err := s.client.NewRequest(ctx, "DELETE", endpoint, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete schedule request: %w", err)
	}
	_, err = s.client.Do(ctx, req, nil)
	return err
}

// CreateTask adds a new task to an existing schedule.
func (s *schedulesService) CreateTask(ctx context.Context, scheduleID int, options api.TaskCreateOptions) (*api.Task, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create task options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/client/servers/%s/schedules/%d/tasks", s.serverIdentifier, scheduleID)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new task request: %w", err)
	}

	res := &api.TaskResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}
	return res.Attributes, nil
}

// UpdateTask modifies an existing task in a schedule.
func (s *schedulesService) UpdateTask(ctx context.Context, scheduleID, taskID int, options api.TaskUpdateOptions) (*api.Task, error) {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update task options: %w", err)
	}

	endpoint := fmt.Sprintf("/api/client/servers/%s/schedules/%d/tasks/%d", s.serverIdentifier, scheduleID, taskID)
	req, err := s.client.NewRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create update task request: %w", err)
	}

	res := &api.TaskResponse{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}
	return res.Attributes, nil
}

// DeleteTask removes a task from a schedule.
func (s *schedulesService) DeleteTask(ctx context.Context, scheduleID, taskID int) error {
	endpoint := fmt.Sprintf("/api/client/servers/%s/schedules/%d/tasks/%d", s.serverIdentifier, scheduleID, taskID)
	req, err := s.client.NewRequest(ctx, "DELETE", endpoint, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete task request: %w", err)
	}
	_, err = s.client.Do(ctx, req, nil)
	return err
}
