package pcloud

import (
	"context"
	"fmt"
	"iter"
	"net/url"
	"strconv"
)

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

func (c *Client) ListFolder(ctx context.Context, folderID uint64, opts *ListFolderOpts) (*Metadata, error) {
	m, err := c.listFolder(ctx, url.Values{"folderid": {strconv.FormatUint(folderID, 10)}}, opts)
	if err != nil {
		return nil, fmt.Errorf("list folder %d: %w", folderID, err)
	}
	return m, nil
}

func (c *Client) ListFolderByPath(ctx context.Context, path string, opts *ListFolderOpts) (*Metadata, error) {
	m, err := c.listFolder(ctx, url.Values{"path": {path}}, opts)
	if err != nil {
		return nil, fmt.Errorf("list folder %s: %w", path, err)
	}
	return m, nil
}

func (c *Client) createFolder(ctx context.Context, params url.Values) (*Metadata, error) {
	var resp metadataResponse
	if err := c.do(ctx, "createfolder", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) CreateFolder(ctx context.Context, parentID uint64, name string) (*Metadata, error) {
	m, err := c.createFolder(ctx, url.Values{
		"folderid": {strconv.FormatUint(parentID, 10)},
		"name":     {name},
	})
	if err != nil {
		return nil, fmt.Errorf("create folder %d/%s: %w", parentID, name, err)
	}
	return m, nil
}

func (c *Client) CreateFolderByPath(ctx context.Context, path string) (*Metadata, error) {
	m, err := c.createFolder(ctx, url.Values{"path": {path}})
	if err != nil {
		return nil, fmt.Errorf("create folder %s: %w", path, err)
	}
	return m, nil
}

func (c *Client) CreateFolderIfNotExists(ctx context.Context, parentID uint64, name string) (*Metadata, error) {
	params := url.Values{
		"folderid": {strconv.FormatUint(parentID, 10)},
		"name":     {name},
	}

	var resp metadataResponse
	if err := c.do(ctx, "createfolderifnotexists", params, &resp); err != nil {
		return nil, fmt.Errorf("create folder if not exists %d/%s: %w", parentID, name, err)
	}
	return &resp.Metadata, nil
}

func (c *Client) RenameFolder(ctx context.Context, folderID uint64, newName string) (*Metadata, error) {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
		"toname":   {newName},
	}

	var resp metadataResponse
	if err := c.do(ctx, "renamefolder", params, &resp); err != nil {
		return nil, fmt.Errorf("rename folder %d: %w", folderID, err)
	}
	return &resp.Metadata, nil
}

func (c *Client) MoveFolder(ctx context.Context, folderID, toFolderID uint64, name string) (*Metadata, error) {
	params := url.Values{
		"folderid":   {strconv.FormatUint(folderID, 10)},
		"tofolderid": {strconv.FormatUint(toFolderID, 10)},
		"toname":     {name},
	}

	var resp metadataResponse
	if err := c.do(ctx, "renamefolder", params, &resp); err != nil {
		return nil, fmt.Errorf("move folder %d: %w", folderID, err)
	}
	return &resp.Metadata, nil
}

func (c *Client) CopyFolder(ctx context.Context, folderID, toFolderID uint64) (*Metadata, error) {
	params := url.Values{
		"folderid":   {strconv.FormatUint(folderID, 10)},
		"tofolderid": {strconv.FormatUint(toFolderID, 10)},
	}

	var resp metadataResponse
	if err := c.do(ctx, "copyfolder", params, &resp); err != nil {
		return nil, fmt.Errorf("copy folder %d: %w", folderID, err)
	}
	return &resp.Metadata, nil
}

func (c *Client) DeleteFolder(ctx context.Context, folderID uint64) error {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
	}

	var resp Error
	if err := c.do(ctx, "deletefolder", params, &resp); err != nil {
		return fmt.Errorf("delete folder %d: %w", folderID, err)
	}
	return nil
}

func (c *Client) DeleteFolderRecursive(ctx context.Context, folderID uint64) error {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
	}

	var resp Error
	if err := c.do(ctx, "deletefolderrecursive", params, &resp); err != nil {
		return fmt.Errorf("delete folder recursive %d: %w", folderID, err)
	}
	return nil
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
