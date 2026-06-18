package pcloud

import (
	"context"
	"fmt"
	"net/url"
)

const (
	minThumbSize = 1
	maxThumbSize = 2048
)

type ThumbOpts struct {
	Crop bool
	Type string
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
	outOfRange := width < minThumbSize || width > maxThumbSize || height < minThumbSize || height > maxThumbSize
	if outOfRange {
		return fmt.Errorf("thumbnail size %dx%d is out of allowed range %d–%d", width, height, minThumbSize, maxThumbSize)
	}
	return nil
}

func (c *Client) getThumbnail(ctx context.Context, params url.Values, width, height int, opts *ThumbOpts) (*FileLink, error) {
	if err := validateThumbSize(width, height); err != nil {
		return nil, err
	}
	params.Set("size", fmt.Sprintf("%dx%d", width, height))
	applyThumbOpts(params, opts)

	var resp fileLinkResponse
	if err := c.doGet(ctx, "getthumb", params, &resp); err != nil {
		return nil, err
	}
	return &resp.FileLink, nil
}

func (c *Client) GetThumbnail(ctx context.Context, fileID uint64, width, height int, opts *ThumbOpts) (*FileLink, error) {
	fl, err := c.getThumbnail(ctx, url.Values{paramFileID: {formatUint(fileID)}}, width, height, opts)
	if err != nil {
		return nil, fmt.Errorf("get thumbnail %d: %w", fileID, err)
	}
	return fl, nil
}

func (c *Client) GetThumbnailByPath(ctx context.Context, path string, width, height int, opts *ThumbOpts) (*FileLink, error) {
	fl, err := c.getThumbnail(ctx, url.Values{paramPath: {path}}, width, height, opts)
	if err != nil {
		return nil, fmt.Errorf("get thumbnail %s: %w", path, err)
	}
	return fl, nil
}
