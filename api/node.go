package api

import "time"

type Node struct {
	ID                 int       `json:"id"`
	UUID               string    `json:"uuid"`
	Public             bool      `json:"public"`
	Name               string    `json:"name"`
	Description        *string   `json:"description"`
	LocationID         int       `json:"location_id"`
	FQDN               string    `json:"fqdn"`
	Scheme             string    `json:"scheme"`
	BehindProxy        bool      `json:"behind_proxy"`
	Memory             int       `json:"memory"`
	MemoryOverallocate int       `json:"memory_overallocate"`
	Disk               int       `json:"disk"`
	DiskOverallocate   int       `json:"disk_overallocate"`
	DaemonListen       int       `json:"daemon_listen"`
	DaemonSFTP         int       `json:"daemon_sftp"`
	DaemonBase         string    `json:"daemon_base"`
	MaintenanceMode    bool      `json:"maintenance_mode"`
	UploadSize         int       `json:"upload_size"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type NodeConfiguration struct {
	Debug     bool             `json:"debug"`
	UUID      string           `json:"uuid"`
	TokenID   string           `json:"token_id"`
	Token     string           `json:"token"`
	API       NodeConfigAPI    `json:"api"`
	System    NodeConfigSystem `json:"system"`
	RemoteURL string           `json:"remote"`
}

type NodeConfigAPI struct {
	Host        string        `json:"host"`
	Port        int           `json:"port"`
	SSL         NodeConfigSSL `json:"ssl"`
	UploadLimit int           `json:"upload_limit"`
}

type NodeConfigSSL struct {
	Enabled  bool   `json:"enabled"`
	CertPath string `json:"cert"`
	KeyPath  string `json:"key"`
}
type NodeConfigSystem struct {
	DataPath string         `json:"data"`
	SFTP     NodeConfigSFTP `json:"sftp"`
}
type NodeConfigSFTP struct {
	BindPort int `json:"bind_port"`
}

type NodeCreateOptions struct {
	Name               string `json:"name"`
	LocationID         int    `json:"location_id"`
	FQDN               string `json:"fqdn"`
	Scheme             string `json:"scheme"` // "http" or "https"
	Memory             int    `json:"memory"`
	MemoryOverallocate int    `json:"memory_overallocate"`
	Disk               int    `json:"disk"`
	DiskOverallocate   int    `json:"disk_overallocate"`
	DaemonSFTP         int    `json:"daemon_sftp"`
	DaemonListen       int    `json:"daemon_listen"`
	// Optional fields
	Description     *string `json:"description,omitempty"`
	BehindProxy     *bool   `json:"behind_proxy,omitempty"`
	MaintenanceMode *bool   `json:"maintenance_mode,omitempty"`
	UploadSize      *int    `json:"upload_size,omitempty"`
}

type NodeDeployableOptions struct {
	PaginationOptions
	Memory int // MB
	Disk   int // MB
}

type NodeUpdateOptions struct {
	Name               string  `json:"name,omitempty"`
	LocationID         int     `json:"location_id,omitempty"`
	FQDN               string  `json:"fqdn,omitempty"`
	Scheme             string  `json:"scheme,omitempty"`
	Memory             int     `json:"memory,omitempty"`
	MemoryOverallocate int     `json:"memory_overallocate,omitempty"`
	Disk               int     `json:"disk,omitempty"`
	DiskOverallocate   int     `json:"disk_overallocate,omitempty"`
	DaemonSFTP         int     `json:"daemon_sftp,omitempty"`
	DaemonListen       int     `json:"daemon_listen,omitempty"`
	Description        *string `json:"description,omitempty"`
	BehindProxy        *bool   `json:"behind_proxy,omitempty"`
	MaintenanceMode    *bool   `json:"maintenance_mode,omitempty"`
	UploadSize         *int    `json:"upload_size,omitempty"`
}
