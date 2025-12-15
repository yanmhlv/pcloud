package pcloud

import (
	"bytes"
	"context"
	"errors"
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

type ProgressFunc func(transferred, total int64)

type UploadOpts struct {
	NoPartial      bool
	RenameIfExists bool
	ModifiedTime   int64
	CreatedTime    int64
	OnProgress     ProgressFunc
}

type DownloadOpts struct {
	OnProgress ProgressFunc
}

type progressReader struct {
	reader     io.Reader
	total      int64
	read       int64
	onProgress ProgressFunc
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	if n > 0 && pr.onProgress != nil {
		pr.read += int64(n)
		pr.onProgress(pr.read, pr.total)
	}
	return n, err
}

func (pr *progressReader) Close() error {
	if c, ok := pr.reader.(io.Closer); ok {
		return c.Close()
	}
	return nil
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

func (c *Client) upload(ctx context.Context, params url.Values, filename string, content io.Reader, opts *UploadOpts) (*Metadata, error) {
	applyUploadOpts(params, opts)

	var contentSize int64 = -1
	if sizer, ok := content.(interface{ Len() int }); ok {
		contentSize = int64(sizer.Len())
	} else if seeker, ok := content.(io.Seeker); ok {
		pos, _ := seeker.Seek(0, io.SeekCurrent)
		end, err := seeker.Seek(0, io.SeekEnd)
		if err == nil {
			contentSize = end - pos
			_, _ = seeker.Seek(pos, io.SeekStart)
		}
	}

	readContent := content
	if opts != nil && opts.OnProgress != nil && contentSize > 0 {
		readContent = &progressReader{
			reader:     content,
			total:      contentSize,
			onProgress: opts.OnProgress,
		}
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, readContent); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	var resp uploadResponse
	if err := c.doPost(ctx, "uploadfile", params, &body, writer.FormDataContentType(), &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	if len(resp.Metadata) == 0 {
		return nil, errors.New("no metadata in response")
	}
	return &resp.Metadata[0], nil
}

func (c *Client) Upload(ctx context.Context, folderID uint64, filename string, content io.Reader, opts *UploadOpts) (*Metadata, error) {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
		"filename": {filename},
	}
	return c.upload(ctx, params, filename, content, opts)
}

func (c *Client) UploadByPath(ctx context.Context, path, filename string, content io.Reader, opts *UploadOpts) (*Metadata, error) {
	params := url.Values{
		"path":     {path},
		"filename": {filename},
	}
	return c.upload(ctx, params, filename, content, opts)
}

func (c *Client) Download(ctx context.Context, fileID uint64, opts *DownloadOpts) (io.ReadCloser, error) {
	link, err := c.GetFileLink(ctx, fileID)
	if err != nil {
		return nil, err
	}
	return c.downloadFromLink(ctx, link, opts)
}

func (c *Client) DownloadByPath(ctx context.Context, path string, opts *DownloadOpts) (io.ReadCloser, error) {
	link, err := c.GetFileLinkByPath(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.downloadFromLink(ctx, link, opts)
}

func (c *Client) downloadFromLink(ctx context.Context, link *FileLink, opts *DownloadOpts) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link.URL(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("download failed: %s", resp.Status)
	}

	if opts != nil && opts.OnProgress != nil {
		return &progressReader{
			reader:     resp.Body,
			total:      resp.ContentLength,
			onProgress: opts.OnProgress,
		}, nil
	}
	return resp.Body, nil
}

func (c *Client) Stat(ctx context.Context, fileID uint64) (*Metadata, error) {
	params := url.Values{
		"fileid": {strconv.FormatUint(fileID, 10)},
	}

	var resp fileResponse
	if err := c.do(ctx, "stat", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) StatByPath(ctx context.Context, path string) (*Metadata, error) {
	params := url.Values{
		"path": {path},
	}

	var resp fileResponse
	if err := c.do(ctx, "stat", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) DeleteFile(ctx context.Context, fileID uint64) error {
	params := url.Values{
		"fileid": {strconv.FormatUint(fileID, 10)},
	}

	var resp Error
	if err := c.do(ctx, "deletefile", params, &resp); err != nil {
		return err
	}
	return resp.Err()
}

func (c *Client) DeleteFileByPath(ctx context.Context, path string) error {
	params := url.Values{
		"path": {path},
	}

	var resp Error
	if err := c.do(ctx, "deletefile", params, &resp); err != nil {
		return err
	}
	return resp.Err()
}

func (c *Client) RenameFile(ctx context.Context, fileID uint64, newName string) (*Metadata, error) {
	params := url.Values{
		"fileid": {strconv.FormatUint(fileID, 10)},
		"toname": {newName},
	}

	var resp fileResponse
	if err := c.do(ctx, "renamefile", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) MoveFile(ctx context.Context, fileID, toFolderID uint64, name string) (*Metadata, error) {
	params := url.Values{
		"fileid":     {strconv.FormatUint(fileID, 10)},
		"tofolderid": {strconv.FormatUint(toFolderID, 10)},
		"toname":     {name},
	}

	var resp fileResponse
	if err := c.do(ctx, "renamefile", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) CopyFile(ctx context.Context, fileID, toFolderID uint64) (*Metadata, error) {
	params := url.Values{
		"fileid":     {strconv.FormatUint(fileID, 10)},
		"tofolderid": {strconv.FormatUint(toFolderID, 10)},
	}

	var resp fileResponse
	if err := c.do(ctx, "copyfile", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}
