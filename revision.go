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

type revertResponse struct {
	Error
	Metadata Metadata `json:"metadata"`
}

func (c *Client) ListRevisions(ctx context.Context, fileID uint64) ([]Revision, error) {
	params := url.Values{
		"fileid": {strconv.FormatUint(fileID, 10)},
	}

	var resp revisionsResponse
	if err := c.do(ctx, "listrevisions", params, &resp); err != nil {
		return nil, err
	}
	return resp.Revisions, nil
}

func (c *Client) ListRevisionsByPath(ctx context.Context, path string) ([]Revision, error) {
	params := url.Values{
		"path": {path},
	}

	var resp revisionsResponse
	if err := c.do(ctx, "listrevisions", params, &resp); err != nil {
		return nil, err
	}
	return resp.Revisions, nil
}

func (c *Client) RevertRevision(ctx context.Context, fileID, revisionID uint64) (*Metadata, error) {
	params := url.Values{
		"fileid":     {strconv.FormatUint(fileID, 10)},
		"revisionid": {strconv.FormatUint(revisionID, 10)},
	}

	var resp revertResponse
	if err := c.do(ctx, "revertrevision", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) RevertRevisionByPath(ctx context.Context, path string, revisionID uint64) (*Metadata, error) {
	params := url.Values{
		"path":       {path},
		"revisionid": {strconv.FormatUint(revisionID, 10)},
	}

	var resp revertResponse
	if err := c.do(ctx, "revertrevision", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}
