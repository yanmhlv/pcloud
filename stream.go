package pcloud

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type FileLinkOpts struct {
	ForceDownload bool
	ContentType   string
	MaxSpeed      uint64
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
		params.Set("maxspeed", strconv.FormatUint(opts.MaxSpeed, 10))
	}
}

type fileLinkResponse struct {
	Error
	FileLink
}

func (c *Client) getFileLink(ctx context.Context, params url.Values, opts *FileLinkOpts) (*FileLink, error) {
	applyLinkOpts(params, opts)
	var resp fileLinkResponse
	if err := c.doGet(ctx, "getfilelink", params, &resp); err != nil {
		return nil, err
	}
	return &resp.FileLink, nil
}

func (c *Client) GetFileLink(ctx context.Context, fileID uint64, opts *FileLinkOpts) (*FileLink, error) {
	fl, err := c.getFileLink(ctx, url.Values{paramFileID: {formatUint(fileID)}}, opts)
	if err != nil {
		return nil, fmt.Errorf("get file link %d: %w", fileID, err)
	}
	return fl, nil
}

func (c *Client) GetFileLinkByPath(ctx context.Context, path string, opts *FileLinkOpts) (*FileLink, error) {
	fl, err := c.getFileLink(ctx, url.Values{paramPath: {path}}, opts)
	if err != nil {
		return nil, fmt.Errorf("get file link %s: %w", path, err)
	}
	return fl, nil
}

func (c *Client) getMediaLink(ctx context.Context, fileID uint64, endpoint string) (*FileLink, error) {
	params := url.Values{
		paramFileID: {formatUint(fileID)},
	}

	var resp fileLinkResponse
	if err := c.doGet(ctx, endpoint, params, &resp); err != nil {
		return nil, err
	}
	return &resp.FileLink, nil
}

func (c *Client) GetVideoLink(ctx context.Context, fileID uint64) (*FileLink, error) {
	fl, err := c.getMediaLink(ctx, fileID, "getvideolink")
	if err != nil {
		return nil, fmt.Errorf("get video link %d: %w", fileID, err)
	}
	return fl, nil
}

func (c *Client) GetAudioLink(ctx context.Context, fileID uint64) (*FileLink, error) {
	fl, err := c.getMediaLink(ctx, fileID, "getaudiolink")
	if err != nil {
		return nil, fmt.Errorf("get audio link %d: %w", fileID, err)
	}
	return fl, nil
}

func (c *Client) GetHLSLink(ctx context.Context, fileID uint64) (*FileLink, error) {
	fl, err := c.getMediaLink(ctx, fileID, "gethlslink")
	if err != nil {
		return nil, fmt.Errorf("get hls link %d: %w", fileID, err)
	}
	return fl, nil
}
