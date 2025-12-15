package pcloud

import (
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

func (c *Client) ListFolder(folderID uint64, opts *ListFolderOpts) (*Metadata, error) {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
	}
	if opts != nil {
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

	var resp folderResponse
	if err := c.do("listfolder", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) ListFolderByPath(path string, opts *ListFolderOpts) (*Metadata, error) {
	params := url.Values{
		"path": {path},
	}
	if opts != nil {
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

	var resp folderResponse
	if err := c.do("listfolder", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) CreateFolder(parentID uint64, name string) (*Metadata, error) {
	params := url.Values{
		"folderid": {strconv.FormatUint(parentID, 10)},
		"name":     {name},
	}

	var resp folderResponse
	if err := c.do("createfolder", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) CreateFolderByPath(path string) (*Metadata, error) {
	params := url.Values{
		"path": {path},
	}

	var resp folderResponse
	if err := c.do("createfolder", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) CreateFolderIfNotExists(parentID uint64, name string) (*Metadata, error) {
	params := url.Values{
		"folderid": {strconv.FormatUint(parentID, 10)},
		"name":     {name},
	}

	var resp folderResponse
	if err := c.do("createfolderifnotexists", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) RenameFolder(folderID uint64, newName string) (*Metadata, error) {
	params := url.Values{
		"folderid":   {strconv.FormatUint(folderID, 10)},
		"toname":     {newName},
	}

	var resp folderResponse
	if err := c.do("renamefolder", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) MoveFolder(folderID, toFolderID uint64) (*Metadata, error) {
	params := url.Values{
		"folderid":   {strconv.FormatUint(folderID, 10)},
		"tofolderid": {strconv.FormatUint(toFolderID, 10)},
	}

	var resp folderResponse
	if err := c.do("renamefolder", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) CopyFolder(folderID, toFolderID uint64) (*Metadata, error) {
	params := url.Values{
		"folderid":   {strconv.FormatUint(folderID, 10)},
		"tofolderid": {strconv.FormatUint(toFolderID, 10)},
	}

	var resp folderResponse
	if err := c.do("copyfolder", params, &resp); err != nil {
		return nil, err
	}
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return &resp.Metadata, nil
}

func (c *Client) DeleteFolder(folderID uint64) error {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
	}

	var resp Error
	if err := c.do("deletefolder", params, &resp); err != nil {
		return err
	}
	return resp.Err()
}

func (c *Client) DeleteFolderRecursive(folderID uint64) error {
	params := url.Values{
		"folderid": {strconv.FormatUint(folderID, 10)},
	}

	var resp Error
	if err := c.do("deletefolderrecursive", params, &resp); err != nil {
		return err
	}
	return resp.Err()
}
