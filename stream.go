package pcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// getvideolink
// getvideolinks
// getaudiolink
// gethlslink
// gettextfile

// GetFileLink; https://docs.pcloud.com/methods/streaming/getfilelink.html
func (c *pCloudClient) GetFileLink(fileID int, path string, forceDownload int, contentType string, maxSpeed int, skipFilename int) ([]string, error) {
	var links []string

	values := url.Values{}

	switch {
	case fileID > 0:
		values.Add("fileid", strconv.Itoa(fileID))
	case path != "":
		values.Add("path", path)
	default:
		return links, errors.New("bad params")
	}

	if forceDownload > 0 {
		values.Add("forcedownload", strconv.Itoa(forceDownload))
	}
	if contentType != "" {
		values.Add("contenttype", contentType)
	}
	if maxSpeed > 0 {
		values.Add("maxspeed", strconv.Itoa(maxSpeed))
	}
	if skipFilename > 0 {
		values.Add("skipfilename", strconv.Itoa(skipFilename))
	}

	resp, err := c.Client.Get(urlBuilder("getfilelink", values))
	if err != nil {
		return links, err
	}
	defer resp.Body.Close()

	result := struct {
		Result int      `json:"result"`
		Error  string   `json:"error"`
		Path   string   `json:"path"`
		Hosts  []string `json:"hosts"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return links, err
	}
	if result.Result > 0 {
		return links, errors.New(result.Error)
	}

	for _, host := range result.Hosts {
		links = append(links, fmt.Sprintf("https://%s%s", host, result.Path))
	}
	return links, nil
}
