package pcloud

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

type fileResponse struct {
	Error
	Metadata Metadata `json:"metadata"`
}

type uploadResponse struct {
	Error
	FileIDs  []uint64   `json:"fileids"`
	Metadata []Metadata `json:"metadata"`
}

type UploadOpts struct {
	NoPartial      bool
	RenameIfExists bool
	ModifiedTime   int64
	CreatedTime    int64
}

func applyUploadOpts(params url.Values, opts *UploadOpts) {
	if opts == nil {
		return
	}
	if opts.NoPartial {
		params.Set("nopartial", "1")
	}
	if opts.RenameIfExists {
		params.Set("renameifexists", "1")
	}
	if opts.ModifiedTime > 0 {
		params.Set("mtime", strconv.FormatInt(opts.ModifiedTime, 10))
	}
	if opts.CreatedTime > 0 {
		params.Set("ctime", strconv.FormatInt(opts.CreatedTime, 10))
	}
}

func (c *Client) upload(params url.Values, filename string, content io.Reader, opts *UploadOpts) (*Metadata, error) {
	applyUploadOpts(params, opts)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, content); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	var resp uploadResponse
	if err := c.doPost("uploadfile", params, &body, writer.FormDataContentType(), &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	if len(resp.Metadata) == 0 {
		return nil, fmt.Errorf("no metadata in response")
	}
	return &resp.Metadata[0], nil
}

func (c *Client) Upload(folderID uint64, filename string, content io.Reader, opts *UploadOpts) (*Metadata, error) {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
		"filename": {filename},
	}
	return c.upload(params, filename, content, opts)
}

func (c *Client) UploadByPath(path, filename string, content io.Reader, opts *UploadOpts) (*Metadata, error) {
	params := url.Values{
		"path":     {path},
		"filename": {filename},
	}
	return c.upload(params, filename, content, opts)
}

func (c *Client) Download(fileID uint64) (io.ReadCloser, error) {
	link, err := c.GetFileLink(fileID)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Get(link.URL())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("download failed: %s", resp.Status)
	}
	return resp.Body, nil
}

func (c *Client) DownloadByPath(path string) (io.ReadCloser, error) {
	link, err := c.GetFileLinkByPath(path)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Get(link.URL())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("download failed: %s", resp.Status)
	}
	return resp.Body, nil
}

func (c *Client) Stat(fileID uint64) (*Metadata, error) {
	params := url.Values{
		"fileid": {strconv.FormatUint(fileID, 10)},
	}

	var resp fileResponse
	if err := c.do("stat", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) StatByPath(path string) (*Metadata, error) {
	params := url.Values{
		"path": {path},
	}

	var resp fileResponse
	if err := c.do("stat", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) DeleteFile(fileID uint64) error {
	params := url.Values{
		"fileid": {strconv.FormatUint(fileID, 10)},
	}

	var resp Error
	if err := c.do("deletefile", params, &resp); err != nil {
		return err
	}
	return resp.Err()
}

func (c *Client) DeleteFileByPath(path string) error {
	params := url.Values{
		"path": {path},
	}

	var resp Error
	if err := c.do("deletefile", params, &resp); err != nil {
		return err
	}
	return resp.Err()
}

func (c *Client) RenameFile(fileID uint64, newName string) (*Metadata, error) {
	params := url.Values{
		"fileid": {strconv.FormatUint(fileID, 10)},
		"toname": {newName},
	}

	var resp fileResponse
	if err := c.do("renamefile", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) MoveFile(fileID, toFolderID uint64) (*Metadata, error) {
	params := url.Values{
		"fileid":     {strconv.FormatUint(fileID, 10)},
		"tofolderid": {strconv.FormatUint(toFolderID, 10)},
	}

	var resp fileResponse
	if err := c.do("renamefile", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) CopyFile(fileID, toFolderID uint64) (*Metadata, error) {
	params := url.Values{
		"fileid":     {strconv.FormatUint(fileID, 10)},
		"tofolderid": {strconv.FormatUint(toFolderID, 10)},
	}

	var resp fileResponse
	if err := c.do("copyfile", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}
