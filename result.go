package pcloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// convertToBuffer; convert http.Response.Body to bytes.Buffer
func convertToBuffer(resp *http.Response, err error) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	if err != nil {
		return buf, err
	}
	defer resp.Body.Close()

	_, err = buf.ReadFrom(resp.Body)
	return buf, err
}

// checkResult; returned error if request is failed or server returned error
func checkResult(resp *http.Response, err error) error {
	buf, err := convertToBuffer(resp, err)
	if err != nil {
		return err
	}

	result := struct {
		Result int    `json:"result"`
		Error  string `json:"error"`
	}{}

	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		return err
	}

	if result.Result != 0 {
		return errors.New(result.Error)
	}

	return nil
}
