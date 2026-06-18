package pcloud

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type PublicLink struct {
	LinkID    uint64   `json:"linkid"`
	Code      string   `json:"code"`
	Link      string   `json:"link"`
	Created   Time     `json:"created"`
	Modified  Time     `json:"modified"`
	Traffic   uint64   `json:"traffic"`
	Downloads uint64   `json:"downloads"`
	Metadata  Metadata `json:"metadata"`
	ShortLink string   `json:"shortlink,omitempty"`
	ShortCode string   `json:"shortcode,omitempty"`
}

type PublicLinkOpts struct {
	MaxDownloads uint64
	MaxTraffic   uint64
	ExpireAt     time.Time
	ShortLink    bool
}

type listPublicLinksResponse struct {
	Error
	PubLinks []PublicLink `json:"publinks"`
}

func applyPublicLinkOpts(params url.Values, opts *PublicLinkOpts) {
	if opts == nil {
		return
	}
	if opts.MaxDownloads > 0 {
		params.Set("maxdownloads", strconv.FormatUint(opts.MaxDownloads, 10))
	}
	if opts.MaxTraffic > 0 {
		params.Set("maxtraffic", strconv.FormatUint(opts.MaxTraffic, 10))
	}
	if !opts.ExpireAt.IsZero() {
		params.Set("expire", strconv.FormatInt(opts.ExpireAt.Unix(), 10))
	}
	if opts.ShortLink {
		params.Set("shortlink", "1")
	}
}

type publicLinkResponse struct {
	Error
	PublicLink
}

func (c *Client) createFilePublicLink(ctx context.Context, params url.Values, opts *PublicLinkOpts) (*PublicLink, error) {
	applyPublicLinkOpts(params, opts)
	var resp publicLinkResponse
	if err := c.doGet(ctx, "getfilepublink", params, &resp); err != nil {
		return nil, err
	}
	return &resp.PublicLink, nil
}

func (c *Client) CreateFilePublicLink(ctx context.Context, fileID uint64, opts *PublicLinkOpts) (*PublicLink, error) {
	pl, err := c.createFilePublicLink(ctx, url.Values{paramFileID: {formatUint(fileID)}}, opts)
	if err != nil {
		return nil, fmt.Errorf("create file public link %d: %w", fileID, err)
	}
	return pl, nil
}

func (c *Client) CreateFilePublicLinkByPath(ctx context.Context, path string, opts *PublicLinkOpts) (*PublicLink, error) {
	pl, err := c.createFilePublicLink(ctx, url.Values{paramPath: {path}}, opts)
	if err != nil {
		return nil, fmt.Errorf("create file public link %s: %w", path, err)
	}
	return pl, nil
}

func (c *Client) createFolderPublicLink(ctx context.Context, params url.Values, opts *PublicLinkOpts) (*PublicLink, error) {
	applyPublicLinkOpts(params, opts)
	var resp publicLinkResponse
	if err := c.doGet(ctx, "getfolderpublink", params, &resp); err != nil {
		return nil, err
	}
	return &resp.PublicLink, nil
}

func (c *Client) CreateFolderPublicLink(ctx context.Context, folderID uint64, opts *PublicLinkOpts) (*PublicLink, error) {
	pl, err := c.createFolderPublicLink(ctx, url.Values{paramFolderID: {formatUint(folderID)}}, opts)
	if err != nil {
		return nil, fmt.Errorf("create folder public link %d: %w", folderID, err)
	}
	return pl, nil
}

func (c *Client) CreateFolderPublicLinkByPath(ctx context.Context, path string, opts *PublicLinkOpts) (*PublicLink, error) {
	pl, err := c.createFolderPublicLink(ctx, url.Values{paramPath: {path}}, opts)
	if err != nil {
		return nil, fmt.Errorf("create folder public link %s: %w", path, err)
	}
	return pl, nil
}

func (c *Client) ListPublicLinks(ctx context.Context) ([]PublicLink, error) {
	var resp listPublicLinksResponse
	if err := c.doGet(ctx, "listpublinks", url.Values{}, &resp); err != nil {
		return nil, fmt.Errorf("list public links: %w", err)
	}
	return resp.PubLinks, nil
}

func (c *Client) DeletePublicLink(ctx context.Context, linkID uint64) error {
	params := url.Values{
		paramLinkID: {formatUint(linkID)},
	}

	var resp Error
	if err := c.doGet(ctx, "deletepublink", params, &resp); err != nil {
		return fmt.Errorf("delete public link %d: %w", linkID, err)
	}
	return nil
}

func (c *Client) ChangePublicLink(ctx context.Context, linkID uint64, opts *PublicLinkOpts) (*PublicLink, error) {
	params := url.Values{
		paramLinkID: {formatUint(linkID)},
	}
	applyPublicLinkOpts(params, opts)

	var resp publicLinkResponse
	if err := c.doGet(ctx, "changepublink", params, &resp); err != nil {
		return nil, fmt.Errorf("change public link %d: %w", linkID, err)
	}
	return &resp.PublicLink, nil
}

func (c *Client) GetPublicLinkInfo(ctx context.Context, code string) (*PublicLink, error) {
	params := url.Values{
		"code": {code},
	}

	var resp publicLinkResponse
	if err := c.doGet(ctx, "showpublink", params, &resp); err != nil {
		return nil, fmt.Errorf("get public link info %s: %w", code, err)
	}
	return &resp.PublicLink, nil
}
