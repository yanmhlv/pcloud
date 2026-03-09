package pcloud

import (
	"context"
	"net/url"
	"strconv"
)

// SearchOpts controls optional parameters for file searches.
type SearchOpts struct {
	FolderID  uint64
	Recursive bool
}

type searchResponse struct {
	Error
	Items []Metadata `json:"items"`
}

// Search searches for files matching query. Pass nil opts for a global search.
func (c *Client) Search(ctx context.Context, query string, opts *SearchOpts) ([]Metadata, error) {
	params := url.Values{"query": {query}}
	if opts != nil {
		if opts.FolderID > 0 {
			params.Set("folderid", strconv.FormatUint(opts.FolderID, 10))
		}
		if opts.Recursive {
			params.Set("recursive", "1")
		}
	}

	var resp searchResponse
	if err := c.do(ctx, "searchfiles", params, &resp); err != nil {
		return nil, err
	}
	return resp.Items, nil
}
