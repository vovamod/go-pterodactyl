package api

import "time"

type ServerRelationships struct {
	// The API returns a paginated-like list for databases.
	// We use a pointer because this field may not always be present.
	Databases *PaginatedResponse[Database] `json:"databases,omitempty"`
}

type Server struct {
	ID            int                 `json:"id"`
	ExternalID    *string             `json:"external_id"` // Should be a pointer for nullable
	UUID          string              `json:"uuid"`
	Identifier    string              `json:"identifier"`
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	Suspended     bool                `json:"suspended"`
	Limits        ServerLimits        `json:"limits"`
	FeatureLimits ServerFeatureLimits `json:"feature_limits"`
	User          int                 `json:"user"`
	Node          int                 `json:"node"`
	Allocation    int                 `json:"allocation"`
	Nest          int                 `json:"nest"`
	Egg           int                 `json:"egg"`
	Container     ServerContainer     `json:"container"`
	UpdatedAt     *time.Time          `json:"updated_at"`
	CreatedAt     time.Time           `json:"created_at"`

	Relationships *ServerRelationships `json:"relationships,omitempty"`
}

type ServerLimits struct {
	Memory int `json:"memory"`
	Swap   int `json:"swap"`
	Disk   int `json:"disk"`
	IO     int `json:"io"`
	CPU    int `json:"cpu"`
}

// ServerFeatureLimits defines the limits on creating additional resources.
type ServerFeatureLimits struct {
	Databases   int `json:"databases"`
	Allocations int `json:"allocations"`
	Backups     int `json:"backups"`
}

// ServerContainer holds the Docker container configuration for a server.
type ServerContainer struct {
	StartupCommand string         `json:"startup_command"`
	Image          string         `json:"image"`
	Installed      BoolInt        `json:"installed"`
	Environment    map[string]any `json:"environment"`
}

type ServerCreateOptions struct {
	Name        string `json:"name"`
	User        int    `json:"user"` // Owner ID
	Nest        int    `json:"nest"`
	Egg         int    `json:"egg"`
	DockerImage string `json:"docker_image,omitempty"`
	Startup     string `json:"startup,omitempty"`

	LocationID *int `json:"location_id,omitempty"`

	NodeID     *int `json:"node_id,omitempty"`
	Allocation *struct {
		Default    int   `json:"default"`
		Additional []int `json:"additional,omitempty"`
	} `json:"allocation,omitempty"`

	Environment   *map[string]string  `json:"environment,omitempty"`
	Limits        ServerLimits        `json:"limits"`
	FeatureLimits ServerFeatureLimits `json:"feature_limits"`
	ExternalID    *string             `json:"external_id,omitempty"`
	Description   *string             `json:"description,omitempty"`
	Deploy        *struct {
		Locations   []int    `json:"locations"`
		DedicatedIP bool     `json:"dedicated_ip"`
		PortRange   []string `json:"port_range"`
	} `json:"deploy,omitempty"`

	// StartWhenCreated specifies if the server should start after being installed.
	// The json tag "start_on_completion" is correct for the Pterodactyl API.
	StartWhenCreated *bool `json:"start_on_completion,omitempty"`
}

type ServerUpdateDetailsOptions struct {
	Name        string  `json:"name,omitempty"`
	User        int     `json:"user,omitempty"` // Owner ID
	ExternalID  *string `json:"external_id,omitempty"`
	Description *string `json:"description,omitempty"`
}

type ServerUpdateBuildOptions struct {
	Allocation    int                 `json:"allocation_id"` // The primary allocation ID
	Memory        int                 `json:"memory,omitempty"`
	Swap          int                 `json:"swap,omitempty"`
	Disk          int                 `json:"disk,omitempty"`
	IO            int                 `json:"io,omitempty"`
	CPU           int                 `json:"cpu,omitempty"`
	Threads       *string             `json:"threads,omitempty"`
	FeatureLimits ServerFeatureLimits `json:"feature_limits,omitempty"`
}

type ServerUpdateStartupOptions struct {
	Startup     string             `json:"startup"`
	Environment *map[string]string `json:"environment,omitempty"`
	Egg         int                `json:"egg"`
	Image       string             `json:"image"`
	SkipScripts bool               `json:"skip_scripts"`
}

type ServerDeleteOptions struct {
	Force bool `json:"force,omitempty"`
}

type ClientServer struct {
	ServerOwner bool   `json:"server_owner"`
	Identifier  string `json:"identifier"`
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Node        string `json:"node"`
	SftpDetails struct {
		IP   string `json:"ip"`
		Port int    `json:"port"`
	} `json:"sftp_details"`
	Description   string              `json:"description"`
	Limits        ServerLimits        `json:"limits"`
	FeatureLimits ServerFeatureLimits `json:"feature_limits"`
	// In the ClientAPI API, `is_suspended` and `is_installing` might appear.
	// Using a custom boolean type is a good proactive measure.
	IsSuspended   BoolInt              `json:"is_suspended"`
	IsInstalling  BoolInt              `json:"is_installing"`
	Relationships *ServerRelationships `json:"relationships,omitempty"`
}

type SendCommandOptions struct {
	Command string `json:"command"`
}

type SetPowerStateOptions struct {
	Signal string `json:"signal"` // "start", "stop", "restart", "kill"
}
