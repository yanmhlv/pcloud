package pcloud

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type uploadResponse struct {
	Error
	FileIDs  []uint64   `json:"fileids"`
	Metadata []Metadata `json:"metadata"`
}

// ProgressFunc is called during upload or download with cumulative bytes transferred and total size.
// Total may be -1 if the size is unknown.
type ProgressFunc func(transferred, total int64)

// UploadOpts controls optional parameters for file uploads.
type UploadOpts struct {
	NoPartial      bool
	RenameIfExists bool
	ModifiedTime   time.Time
	CreatedTime    time.Time
	OnProgress     ProgressFunc
}

// DownloadOpts controls optional parameters for file downloads.
type DownloadOpts struct {
	OnProgress ProgressFunc
}

type progressReader struct {
	rc         io.ReadCloser
	total      int64
	read       int64
	onProgress ProgressFunc
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.rc.Read(p)
	if n > 0 && pr.onProgress != nil {
		pr.read += int64(n)
		pr.onProgress(pr.read, pr.total)
	}
	return n, err
}

func (pr *progressReader) Close() error {
	return pr.rc.Close()
}

func getContentSize(r io.Reader) (int64, error) {
	seeker, ok := r.(io.Seeker)
	if !ok {
		return -1, nil
	}
	pos, err := seeker.Seek(0, io.SeekCurrent)
	if err != nil {
		return -1, fmt.Errorf("seek current: %w", err)
	}
	end, err := seeker.Seek(0, io.SeekEnd)
	if err != nil {
		return -1, fmt.Errorf("seek end: %w", err)
	}
	if _, err := seeker.Seek(pos, io.SeekStart); err != nil {
		return -1, fmt.Errorf("seek start: %w", err)
	}
	return end - pos, nil
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
	if !opts.ModifiedTime.IsZero() {
		params.Set("mtime", strconv.FormatInt(opts.ModifiedTime.Unix(), 10))
	}
	if !opts.CreatedTime.IsZero() {
		params.Set("ctime", strconv.FormatInt(opts.CreatedTime.Unix(), 10))
	}
}

func (c *Client) upload(ctx context.Context, params url.Values, filename string, content io.Reader, opts *UploadOpts) (*Metadata, error) {
	applyUploadOpts(params, opts)

	var contentSize int64 = -1
	if sizer, ok := content.(interface{ Len() int }); ok {
		contentSize = int64(sizer.Len())
	}

	if contentSize < 0 {
		size, err := getContentSize(content)
		if err != nil {
			return nil, err
		}
		contentSize = size
	}

	readContent := content
	if opts != nil && opts.OnProgress != nil && contentSize > 0 {
		readContent = &progressReader{
			rc:         io.NopCloser(content),
			total:      contentSize,
			onProgress: opts.OnProgress,
		}
	}

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		err := func() error {
			part, err := writer.CreateFormFile("file", filename)
			if err != nil {
				return err
			}
			if _, err := io.Copy(part, readContent); err != nil {
				return err
			}
			return writer.Close()
		}()
		_ = pw.CloseWithError(err)
		errCh <- err
	}()

	var resp uploadResponse
	if err := c.doPost(ctx, "uploadfile", params, pr, writer.FormDataContentType(), &resp); err != nil {
		_ = pr.CloseWithError(err)
		return nil, err
	}
	if err := <-errCh; err != nil {
		return nil, err
	}
	if len(resp.Metadata) == 0 {
		return nil, errors.New("no metadata in response")
	}
	return &resp.Metadata[0], nil
}

// Upload uploads content to folderID as filename.
func (c *Client) Upload(ctx context.Context, folderID uint64, filename string, content io.Reader, opts *UploadOpts) (*Metadata, error) {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
		"filename": {filename},
	}
	m, err := c.upload(ctx, params, filename, content, opts)
	if err != nil {
		return nil, fmt.Errorf("upload to folder %d: %w", folderID, err)
	}
	return m, nil
}

// UploadByPath uploads content to the folder at path as filename.
func (c *Client) UploadByPath(ctx context.Context, path, filename string, content io.Reader, opts *UploadOpts) (*Metadata, error) {
	params := url.Values{
		"path":     {path},
		"filename": {filename},
	}
	m, err := c.upload(ctx, params, filename, content, opts)
	if err != nil {
		return nil, fmt.Errorf("upload to %s: %w", path, err)
	}
	return m, nil
}

// Download downloads a file by ID and returns the response body.
// The caller must close the returned ReadCloser.
func (c *Client) Download(ctx context.Context, fileID uint64, opts *DownloadOpts) (io.ReadCloser, error) {
	link, err := c.GetFileLink(ctx, fileID, nil)
	if err != nil {
		return nil, fmt.Errorf("download file %d: %w", fileID, err)
	}
	rc, err := c.downloadFromLink(ctx, link, opts)
	if err != nil {
		return nil, fmt.Errorf("download file %d: %w", fileID, err)
	}
	return rc, nil
}

// DownloadByPath downloads a file by path and returns the response body.
// The caller must close the returned ReadCloser.
func (c *Client) DownloadByPath(ctx context.Context, path string, opts *DownloadOpts) (io.ReadCloser, error) {
	link, err := c.GetFileLinkByPath(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", path, err)
	}
	rc, err := c.downloadFromLink(ctx, link, opts)
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", path, err)
	}
	return rc, nil
}

func (c *Client) downloadFromLink(ctx context.Context, link *FileLink, opts *DownloadOpts) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link.URL(), nil)
	if err != nil {
		return nil, err
	}

	c.mu.RLock()
	httpCl := c.httpClient
	c.mu.RUnlock()

	resp, err := httpCl.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("download failed: %s", resp.Status)
	}

	if opts != nil && opts.OnProgress != nil {
		return &progressReader{
			rc:         resp.Body,
			total:      resp.ContentLength,
			onProgress: opts.OnProgress,
		}, nil
	}
	return resp.Body, nil
}

func (c *Client) stat(ctx context.Context, params url.Values) (*Metadata, error) {
	var resp metadataResponse
	if err := c.do(ctx, "stat", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

// Stat returns metadata for a file identified by numeric ID.
func (c *Client) Stat(ctx context.Context, fileID uint64) (*Metadata, error) {
	m, err := c.stat(ctx, url.Values{"fileid": {strconv.FormatUint(fileID, 10)}})
	if err != nil {
		return nil, fmt.Errorf("stat file %d: %w", fileID, err)
	}
	return m, nil
}

// StatByPath returns metadata for a file identified by path.
func (c *Client) StatByPath(ctx context.Context, path string) (*Metadata, error) {
	m, err := c.stat(ctx, url.Values{"path": {path}})
	if err != nil {
		return nil, fmt.Errorf("stat %s: %w", path, err)
	}
	return m, nil
}

func (c *Client) deleteFile(ctx context.Context, params url.Values) error {
	var resp Error
	return c.do(ctx, "deletefile", params, &resp)
}

// DeleteFile deletes a file by numeric ID.
func (c *Client) DeleteFile(ctx context.Context, fileID uint64) error {
	if err := c.deleteFile(ctx, url.Values{"fileid": {strconv.FormatUint(fileID, 10)}}); err != nil {
		return fmt.Errorf("delete file %d: %w", fileID, err)
	}
	return nil
}

// DeleteFileByPath deletes a file by path.
func (c *Client) DeleteFileByPath(ctx context.Context, path string) error {
	if err := c.deleteFile(ctx, url.Values{"path": {path}}); err != nil {
		return fmt.Errorf("delete %s: %w", path, err)
	}
	return nil
}

// RenameFile renames a file in-place.
func (c *Client) RenameFile(ctx context.Context, fileID uint64, newName string) (*Metadata, error) {
	params := url.Values{
		"fileid": {strconv.FormatUint(fileID, 10)},
		"toname": {newName},
	}

	var resp metadataResponse
	if err := c.do(ctx, "renamefile", params, &resp); err != nil {
		return nil, fmt.Errorf("rename file %d: %w", fileID, err)
	}
	return &resp.Metadata, nil
}

// MoveFile moves a file to a different folder and optionally renames it.
func (c *Client) MoveFile(ctx context.Context, fileID, toFolderID uint64, name string) (*Metadata, error) {
	params := url.Values{
		"fileid":     {strconv.FormatUint(fileID, 10)},
		"tofolderid": {strconv.FormatUint(toFolderID, 10)},
		"toname":     {name},
	}

	var resp metadataResponse
	if err := c.do(ctx, "renamefile", params, &resp); err != nil {
		return nil, fmt.Errorf("move file %d: %w", fileID, err)
	}
	return &resp.Metadata, nil
}

// CopyFile copies a file into toFolderID.
func (c *Client) CopyFile(ctx context.Context, fileID, toFolderID uint64) (*Metadata, error) {
	params := url.Values{
		"fileid":     {strconv.FormatUint(fileID, 10)},
		"tofolderid": {strconv.FormatUint(toFolderID, 10)},
	}

	var resp metadataResponse
	if err := c.do(ctx, "copyfile", params, &resp); err != nil {
		return nil, fmt.Errorf("copy file %d: %w", fileID, err)
	}
	return &resp.Metadata, nil
}
