package api

import "time"

type Egg struct {
	ID           int               `json:"id"`
	UUID         string            `json:"uuid"`
	Name         string            `json:"name"`
	NestID       int               `json:"nest"`
	Author       string            `json:"author"`
	Description  string            `json:"description"`
	DockerImage  string            `json:"docker_image"`
	DockerImages map[string]string `json:"docker_images"`
	Config       *EggConfig        `json:"config"`
	Startup      string            `json:"startup"`
	Script       *EggScript        `json:"script"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type EggConfig struct {
	Files        map[string]any `json:"files"`
	Startup      *EggStartup    `json:"startup"`
	Stop         string         `json:"stop"`
	Logs         any            `json:"logs"`
	FileDenylist []string       `json:"file_denylist"`
	Extends      any            `json:"extends"`
}

type EggStartup struct {
	Done            string   `json:"done"`
	UserInteraction []string `json:"user_interaction,omitempty"`
}

type EggScript struct {
	Privileged bool   `json:"privileged"`
	Install    string `json:"install"`
	Entry      string `json:"entry"`
	Container  string `json:"container"`
	Extends    any    `json:"extends"`
}
