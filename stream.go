package pcloud

import (
	"net/url"
	"strconv"
)

type FileLinkOpts struct {
	ForceDownload bool
	ContentType   string
	MaxSpeed      int
}

func (c *Client) GetFileLink(fileID uint64) (*FileLink, error) {
	return c.GetFileLinkWithOpts(fileID, nil)
}

func (c *Client) GetFileLinkByPath(path string) (*FileLink, error) {
	return c.GetFileLinkByPathWithOpts(path, nil)
}

func (c *Client) GetFileLinkWithOpts(fileID uint64, opts *FileLinkOpts) (*FileLink, error) {
	params := url.Values{
		"fileid": {strconv.FormatUint(fileID, 10)},
	}
	applyLinkOpts(params, opts)

	var resp FileLink
	if err := c.do("getfilelink", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetFileLinkByPathWithOpts(path string, opts *FileLinkOpts) (*FileLink, error) {
	params := url.Values{
		"path": {path},
	}
	applyLinkOpts(params, opts)

	var resp FileLink
	if err := c.do("getfilelink", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetVideoLink(fileID uint64) (*FileLink, error) {
	params := url.Values{
		"fileid": {strconv.FormatUint(fileID, 10)},
	}

	var resp FileLink
	if err := c.do("getvideolink", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetAudioLink(fileID uint64) (*FileLink, error) {
	params := url.Values{
		"fileid": {strconv.FormatUint(fileID, 10)},
	}

	var resp FileLink
	if err := c.do("getaudiolink", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetHLSLink(fileID uint64) (*FileLink, error) {
	params := url.Values{
		"fileid": {strconv.FormatUint(fileID, 10)},
	}

	var resp FileLink
	if err := c.do("gethlslink", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp, nil
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
