package pcloud

import (
	"context"
	"net/url"
	"strconv"
)

// FileLinkOpts controls optional parameters for file link requests.
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

func (c *Client) getFileLink(ctx context.Context, params url.Values, opts *FileLinkOpts) (*FileLink, error) {
	applyLinkOpts(params, opts)
	var resp FileLink
	if err := c.do(ctx, "getfilelink", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetFileLink returns a temporary direct-download URL for a file by numeric ID.
func (c *Client) GetFileLink(ctx context.Context, fileID uint64) (*FileLink, error) {
	return c.getFileLink(ctx, url.Values{"fileid": {strconv.FormatUint(fileID, 10)}}, nil)
}

// GetFileLinkByPath returns a temporary direct-download URL for a file by path.
func (c *Client) GetFileLinkByPath(ctx context.Context, path string) (*FileLink, error) {
	return c.getFileLink(ctx, url.Values{"path": {path}}, nil)
}

// GetFileLinkWithOpts returns a direct-download URL for a file by numeric ID with options.
func (c *Client) GetFileLinkWithOpts(ctx context.Context, fileID uint64, opts *FileLinkOpts) (*FileLink, error) {
	return c.getFileLink(ctx, url.Values{"fileid": {strconv.FormatUint(fileID, 10)}}, opts)
}

// GetFileLinkByPathWithOpts returns a direct-download URL for a file by path with options.
func (c *Client) GetFileLinkByPathWithOpts(ctx context.Context, path string, opts *FileLinkOpts) (*FileLink, error) {
	return c.getFileLink(ctx, url.Values{"path": {path}}, opts)
}

func (c *Client) getMediaLink(ctx context.Context, fileID uint64, endpoint string) (*FileLink, error) {
	params := url.Values{
		"fileid": {strconv.FormatUint(fileID, 10)},
	}

	var resp FileLink
	if err := c.do(ctx, endpoint, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetVideoLink returns a streaming URL for a video file.
func (c *Client) GetVideoLink(ctx context.Context, fileID uint64) (*FileLink, error) {
	return c.getMediaLink(ctx, fileID, "getvideolink")
}

// GetAudioLink returns a streaming URL for an audio file.
func (c *Client) GetAudioLink(ctx context.Context, fileID uint64) (*FileLink, error) {
	return c.getMediaLink(ctx, fileID, "getaudiolink")
}

// GetHLSLink returns an HLS streaming URL for a video file.
func (c *Client) GetHLSLink(ctx context.Context, fileID uint64) (*FileLink, error) {
	return c.getMediaLink(ctx, fileID, "gethlslink")
}
