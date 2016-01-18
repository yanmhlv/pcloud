package pcloud

import (
	"fmt"
	"net/http"
    "fmt"
)

func (c *Client) ListFolder(path string, folderid int, recursive int, showdeleted int, nofiles int, noshares int) (map[int]string, error) {
	// if path == "" and folderid < 0 {
	//     return nil, errors.New("Bad path or folderid!")
	// }

	urlStr := "/listfolder?%s=%s"
	switch {
	case path != "":
		urlStr = fmt.Sprintf(urlStr, "path", path)
	case folderid >= 0:
		urlStr = fmt.Sprintf(urlStr, "folderid", folderid)
    default:
        return nil, errors.New("bad path or folderid")
	}

    for k, v := map[string]int{
        "recursive": recursive,
        "showdeleted": showdeleted,
        "nofiles": nofiles,
        "noshares": noshares,
    } {
        if v > 0 {
            urlStr += fmt.Sprintf("&%s=%d", k, v)
        }

    }

	req, err := http.NewRequest("GET", urlStr, body)
	if err != nil {
		return nil, err
	}

	c.prepareRequest(req)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	c.prepareResponse(resp, v)
}
