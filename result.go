package pcloud

import (
	"encoding/json"
	"errors"
	"net/http"
)

// checkResult; returned error if request is failed or server returned error
func checkResult(resp *http.Response, err error) error {
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
