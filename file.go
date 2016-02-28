package pcloud

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

// uploadprogress
// downloadfile
// checksumfile

func (c *pCloudClient) UploadFile(reader io.Reader, path string, folderID int, filename string, noPartial int, progressHash string, renameIfExists int) error {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	values := url.Values{
		"auth": {*c.Auth},
	}

	if noPartial > 0 {
		values["nopartial"] = []string{strconv.Itoa(noPartial)}
	}
	if progressHash != "" {
		values["progresshash"] = []string{progressHash}
	}
	if renameIfExists > 0 {
		values["renameifexists"] = []string{strconv.Itoa(renameIfExists)}
	}

	switch {
	case path != "":
		values["path"] = []string{path}
	case folderID >= 0:
		values["folderid"] = []string{strconv.Itoa(folderID)}
	default:
		return errors.New("bad params")
	}

	if filename == "" {
		return errors.New("bad params")
	}

	fw, err := w.CreateFormFile(filename, filename)
	if err != nil {
		return err
	}
	if _, err := io.Copy(fw, reader); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", urlBuilder("uploadfile", values), &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	return checkResult(c.Client.Do(req))
}

func (c *pCloudClient) CopyFile(fileID int, path string, toFolderID int, toName string, toPath string) error {
	values := url.Values{
		"auth": {*c.Auth},
	}

	switch {
	case fileID > 0:
		values["fileid"] = []string{strconv.Itoa(fileID)}
	case path != "":
		values["path"] = []string{path}
	default:
		return errors.New("bad params")
	}

	switch {
	case toFolderID > 0 && toName != "":
		values["tofolderid"] = []string{strconv.Itoa(toFolderID)}
		values["toname"] = []string{toName}
	case toPath != "":
		values["topath"] = []string{toPath}
	default:
		return errors.New("bad params")
	}

	return checkResult(c.Client.Get(urlBuilder("copyfile", values)))
}

func (c *pCloudClient) DeleteFile(fileID int, path string) error {
	values := url.Values{
		"auth": {*c.Auth},
	}

	switch {
	case fileID > 0:
		values["fileid"] = []string{strconv.Itoa(fileID)}
	case path != "":
		values["path"] = []string{path}
	default:
		return errors.New("bad params")
	}

	return checkResult(c.Client.Get(urlBuilder("deletefile", values)))
}

func (c *pCloudClient) RenameFile(fileID int, path string, toPath string, toFolderID int, toName string) error {
	values := url.Values{
		"auth": {*c.Auth},
	}

	switch {
	case fileID > 0:
		values["fileid"] = []string{strconv.Itoa(fileID)}
	case path != "":
		values["path"] = []string{path}
	default:
		return errors.New("bad params")
	}

	switch {
	case toPath != "":
		values["topath"] = []string{toPath}
	case toFolderID > 0 && toName != "":
		values["toname"] = []string{toName}
		values["tofolderid"] = []string{strconv.Itoa(toFolderID)}
	default:
		return errors.New("bad params")
	}

	return checkResult(c.Client.Get(urlBuilder("renamefile", values)))
}
