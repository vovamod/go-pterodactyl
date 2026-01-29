package crud

import (
	"context"
	"fmt"

	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/internal/requester"
)

func List[T any](ctx context.Context, c requester.Requester, path string,
	opt *api.PaginationOptions) ([]*T, *api.Meta, error) {

	req, err := c.NewRequest(ctx, "GET", path, nil, opt)
	if err != nil {
		return nil, nil, err
	}

	resp := &api.PaginatedResponse[T]{}
	if _, err = c.Do(ctx, req, resp); err != nil {
		return nil, nil, err
	}

	out := make([]*T, len(resp.Data))
	for i, item := range resp.Data {
		out[i] = item.Attributes
	}
	return out, &resp.Meta, nil
}

func ListAll[T any](
	ctx context.Context,
	c requester.Requester,
	path string,
	perPage int,
) ([]*T, error) {

	if perPage <= 0 {
		perPage = 100
	}

	all := make([]*T, 0, perPage)
	opts := &api.PaginationOptions{PerPage: perPage, Page: 1}

	for {
		items, meta, err := List[T](ctx, c, path, opts)
		if err != nil {
			return nil, err
		}
		all = append(all, items...)

		if meta.Pagination.CurrentPage >= meta.Pagination.TotalPages {
			break
		}
		opts.Page++
	}

	return all, nil
}

func Get[T any](ctx context.Context, c requester.Requester, path string, id int) (*T, error) {
	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s/%d", path, id), nil, nil)
	if err != nil {
		return nil, err
	}

	resp := &api.ListItem[T]{}
	if _, err = c.Do(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp.Attributes, nil
}

func Delete[T any](ctx context.Context, c requester.Requester, path string, id int) error {
	req, err := c.NewRequest(ctx, "DELETE", fmt.Sprintf("%s/%d", path, id), nil, nil)
	if err != nil {
		return err
	}
	_, err = c.Do(ctx, req, nil)
	return err
}
