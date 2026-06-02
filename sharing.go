package pcloud

import (
	"context"
	"fmt"
	"net/url"
)

type SharePermissions struct {
	CanRead   bool
	CanCreate bool
	CanModify bool
	CanDelete bool
}

type Share struct {
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

type ShareOpts struct {
	Note string
}

type shareResponse struct {
	Error
	Share
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

func (c *Client) ShareFolder(ctx context.Context, folderID uint64, email string, perms SharePermissions, opts *ShareOpts) (*Share, error) {
	params := url.Values{
		paramFolderID: {formatUint(folderID)},
		paramMail:     {email},
	}
	s, err := c.shareFolder(ctx, params, perms, opts)
	if err != nil {
		return nil, fmt.Errorf("share folder %d: %w", folderID, err)
	}
	return s, nil
}

func (c *Client) ShareFolderByPath(ctx context.Context, path string, email string, perms SharePermissions, opts *ShareOpts) (*Share, error) {
	params := url.Values{
		paramPath: {path},
		paramMail: {email},
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

	var resp shareResponse
	if err := c.do(ctx, "sharefolder", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Share, nil
}

type Shares struct {
	Active   []Share
	Requests []Share
}

func (c *Client) ListShares(ctx context.Context) (Shares, error) {
	var resp listSharesResponse
	if err := c.do(ctx, "listshares", url.Values{}, &resp); err != nil {
		return Shares{}, fmt.Errorf("list shares: %w", err)
	}
	return Shares{Active: resp.Shares, Requests: resp.Requests}, nil
}

func (c *Client) AcceptShare(ctx context.Context, shareRequestID uint64) error {
	params := url.Values{
		paramShareRequestID: {formatUint(shareRequestID)},
	}

	var resp Error
	if err := c.do(ctx, "acceptshare", params, &resp); err != nil {
		return fmt.Errorf("accept share %d: %w", shareRequestID, err)
	}
	return nil
}

func (c *Client) DeclineShare(ctx context.Context, shareRequestID uint64) error {
	params := url.Values{
		paramShareRequestID: {formatUint(shareRequestID)},
	}

	var resp Error
	if err := c.do(ctx, "declineshare", params, &resp); err != nil {
		return fmt.Errorf("decline share %d: %w", shareRequestID, err)
	}
	return nil
}

func (c *Client) RemoveShare(ctx context.Context, shareID uint64) error {
	params := url.Values{
		paramShareID: {formatUint(shareID)},
	}

	var resp Error
	if err := c.do(ctx, "removeshare", params, &resp); err != nil {
		return fmt.Errorf("remove share %d: %w", shareID, err)
	}
	return nil
}

func (c *Client) CancelShareRequest(ctx context.Context, shareRequestID uint64) error {
	params := url.Values{
		paramShareRequestID: {formatUint(shareRequestID)},
	}

	var resp Error
	if err := c.do(ctx, "cancelsharerequest", params, &resp); err != nil {
		return fmt.Errorf("cancel share request %d: %w", shareRequestID, err)
	}
	return nil
}

func (c *Client) ChangeShare(ctx context.Context, shareID uint64, perms SharePermissions) error {
	params := url.Values{
		paramShareID: {formatUint(shareID)},
	}
	applyPermissions(params, perms)

	var resp Error
	if err := c.do(ctx, "changeshare", params, &resp); err != nil {
		return fmt.Errorf("change share %d: %w", shareID, err)
	}
	return nil
}
