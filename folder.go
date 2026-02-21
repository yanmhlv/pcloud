package pcloud

import (
	"context"
	"iter"
	"net/url"
	"strconv"
)

// ListFolderOpts controls optional parameters for listing folder contents.
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

func (c *Client) listFolder(ctx context.Context, params url.Values, opts *ListFolderOpts) (*Metadata, error) {
	applyListFolderOpts(params, opts)
	var resp metadataResponse
	if err := c.do(ctx, "listfolder", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

// ListFolder returns the contents of a folder identified by numeric ID.
func (c *Client) ListFolder(ctx context.Context, folderID uint64, opts *ListFolderOpts) (*Metadata, error) {
	return c.listFolder(ctx, url.Values{"folderid": {strconv.FormatUint(folderID, 10)}}, opts)
}

// ListFolderByPath returns the contents of a folder identified by path.
func (c *Client) ListFolderByPath(ctx context.Context, path string, opts *ListFolderOpts) (*Metadata, error) {
	return c.listFolder(ctx, url.Values{"path": {path}}, opts)
}

func (c *Client) statFolder(ctx context.Context, params url.Values) (*Metadata, error) {
	var resp metadataResponse
	if err := c.do(ctx, "stat", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

// StatFolder returns metadata for a folder identified by numeric ID.
func (c *Client) StatFolder(ctx context.Context, folderID uint64) (*Metadata, error) {
	return c.statFolder(ctx, url.Values{"folderid": {strconv.FormatUint(folderID, 10)}})
}

// StatFolderByPath returns metadata for a folder identified by path.
func (c *Client) StatFolderByPath(ctx context.Context, path string) (*Metadata, error) {
	return c.statFolder(ctx, url.Values{"path": {path}})
}

func (c *Client) createFolder(ctx context.Context, params url.Values) (*Metadata, error) {
	var resp metadataResponse
	if err := c.do(ctx, "createfolder", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

// CreateFolder creates a new folder inside parentID with the given name.
func (c *Client) CreateFolder(ctx context.Context, parentID uint64, name string) (*Metadata, error) {
	return c.createFolder(ctx, url.Values{
		"folderid": {strconv.FormatUint(parentID, 10)},
		"name":     {name},
	})
}

// CreateFolderByPath creates a folder at the given absolute path.
func (c *Client) CreateFolderByPath(ctx context.Context, path string) (*Metadata, error) {
	return c.createFolder(ctx, url.Values{"path": {path}})
}

// CreateFolderIfNotExists creates a folder only if it does not already exist.
func (c *Client) CreateFolderIfNotExists(ctx context.Context, parentID uint64, name string) (*Metadata, error) {
	params := url.Values{
		"folderid": {strconv.FormatUint(parentID, 10)},
		"name":     {name},
	}

	var resp metadataResponse
	if err := c.do(ctx, "createfolderifnotexists", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

// RenameFolder renames a folder in-place.
func (c *Client) RenameFolder(ctx context.Context, folderID uint64, newName string) (*Metadata, error) {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
		"toname":   {newName},
	}

	var resp metadataResponse
	if err := c.do(ctx, "renamefolder", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

// MoveFolder moves a folder to a different parent and optionally renames it.
func (c *Client) MoveFolder(ctx context.Context, folderID, toFolderID uint64, name string) (*Metadata, error) {
	params := url.Values{
		"folderid":   {strconv.FormatUint(folderID, 10)},
		"tofolderid": {strconv.FormatUint(toFolderID, 10)},
		"toname":     {name},
	}

	var resp metadataResponse
	if err := c.do(ctx, "renamefolder", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

// CopyFolder copies a folder into toFolderID.
func (c *Client) CopyFolder(ctx context.Context, folderID, toFolderID uint64) (*Metadata, error) {
	params := url.Values{
		"folderid":   {strconv.FormatUint(folderID, 10)},
		"tofolderid": {strconv.FormatUint(toFolderID, 10)},
	}

	var resp metadataResponse
	if err := c.do(ctx, "copyfolder", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

// DeleteFolder deletes an empty folder.
func (c *Client) DeleteFolder(ctx context.Context, folderID uint64) error {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
	}

	var resp Error
	return c.do(ctx, "deletefolder", params, &resp)
}

// DeleteFolderRecursive deletes a folder and all its contents.
func (c *Client) DeleteFolderRecursive(ctx context.Context, folderID uint64) error {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
	}

	var resp Error
	return c.do(ctx, "deletefolderrecursive", params, &resp)
}

func walkContents(ctx context.Context, contents []Metadata, yield func(Metadata, error) bool) {
	var walk func(items []Metadata) bool
	walk = func(items []Metadata) bool {
		if err := ctx.Err(); err != nil {
			yield(Metadata{}, err)
			return false
		}
		for _, item := range items {
			if err := ctx.Err(); err != nil {
				yield(Metadata{}, err)
				return false
			}
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

// Walk returns an iterator that yields every file and folder under folderID recursively.
// The caller may break early. Context cancellation stops iteration.
func (c *Client) Walk(ctx context.Context, folderID uint64) iter.Seq2[Metadata, error] {
	return func(yield func(Metadata, error) bool) {
		folder, err := c.ListFolder(ctx, folderID, &ListFolderOpts{Recursive: true})
		if err != nil {
			yield(Metadata{}, err)
			return
		}
		walkContents(ctx, folder.Contents, yield)
	}
}

// WalkByPath returns an iterator that yields every file and folder under path recursively.
func (c *Client) WalkByPath(ctx context.Context, path string) iter.Seq2[Metadata, error] {
	return func(yield func(Metadata, error) bool) {
		folder, err := c.ListFolderByPath(ctx, path, &ListFolderOpts{Recursive: true})
		if err != nil {
			yield(Metadata{}, err)
			return
		}
		walkContents(ctx, folder.Contents, yield)
	}
}
