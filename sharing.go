package pcloud

import (
	"context"
	"net/url"
	"strconv"
)

type SharePermissions struct {
	CanRead   bool
	CanCreate bool
	CanModify bool
	CanDelete bool
}

type Share struct {
	Error
	ShareID         uint64 `json:"shareid"`
	ShareRequestID  uint64 `json:"sharerequestid"`
	FolderID        uint64 `json:"folderid"`
	ToEmail         string `json:"tomail"`
	ToUserID        uint64 `json:"touserid"`
	FromUserID      uint64 `json:"fromuserid"`
	CanRead         bool   `json:"canread"`
	CanCreate       bool   `json:"cancreate"`
	CanModify       bool   `json:"canmodify"`
	CanDelete       bool   `json:"candelete"`
	Created         Time   `json:"created"`
	Message         string `json:"message,omitempty"`
	ShareName       string `json:"sharename,omitempty"`
	Accepted        bool   `json:"accepted,omitempty"`
	IncomingRequest bool   `json:"incoming,omitempty"`
}

type ShareOpts struct {
	Message string
}

type listSharesResponse struct {
	Error
	Shares   []Share `json:"shares"`
	Requests []Share `json:"requests"`
}

func permissionsToAccess(perms SharePermissions) string {
	access := 0
	if perms.CanModify {
		access = 2
	} else if perms.CanCreate {
		access = 1
	}
	return strconv.Itoa(access)
}

func (c *Client) ShareFolder(ctx context.Context, folderID uint64, email string, perms SharePermissions, opts *ShareOpts) (*Share, error) {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
		"mail":     {email},
	}

	if perms.CanRead {
		params.Set("canread", "1")
	}
	if perms.CanCreate {
		params.Set("cancreate", "1")
	}
	if perms.CanModify {
		params.Set("canmodify", "1")
	}
	if perms.CanDelete {
		params.Set("candelete", "1")
	}
	if opts != nil && opts.Message != "" {
		params.Set("message", opts.Message)
	}

	var resp Share
	if err := c.do(ctx, "sharefolder", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) ShareFolderByPath(ctx context.Context, path string, email string, perms SharePermissions, opts *ShareOpts) (*Share, error) {
	params := url.Values{
		"path": {path},
		"mail": {email},
	}

	if perms.CanRead {
		params.Set("canread", "1")
	}
	if perms.CanCreate {
		params.Set("cancreate", "1")
	}
	if perms.CanModify {
		params.Set("canmodify", "1")
	}
	if perms.CanDelete {
		params.Set("candelete", "1")
	}
	if opts != nil && opts.Message != "" {
		params.Set("message", opts.Message)
	}

	var resp Share
	if err := c.do(ctx, "sharefolder", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) ListShares(ctx context.Context) ([]Share, []Share, error) {
	var resp listSharesResponse
	if err := c.do(ctx, "listshares", url.Values{}, &resp); err != nil {
		return nil, nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, nil, err
	}
	return resp.Shares, resp.Requests, nil
}

func (c *Client) AcceptShare(ctx context.Context, shareRequestID uint64) error {
	params := url.Values{
		"sharerequestid": {strconv.FormatUint(shareRequestID, 10)},
	}

	var resp Error
	if err := c.do(ctx, "acceptshare", params, &resp); err != nil {
		return err
	}
	return resp.Err()
}

func (c *Client) DeclineShare(ctx context.Context, shareRequestID uint64) error {
	params := url.Values{
		"sharerequestid": {strconv.FormatUint(shareRequestID, 10)},
	}

	var resp Error
	if err := c.do(ctx, "declineshare", params, &resp); err != nil {
		return err
	}
	return resp.Err()
}

func (c *Client) RemoveShare(ctx context.Context, shareID uint64) error {
	params := url.Values{
		"shareid": {strconv.FormatUint(shareID, 10)},
	}

	var resp Error
	if err := c.do(ctx, "removeshare", params, &resp); err != nil {
		return err
	}
	return resp.Err()
}

func (c *Client) CancelShareRequest(ctx context.Context, shareRequestID uint64) error {
	params := url.Values{
		"sharerequestid": {strconv.FormatUint(shareRequestID, 10)},
	}

	var resp Error
	if err := c.do(ctx, "cancelsharerequest", params, &resp); err != nil {
		return err
	}
	return resp.Err()
}

func (c *Client) ChangeShare(ctx context.Context, shareID uint64, perms SharePermissions) error {
	params := url.Values{
		"shareid": {strconv.FormatUint(shareID, 10)},
	}

	if perms.CanRead {
		params.Set("canread", "1")
	} else {
		params.Set("canread", "0")
	}
	if perms.CanCreate {
		params.Set("cancreate", "1")
	} else {
		params.Set("cancreate", "0")
	}
	if perms.CanModify {
		params.Set("canmodify", "1")
	} else {
		params.Set("canmodify", "0")
	}
	if perms.CanDelete {
		params.Set("candelete", "1")
	} else {
		params.Set("candelete", "0")
	}

	var resp Error
	if err := c.do(ctx, "changeshare", params, &resp); err != nil {
		return err
	}
	return resp.Err()
}
