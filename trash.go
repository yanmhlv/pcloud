package pcloud

import (
	"context"
	"net/url"
	"strconv"
)

// TrashItem describes a file or folder in the pCloud trash.
type TrashItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	IsFolder bool   `json:"isfolder"`
	FileID   uint64 `json:"fileid,omitempty"`
	FolderID uint64 `json:"folderid,omitempty"`
	Size     uint64 `json:"size,omitempty"`
	Deleted  Time   `json:"deleted"`
}

type trashListResponse struct {
	Error
	Items []TrashItem `json:"items"`
}

// ListTrash returns all files and folders currently in the trash.
func (c *Client) ListTrash(ctx context.Context) ([]TrashItem, error) {
	var resp trashListResponse
	if err := c.do(ctx, "trash_list", url.Values{}, &resp); err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func (c *Client) restoreFromTrash(ctx context.Context, params url.Values) (*Metadata, error) {
	var resp metadataResponse
	if err := c.do(ctx, "trash_restoretofolder", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

// RestoreFromTrash restores a trashed file to its original location by file ID.
func (c *Client) RestoreFromTrash(ctx context.Context, fileID uint64) (*Metadata, error) {
	return c.restoreFromTrash(ctx, url.Values{"fileid": {strconv.FormatUint(fileID, 10)}})
}

// RestoreFromTrashByPath restores a trashed item by its path.
func (c *Client) RestoreFromTrashByPath(ctx context.Context, path string) (*Metadata, error) {
	return c.restoreFromTrash(ctx, url.Values{"path": {path}})
}

// EmptyTrash permanently deletes all items in the trash.
func (c *Client) EmptyTrash(ctx context.Context) error {
	var resp Error
	return c.do(ctx, "trash_clear", url.Values{}, &resp)
}
