package pcloud

import (
	"context"
	"net/url"
	"strconv"
)

type revisionsResponse struct {
	Error
	Revisions []Revision `json:"revisions"`
}

func (c *Client) listRevisions(ctx context.Context, params url.Values) ([]Revision, error) {
	var resp revisionsResponse
	if err := c.do(ctx, "listrevisions", params, &resp); err != nil {
		return nil, err
	}
	return resp.Revisions, nil
}

// ListRevisions returns all revisions for a file identified by numeric ID.
func (c *Client) ListRevisions(ctx context.Context, fileID uint64) ([]Revision, error) {
	return c.listRevisions(ctx, url.Values{"fileid": {strconv.FormatUint(fileID, 10)}})
}

// ListRevisionsByPath returns all revisions for a file identified by path.
func (c *Client) ListRevisionsByPath(ctx context.Context, path string) ([]Revision, error) {
	return c.listRevisions(ctx, url.Values{"path": {path}})
}

func (c *Client) revertRevision(ctx context.Context, params url.Values) (*Metadata, error) {
	var resp metadataResponse
	if err := c.do(ctx, "revertrevision", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

// RevertRevision restores a file to a previous revision identified by file and revision IDs.
func (c *Client) RevertRevision(ctx context.Context, fileID, revisionID uint64) (*Metadata, error) {
	return c.revertRevision(ctx, url.Values{
		"fileid":     {strconv.FormatUint(fileID, 10)},
		"revisionid": {strconv.FormatUint(revisionID, 10)},
	})
}

// RevertRevisionByPath restores a file to a previous revision identified by path and revision ID.
func (c *Client) RevertRevisionByPath(ctx context.Context, path string, revisionID uint64) (*Metadata, error) {
	return c.revertRevision(ctx, url.Values{
		"path":       {path},
		"revisionid": {strconv.FormatUint(revisionID, 10)},
	})
}
