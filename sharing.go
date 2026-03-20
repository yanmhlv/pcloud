package pcloud

import (
	"context"
	"fmt"
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
	Note            string `json:"message,omitempty"`
	ShareName       string `json:"sharename,omitempty"`
	Accepted        bool   `json:"accepted,omitempty"`
	IncomingRequest bool   `json:"incoming,omitempty"`
}

// ShareOpts controls optional parameters for sharing a folder.
type ShareOpts struct {
	Note string
}

type listSharesResponse struct {
	Error
	Shares   []Share `json:"shares"`
	Requests []Share `json:"requests"`
}

func boolParam(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func applyPermissions(params url.Values, perms SharePermissions) {
	params.Set("canread", boolParam(perms.CanRead))
	params.Set("cancreate", boolParam(perms.CanCreate))
	params.Set("canmodify", boolParam(perms.CanModify))
	params.Set("candelete", boolParam(perms.CanDelete))
}

// ShareFolder shares a folder identified by numeric ID with another user by email.
func (c *Client) ShareFolder(ctx context.Context, folderID uint64, email string, perms SharePermissions, opts *ShareOpts) (*Share, error) {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
		"mail":     {email},
	}
	s, err := c.shareFolder(ctx, params, perms, opts)
	if err != nil {
		return nil, fmt.Errorf("share folder %d: %w", folderID, err)
	}
	return s, nil
}

// ShareFolderByPath shares a folder identified by path with another user by email.
func (c *Client) ShareFolderByPath(ctx context.Context, path string, email string, perms SharePermissions, opts *ShareOpts) (*Share, error) {
	params := url.Values{
		"path": {path},
		"mail": {email},
	}
	s, err := c.shareFolder(ctx, params, perms, opts)
	if err != nil {
		return nil, fmt.Errorf("share folder %s: %w", path, err)
	}
	return s, nil
}

func (c *Client) shareFolder(ctx context.Context, params url.Values, perms SharePermissions, opts *ShareOpts) (*Share, error) {
	applyPermissions(params, perms)
	if opts != nil && opts.Note != "" {
		params.Set("message", opts.Note)
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
		return nil, nil, fmt.Errorf("list shares: %w", err)
	}
	return resp.Shares, resp.Requests, nil
}

// AcceptShare accepts an incoming share request by its ID.
func (c *Client) AcceptShare(ctx context.Context, shareRequestID uint64) error {
	params := url.Values{
		"sharerequestid": {strconv.FormatUint(shareRequestID, 10)},
	}

	var resp Error
	if err := c.do(ctx, "acceptshare", params, &resp); err != nil {
		return fmt.Errorf("accept share %d: %w", shareRequestID, err)
	}
	return nil
}

// DeclineShare declines an incoming share request by its ID.
func (c *Client) DeclineShare(ctx context.Context, shareRequestID uint64) error {
	params := url.Values{
		"sharerequestid": {strconv.FormatUint(shareRequestID, 10)},
	}

	var resp Error
	if err := c.do(ctx, "declineshare", params, &resp); err != nil {
		return fmt.Errorf("decline share %d: %w", shareRequestID, err)
	}
	return nil
}

// RemoveShare removes an active share by its ID.
func (c *Client) RemoveShare(ctx context.Context, shareID uint64) error {
	params := url.Values{
		"shareid": {strconv.FormatUint(shareID, 10)},
	}

	var resp Error
	if err := c.do(ctx, "removeshare", params, &resp); err != nil {
		return fmt.Errorf("remove share %d: %w", shareID, err)
	}
	return nil
}

// CancelShareRequest cancels an outgoing share request by its ID.
func (c *Client) CancelShareRequest(ctx context.Context, shareRequestID uint64) error {
	params := url.Values{
		"sharerequestid": {strconv.FormatUint(shareRequestID, 10)},
	}

	var resp Error
	if err := c.do(ctx, "cancelsharerequest", params, &resp); err != nil {
		return fmt.Errorf("cancel share request %d: %w", shareRequestID, err)
	}
	return nil
}

// ChangeShare updates the permissions on an existing share.
func (c *Client) ChangeShare(ctx context.Context, shareID uint64, perms SharePermissions) error {
	params := url.Values{
		"shareid": {strconv.FormatUint(shareID, 10)},
	}
	applyPermissions(params, perms)

	var resp Error
	if err := c.do(ctx, "changeshare", params, &resp); err != nil {
		return fmt.Errorf("change share %d: %w", shareID, err)
	}
	return nil
}
