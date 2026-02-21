package pcloud

import (
	"context"
	"net/url"
	"strconv"
)

// SharePermissions defines the access rights granted in a folder share.
type SharePermissions struct {
	CanRead   bool
	CanCreate bool
	CanModify bool
	CanDelete bool
}

// Share represents a folder sharing record, including both active shares and pending requests.
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

// ShareOpts controls optional parameters for sharing a folder.
type ShareOpts struct {
	Message string
}

type listSharesResponse struct {
	Error
	Shares   []Share `json:"shares"`
	Requests []Share `json:"requests"`
}

func applyPermissions(params url.Values, perms SharePermissions) {
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
}

// ShareFolder shares a folder identified by numeric ID with another user by email.
func (c *Client) ShareFolder(ctx context.Context, folderID uint64, email string, perms SharePermissions, opts *ShareOpts) (*Share, error) {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
		"mail":     {email},
	}
	return c.shareFolder(ctx, params, perms, opts)
}

// ShareFolderByPath shares a folder identified by path with another user by email.
func (c *Client) ShareFolderByPath(ctx context.Context, path string, email string, perms SharePermissions, opts *ShareOpts) (*Share, error) {
	params := url.Values{
		"path": {path},
		"mail": {email},
	}
	return c.shareFolder(ctx, params, perms, opts)
}

func (c *Client) shareFolder(ctx context.Context, params url.Values, perms SharePermissions, opts *ShareOpts) (*Share, error) {
	applyPermissions(params, perms)
	if opts != nil && opts.Message != "" {
		params.Set("message", opts.Message)
	}

	var resp Share
	if err := c.do(ctx, "sharefolder", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListShares returns active shares and pending share requests for the authenticated user.
func (c *Client) ListShares(ctx context.Context) ([]Share, []Share, error) {
	var resp listSharesResponse
	if err := c.do(ctx, "listshares", url.Values{}, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Shares, resp.Requests, nil
}

// AcceptShare accepts an incoming share request by its ID.
func (c *Client) AcceptShare(ctx context.Context, shareRequestID uint64) error {
	params := url.Values{
		"sharerequestid": {strconv.FormatUint(shareRequestID, 10)},
	}

	var resp Error
	return c.do(ctx, "acceptshare", params, &resp)
}

// DeclineShare declines an incoming share request by its ID.
func (c *Client) DeclineShare(ctx context.Context, shareRequestID uint64) error {
	params := url.Values{
		"sharerequestid": {strconv.FormatUint(shareRequestID, 10)},
	}

	var resp Error
	return c.do(ctx, "declineshare", params, &resp)
}

// RemoveShare removes an active share by its ID.
func (c *Client) RemoveShare(ctx context.Context, shareID uint64) error {
	params := url.Values{
		"shareid": {strconv.FormatUint(shareID, 10)},
	}

	var resp Error
	return c.do(ctx, "removeshare", params, &resp)
}

// CancelShareRequest cancels an outgoing share request by its ID.
func (c *Client) CancelShareRequest(ctx context.Context, shareRequestID uint64) error {
	params := url.Values{
		"sharerequestid": {strconv.FormatUint(shareRequestID, 10)},
	}

	var resp Error
	return c.do(ctx, "cancelsharerequest", params, &resp)
}

// ChangeShare updates the permissions on an existing share.
func (c *Client) ChangeShare(ctx context.Context, shareID uint64, perms SharePermissions) error {
	params := url.Values{
		"shareid": {strconv.FormatUint(shareID, 10)},
	}
	applyPermissions(params, perms)

	var resp Error
	return c.do(ctx, "changeshare", params, &resp)
}
