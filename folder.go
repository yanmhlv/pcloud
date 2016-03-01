package pcloud

import (
	"errors"
	"net/url"
	"strconv"
)

// CreateFolder; https://docs.pcloud.com/methods/folder/createfolder.html
func (c *pCloudClient) CreateFolder(path string, folderID int, name string) error {
	values := url.Values{
		"auth": {*c.Auth},
	}

	switch {
	case path != "":
		values.Add("path", path)
	case folderID >= 0 && name != "":
		values.Add("folderid", strconv.Itoa(folderID))
		values.Add("name", name)
	default:
		return errors.New("bad params")
	}

	return checkResult(c.Client.Get(urlBuilder("createfolder", values)))
}

// func (c *pCloudClient) ListFolder() error {
//  u := (&url.URL{
//      Scheme:   apiScheme,
//      Host:     apiHost,
//      Path:     "listfolder",
//      RawQuery: url.Values{}.Encode(),
//  }).String()
//  return nil
// }

// RenameFolder; https://docs.pcloud.com/methods/folder/renamefolder.html
func (c *pCloudClient) RenameFolder(folderID int, path string, topath string) error {
	values := url.Values{
		"auth":   {*c.Auth},
		"topath": {topath},
	}

	switch {
	case folderID >= 0:
		values.Add("folderid", strconv.Itoa(folderID))
	case path != "":
		values.Add("path", path)
	default:
		return errors.New("bad params")
	}

	return checkResult(c.Client.Get(urlBuilder("renamefolder", values)))
}

// DeleteFolder; https://docs.pcloud.com/methods/folder/deletefolder.html
func (c *pCloudClient) DeleteFolder(path string, folderID int) error {
	values := url.Values{
		"auth": {*c.Auth},
	}

	switch {
	case path != "":
		values.Add("path", path)
	case folderID >= 0:
		values.Add("folderid", strconv.Itoa(folderID))
	default:
		return errors.New("bad params")
	}

	return checkResult(c.Client.Get(urlBuilder("deletefolder", values)))
}

// DeleteFolderRecursive; https://docs.pcloud.com/methods/folder/deletefolderrecursive.html
func (c *pCloudClient) DeleteFolderRecursive(path string, folderID int) error {
	values := url.Values{
		"auth": {*c.Auth},
	}

	switch {
	case path != "":
		values.Add("path", path)
	case folderID >= 0:
		values.Add("folderid", strconv.Itoa(folderID))
	default:
		return errors.New("bad params")
	}

	return checkResult(c.Client.Get(urlBuilder("deletefolderrecursive", values)))
}
