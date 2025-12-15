package pcloud

import (
	"context"
	"net/url"
	"strconv"
)

type PublicLink struct {
	Error
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
	MaxDownloads int
	MaxTraffic   uint64
	ExpireTime   int64
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
		params.Set("maxdownloads", strconv.Itoa(opts.MaxDownloads))
	}
	if opts.MaxTraffic > 0 {
		params.Set("maxtraffic", strconv.FormatUint(opts.MaxTraffic, 10))
	}
	if opts.ExpireTime > 0 {
		params.Set("expire", strconv.FormatInt(opts.ExpireTime, 10))
	}
	if opts.ShortLink {
		params.Set("shortlink", "1")
	}
}

func (c *Client) CreateFilePublicLink(ctx context.Context, fileID uint64, opts *PublicLinkOpts) (*PublicLink, error) {
	params := url.Values{
		"fileid": {strconv.FormatUint(fileID, 10)},
	}
	applyPublicLinkOpts(params, opts)

	var resp PublicLink
	if err := c.do(ctx, "getfilepublink", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateFilePublicLinkByPath(ctx context.Context, path string, opts *PublicLinkOpts) (*PublicLink, error) {
	params := url.Values{
		"path": {path},
	}
	applyPublicLinkOpts(params, opts)

	var resp PublicLink
	if err := c.do(ctx, "getfilepublink", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateFolderPublicLink(ctx context.Context, folderID uint64, opts *PublicLinkOpts) (*PublicLink, error) {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
	}
	applyPublicLinkOpts(params, opts)

	var resp PublicLink
	if err := c.do(ctx, "getfolderpublink", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateFolderPublicLinkByPath(ctx context.Context, path string, opts *PublicLinkOpts) (*PublicLink, error) {
	params := url.Values{
		"path": {path},
	}
	applyPublicLinkOpts(params, opts)

	var resp PublicLink
	if err := c.do(ctx, "getfolderpublink", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) ListPublicLinks(ctx context.Context) ([]PublicLink, error) {
	var resp listPublicLinksResponse
	if err := c.do(ctx, "listpublinks", url.Values{}, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return resp.PubLinks, nil
}

func (c *Client) DeletePublicLink(ctx context.Context, linkID uint64) error {
	params := url.Values{
		"linkid": {strconv.FormatUint(linkID, 10)},
	}

	var resp Error
	if err := c.do(ctx, "deletepublink", params, &resp); err != nil {
		return err
	}
	return resp.Err()
}

func (c *Client) ChangePublicLink(ctx context.Context, linkID uint64, opts *PublicLinkOpts) (*PublicLink, error) {
	params := url.Values{
		"linkid": {strconv.FormatUint(linkID, 10)},
	}
	applyPublicLinkOpts(params, opts)

	var resp PublicLink
	if err := c.do(ctx, "changepublink", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetPublicLinkInfo(ctx context.Context, code string) (*PublicLink, error) {
	params := url.Values{
		"code": {code},
	}

	var resp PublicLink
	if err := c.do(ctx, "showpublink", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp, nil
}
