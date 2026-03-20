package pcloud

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ThumbOpts controls optional parameters for thumbnail requests.
type ThumbOpts struct {
	Crop bool
	Type string // "png" or "jpeg"; defaults to "png" if empty
}

func applyThumbOpts(params url.Values, opts *ThumbOpts) {
	if opts == nil {
		return
	}
	if opts.Crop {
		params.Set("crop", "1")
	}
	if opts.Type != "" {
		params.Set("type", opts.Type)
	}
}

func validateThumbSize(width, height int) error {
	if width < 1 || width > 2048 || height < 1 || height > 2048 {
		return fmt.Errorf("thumbnail size %dx%d is out of allowed range 1–2048", width, height)
	}
	return nil
}

func (c *Client) getThumbnail(ctx context.Context, params url.Values, width, height int, opts *ThumbOpts) (*FileLink, error) {
	if err := validateThumbSize(width, height); err != nil {
		return nil, err
	}
	params.Set("size", fmt.Sprintf("%dx%d", width, height))
	applyThumbOpts(params, opts)

	var resp FileLink
	if err := c.do(ctx, "getthumb", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetThumbnail returns a thumbnail download link for a file by numeric ID.
// Width and height must each be between 1 and 2048 pixels.
func (c *Client) GetThumbnail(ctx context.Context, fileID uint64, width, height int, opts *ThumbOpts) (*FileLink, error) {
	fl, err := c.getThumbnail(ctx, url.Values{"fileid": {strconv.FormatUint(fileID, 10)}}, width, height, opts)
	if err != nil {
		return nil, fmt.Errorf("get thumbnail %d: %w", fileID, err)
	}
	return fl, nil
}

// GetThumbnailByPath returns a thumbnail download link for a file by path.
// Width and height must each be between 1 and 2048 pixels.
func (c *Client) GetThumbnailByPath(ctx context.Context, path string, width, height int, opts *ThumbOpts) (*FileLink, error) {
	fl, err := c.getThumbnail(ctx, url.Values{"path": {path}}, width, height, opts)
	if err != nil {
		return nil, fmt.Errorf("get thumbnail %s: %w", path, err)
	}
	return fl, nil
}
