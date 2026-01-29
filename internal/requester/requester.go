package requester

import (
	"context"
	"github.com/vovamod/go-pterodactyl/api"
	"io"
	"net/http"
)

type Requester interface {
	NewRequest(ctx context.Context, method, endpoint string, body io.Reader, options *api.PaginationOptions) (*http.Request, error)
	Do(ctx context.Context, req *http.Request, v any) (*http.Response, error)
}
