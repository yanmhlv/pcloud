package pcloud

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

func (c *pCloudClient) CreateFolder(path string, folderID int, name string) error {
	values := url.Values{
		"auth": {*c.Auth},
	}

	switch {
	case path != "":
		values["path"] = []string{path}
	case folderID >= 0 && name != "":
		values["folderid"] = []string{strconv.Itoa(folderID)}
		values["name"] = []string{name}
	default:
		return errors.New("bad params")
	}

	resp, err := c.Client.Get(urlBuilder("createfolder", values))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	result := struct {
		Result int    `json:"result"`
		Error  string `json:"error"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Result != 0 {
		return errors.New(result.Error)
	}

	return nil
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

func (c *pCloudClient) RenameFolder(folderID int, path string, topath string) error {
	values := url.Values{
		"auth":   {*c.Auth},
		"topath": {topath},
	}

	switch {
	case folderID >= 0:
		values["folderid"] = []string{strconv.Itoa(folderID)}
	case path != "":
		values["path"] = []string{path}
	default:
		return errors.New("bad params")
	}

	resp, err := c.Client.Get(urlBuilder("renamefolder", values))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	result := struct {
		Result int    `json:"result"`
		Error  string `json:"error"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Result != 0 {
		return errors.New(result.Error)
	}

	return nil
}

func (c *pCloudClient) DeleteFolder(path string, folderID int) error {
	values := url.Values{
		"auth": {*c.Auth},
	}

	switch {
	case path != "":
		values["path"] = []string{path}
	case folderID >= 0:
		values["folderid"] = []string{strconv.Itoa(folderID)}
	default:
		return errors.New("bad params")
	}

	resp, err := c.Client.Get(urlBuilder("deletefolder", values))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	result := struct {
		Result int    `json:"result"`
		Error  string `json:"error"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Result != 0 {
		return errors.New(result.Error)
	}

	return nil
}

func (c *pCloudClient) DeleteFolderRecursive(path string, folderID int) error {
	values := url.Values{
		"auth": {*c.Auth},
	}

	switch {
	case path != "":
		values["path"] = []string{path}
	case folderID >= 0:
		values["folderid"] = []string{strconv.Itoa(folderID)}
	default:
		return errors.New("bad params")
	}

	resp, err := c.Client.Get(urlBuilder("deletefolderrecursive", values))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	result := struct {
		Result int    `json:"result"`
		Error  string `json:"error"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Result != 0 {
		return errors.New(result.Error)
	}

	return nil
}
