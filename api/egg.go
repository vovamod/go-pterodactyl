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

type EggVariable struct {
	ID           int    `json:"id"`
	EggID        int    `json:"egg_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	EnvVariable  string `json:"env_variable"`
	DefaultValue string `json:"default_value"`
	UserViewable bool   `json:"user_viewable"`
	UserEditable bool   `json:"user_editable"`
	Rules        string `json:"rules"`
}

type EggWithVariablesResponse struct {
	Attributes    *Egg `json:"attributes"`
	Relationships struct {
		Variables struct {
			Data []ListItem[EggVariable] `json:"data"`
		} `json:"variables"`
	} `json:"relationships"`
}
