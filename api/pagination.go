package api

type PaginatedResponse[T any] struct {
	Object string         `json:"object"` // Will be "list"
	Data   []*ListItem[T] `json:"data"`
	Meta   Meta           `json:"meta"`
}

type ListItem[T any] struct {
	Object     string `json:"object"` // e.g., "user", "server"
	Attributes *T     `json:"attributes"`
}

type Meta struct {
	Pagination Pagination `json:"pagination"`
}

// Pagination holds the detailed page information.
type Pagination struct {
	Total       int `json:"total"`
	Count       int `json:"count"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	// Links can be complex, often just an empty object `{}`, so `any` is safe.
	Links any `json:"links,omitempty"`
}
type PaginationOptions struct {
	Page    int      `json:"page"`
	PerPage int      `json:"per_page"`
	Include []string `json:"include,omitempty"`

	FilterEmail      string `json:"filter_email,omitempty"`
	FilterUUID       string `json:"filter_uuid,omitempty"`
	FilterUsername   string `json:"filter_username,omitempty"`
	FilterExternalID string `json:"filter_external_id,omitempty"`
	Sort             string `json:"sort,omitempty"`
}
