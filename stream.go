package pcloud

import (
	"context"
	"net/url"
	"strconv"
)

type FileLinkOpts struct {
	ForceDownload bool
	ContentType   string
	MaxSpeed      int
}

func (c *Client) GetFileLink(ctx context.Context, fileID uint64) (*FileLink, error) {
	return c.GetFileLinkWithOpts(ctx, fileID, nil)
}

func (c *Client) GetFileLinkByPath(ctx context.Context, path string) (*FileLink, error) {
	return c.GetFileLinkByPathWithOpts(ctx, path, nil)
}

func (c *Client) GetFileLinkWithOpts(ctx context.Context, fileID uint64, opts *FileLinkOpts) (*FileLink, error) {
	params := url.Values{
		"fileid": {strconv.FormatUint(fileID, 10)},
	}
	applyLinkOpts(params, opts)

	var resp FileLink
	if err := c.do(ctx, "getfilelink", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetFileLinkByPathWithOpts(ctx context.Context, path string, opts *FileLinkOpts) (*FileLink, error) {
	params := url.Values{
		"path": {path},
	}
	applyLinkOpts(params, opts)

	var resp FileLink
	if err := c.do(ctx, "getfilelink", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) getMediaLink(ctx context.Context, fileID uint64, endpoint string) (*FileLink, error) {
	params := url.Values{
		"fileid": {strconv.FormatUint(fileID, 10)},
	}

	var resp FileLink
	if err := c.do(ctx, endpoint, params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetVideoLink(ctx context.Context, fileID uint64) (*FileLink, error) {
	return c.getMediaLink(ctx, fileID, "getvideolink")
}

func (c *Client) GetAudioLink(ctx context.Context, fileID uint64) (*FileLink, error) {
	return c.getMediaLink(ctx, fileID, "getaudiolink")
}

func (c *Client) GetHLSLink(ctx context.Context, fileID uint64) (*FileLink, error) {
	return c.getMediaLink(ctx, fileID, "gethlslink")
}

func applyLinkOpts(params url.Values, opts *FileLinkOpts) {
	if opts == nil {
		return
	}
	if opts.ForceDownload {
		params.Set("forcedownload", "1")
	}
	if opts.ContentType != "" {
		params.Set("contenttype", opts.ContentType)
	}
	if opts.MaxSpeed > 0 {
		params.Set("maxspeed", strconv.Itoa(opts.MaxSpeed))
	}
}
