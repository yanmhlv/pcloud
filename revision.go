package pcloud

import (
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

func (c *Client) ListRevisions(fileID uint64) ([]Revision, error) {
	params := url.Values{
		"fileid": {strconv.FormatUint(fileID, 10)},
	}

	var resp revisionsResponse
	if err := c.do("listrevisions", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return resp.Revisions, nil
}

func (c *Client) ListRevisionsByPath(path string) ([]Revision, error) {
	params := url.Values{
		"path": {path},
	}

	var resp revisionsResponse
	if err := c.do("listrevisions", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return resp.Revisions, nil
}

func (c *Client) RevertRevision(fileID, revisionID uint64) (*Metadata, error) {
	params := url.Values{
		"fileid":     {strconv.FormatUint(fileID, 10)},
		"revisionid": {strconv.FormatUint(revisionID, 10)},
	}

	var resp revertResponse
	if err := c.do("revertrevision", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) RevertRevisionByPath(path string, revisionID uint64) (*Metadata, error) {
	params := url.Values{
		"path":       {path},
		"revisionid": {strconv.FormatUint(revisionID, 10)},
	}

	var resp revertResponse
	if err := c.do("revertrevision", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}
