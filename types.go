package pcloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Common pCloud API error codes returned in Error.Result.
const (
	ErrAuthRequired  = 1000
	ErrInvalidName   = 2001
	ErrAccessDenied  = 2003
	ErrAlreadyExists = 2004
	ErrNotFound      = 2005
)

type apiError interface {
	Err() error
}

type metadataResponse struct {
	Error
	Metadata Metadata `json:"metadata"`
}

// Time wraps time.Time to handle pCloud's RFC1123Z date format in JSON.
type Time struct {
	time.Time
}

// MarshalJSON encodes Time as RFC1123Z string, or "" for zero value.
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`""`), nil
	}
	return json.Marshal(t.Format(time.RFC1123Z))
}

// UnmarshalJSON decodes a pCloud RFC1123Z date string into Time.
func (t *Time) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	parsed, err := time.Parse(time.RFC1123Z, s)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

// Hash is a file content hash. pCloud may return it as a string or number.
type Hash string

// UnmarshalJSON decodes a pCloud hash value, which may be a JSON string or number.
func (h *Hash) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*h = Hash(s)
		return nil
	}
	var n uint64
	if err := json.Unmarshal(data, &n); err == nil {
		*h = Hash(strconv.FormatUint(n, 10))
		return nil
	}
	return fmt.Errorf("hash: cannot unmarshal %s", string(data))
}

// Error represents a pCloud API error response.
// Result holds the numeric code; zero means success.
type Error struct {
	Result  int    `json:"result"`
	Message string `json:"error"`
}

// Err returns nil when Result is 0, otherwise returns the Error itself.
func (e *Error) Err() error {
	if e.Result == 0 {
		return nil
	}
	return e
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("pcloud error %d", e.Result)
	}
	return fmt.Sprintf("pcloud error %d: %s", e.Result, e.Message)
}

// Metadata describes a pCloud file or folder.
type Metadata struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	Created     Time       `json:"created"`
	Modified    Time       `json:"modified"`
	IsFolder    bool       `json:"isfolder"`
	IsMine      bool       `json:"ismine"`
	IsShared    bool       `json:"isshared"`
	Icon        string     `json:"icon"`
	FileID      uint64     `json:"fileid,omitempty"`
	FolderID    uint64     `json:"folderid,omitempty"`
	ParentID    uint64     `json:"parentfolderid,omitempty"`
	Size        uint64     `json:"size,omitempty"`
	ContentType string     `json:"contenttype,omitempty"`
	Hash        Hash       `json:"hash,omitzero"`
	Category    int        `json:"category,omitempty"`
	Thumb       bool       `json:"thumb,omitempty"`
	Contents    []Metadata `json:"contents,omitzero"`
}

// Revision describes a historical version of a pCloud file.
type Revision struct {
	RevisionID uint64 `json:"revisionid"`
	Size       uint64 `json:"size"`
	Hash       Hash   `json:"hash"`
	Created    Time   `json:"created"`
}

// UserInfo holds account details returned by the userinfo endpoint.
type UserInfo struct {
	Error
	UserID         uint64 `json:"userid"`
	Email          string `json:"email"`
	EmailVerified  bool   `json:"emailverified"`
	Registered     Time   `json:"registered"`
	Language       string `json:"language"`
	Premium        bool   `json:"premium"`
	PremiumExpires Time   `json:"premiumexpires,omitzero"`
	Quota          uint64 `json:"quota"`
	UsedQuota      uint64 `json:"usedquota"`
}

// FileLink holds a temporary direct-download URL returned by streaming endpoints.
type FileLink struct {
	Error
	Path    string   `json:"path"`
	Expires string   `json:"expires"`
	Hosts   []string `json:"hosts"`
}

// URL returns the first available download URL for this FileLink.
func (f *FileLink) URL() string {
	if len(f.Hosts) == 0 {
		return ""
	}
	return "https://" + f.Hosts[0] + f.Path
}
