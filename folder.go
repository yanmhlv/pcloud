package pcloud

import (
	"context"
	"iter"
	"net/url"
	"strconv"
)

type folderResponse struct {
	Error
	Metadata Metadata `json:"metadata"`
}

type ListFolderOpts struct {
	Recursive   bool
	ShowDeleted bool
	NoFiles     bool
	NoShares    bool
}

func applyListFolderOpts(params url.Values, opts *ListFolderOpts) {
	if opts == nil {
		return
	}
	if opts.Recursive {
		params.Set("recursive", "1")
	}
	if opts.ShowDeleted {
		params.Set("showdeleted", "1")
	}
	if opts.NoFiles {
		params.Set("nofiles", "1")
	}
	if opts.NoShares {
		params.Set("noshares", "1")
	}
}

func (c *Client) ListFolder(ctx context.Context, folderID uint64, opts *ListFolderOpts) (*Metadata, error) {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
	}
	applyListFolderOpts(params, opts)

	var resp folderResponse
	if err := c.do(ctx, "listfolder", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) ListFolderByPath(ctx context.Context, path string, opts *ListFolderOpts) (*Metadata, error) {
	params := url.Values{
		"path": {path},
	}
	applyListFolderOpts(params, opts)

	var resp folderResponse
	if err := c.do(ctx, "listfolder", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) StatFolder(ctx context.Context, folderID uint64) (*Metadata, error) {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
	}

	var resp folderResponse
	if err := c.do(ctx, "stat", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) CreateFolder(ctx context.Context, parentID uint64, name string) (*Metadata, error) {
	params := url.Values{
		"folderid": {strconv.FormatUint(parentID, 10)},
		"name":     {name},
	}

	var resp folderResponse
	if err := c.do(ctx, "createfolder", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) CreateFolderByPath(ctx context.Context, path string) (*Metadata, error) {
	params := url.Values{
		"path": {path},
	}

	var resp folderResponse
	if err := c.do(ctx, "createfolder", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) CreateFolderIfNotExists(ctx context.Context, parentID uint64, name string) (*Metadata, error) {
	params := url.Values{
		"folderid": {strconv.FormatUint(parentID, 10)},
		"name":     {name},
	}

	var resp folderResponse
	if err := c.do(ctx, "createfolderifnotexists", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) RenameFolder(ctx context.Context, folderID uint64, newName string) (*Metadata, error) {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
		"toname":   {newName},
	}

	var resp folderResponse
	if err := c.do(ctx, "renamefolder", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) MoveFolder(ctx context.Context, folderID, toFolderID uint64, name string) (*Metadata, error) {
	params := url.Values{
		"folderid":   {strconv.FormatUint(folderID, 10)},
		"tofolderid": {strconv.FormatUint(toFolderID, 10)},
		"toname":     {name},
	}

	var resp folderResponse
	if err := c.do(ctx, "renamefolder", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) CopyFolder(ctx context.Context, folderID, toFolderID uint64) (*Metadata, error) {
	params := url.Values{
		"folderid":   {strconv.FormatUint(folderID, 10)},
		"tofolderid": {strconv.FormatUint(toFolderID, 10)},
	}

	var resp folderResponse
	if err := c.do(ctx, "copyfolder", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) DeleteFolder(ctx context.Context, folderID uint64) error {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
	}

	var resp Error
	if err := c.do(ctx, "deletefolder", params, &resp); err != nil {
		return err
	}
	return resp.Err()
}

func (c *Client) DeleteFolderRecursive(ctx context.Context, folderID uint64) error {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
	}

	var resp Error
	if err := c.do(ctx, "deletefolderrecursive", params, &resp); err != nil {
		return err
	}
	return resp.Err()
}

func walkContents(contents []Metadata, yield func(Metadata, error) bool) {
	var walk func(items []Metadata) bool
	walk = func(items []Metadata) bool {
		for _, item := range items {
			if !yield(item, nil) {
				return false
			}
			if item.IsFolder && len(item.Contents) > 0 {
				if !walk(item.Contents) {
					return false
				}
			}
		}
		return true
	}
	walk(contents)
}

func (c *Client) Walk(ctx context.Context, folderID uint64) iter.Seq2[Metadata, error] {
	return func(yield func(Metadata, error) bool) {
		folder, err := c.ListFolder(ctx, folderID, &ListFolderOpts{Recursive: true})
		if err != nil {
			yield(Metadata{}, err)
			return
		}
		walkContents(folder.Contents, yield)
	}
}

func (c *Client) WalkByPath(ctx context.Context, path string) iter.Seq2[Metadata, error] {
	return func(yield func(Metadata, error) bool) {
		folder, err := c.ListFolderByPath(ctx, path, &ListFolderOpts{Recursive: true})
		if err != nil {
			yield(Metadata{}, err)
			return
		}
		walkContents(folder.Contents, yield)
	}
}
