package pcloud

import (
	"context"
	"fmt"
	"net/url"
)

type SearchOpts struct {
	FolderID  uint64
	Recursive bool
}

type searchResponse struct {
	Error
	Items []Metadata `json:"items"`
}

func (c *Client) Search(ctx context.Context, query string, opts *SearchOpts) ([]Metadata, error) {
	params := url.Values{"query": {query}}
	if opts != nil {
		if opts.FolderID > 0 {
			params.Set(paramFolderID, formatUint(opts.FolderID))
		}
		if opts.Recursive {
			params.Set("recursive", "1")
		}
	}

	var resp searchResponse
	if err := c.do(ctx, "searchfiles", params, &resp); err != nil {
		return nil, fmt.Errorf("search %q: %w", query, err)
	}
	return resp.Items, nil
}
