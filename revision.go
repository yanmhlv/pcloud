package pcloud

import (
	"context"
	"fmt"
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

func (c *Client) ListRevisions(ctx context.Context, fileID uint64) ([]Revision, error) {
	r, err := c.listRevisions(ctx, url.Values{"fileid": {strconv.FormatUint(fileID, 10)}})
	if err != nil {
		return nil, fmt.Errorf("list revisions %d: %w", fileID, err)
	}
	return r, nil
}

func (c *Client) ListRevisionsByPath(ctx context.Context, path string) ([]Revision, error) {
	r, err := c.listRevisions(ctx, url.Values{"path": {path}})
	if err != nil {
		return nil, fmt.Errorf("list revisions %s: %w", path, err)
	}
	return r, nil
}

func (c *Client) revertRevision(ctx context.Context, params url.Values) (*Metadata, error) {
	var resp metadataResponse
	if err := c.do(ctx, "revertrevision", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) RevertRevision(ctx context.Context, fileID, revisionID uint64) (*Metadata, error) {
	m, err := c.revertRevision(ctx, url.Values{
		"fileid":     {strconv.FormatUint(fileID, 10)},
		"revisionid": {strconv.FormatUint(revisionID, 10)},
	})
	if err != nil {
		return nil, fmt.Errorf("revert revision %d/%d: %w", fileID, revisionID, err)
	}
	return m, nil
}

func (c *Client) RevertRevisionByPath(ctx context.Context, path string, revisionID uint64) (*Metadata, error) {
	m, err := c.revertRevision(ctx, url.Values{
		"path":       {path},
		"revisionid": {strconv.FormatUint(revisionID, 10)},
	})
	if err != nil {
		return nil, fmt.Errorf("revert revision %s/%d: %w", path, revisionID, err)
	}
	return m, nil
}
