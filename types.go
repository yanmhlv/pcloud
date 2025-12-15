package pcloud

import "time"

type Error struct {
	Result  int    `json:"result"`
	Message string `json:"error"`
}

func (e *Error) Err() error {
	if e.Result == 0 {
		return nil
	}
	return e
}

func (e *Error) Error() string {
	return e.Message
}

type Metadata struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	Created     time.Time  `json:"created"`
	Modified    time.Time  `json:"modified"`
	IsFolder    bool       `json:"isfolder"`
	IsMine      bool       `json:"ismine"`
	IsShared    bool       `json:"isshared"`
	Icon        string     `json:"icon"`
	FileID      uint64     `json:"fileid,omitempty"`
	FolderID    uint64     `json:"folderid,omitempty"`
	ParentID    uint64     `json:"parentfolderid,omitempty"`
	Size        uint64     `json:"size,omitempty"`
	ContentType string     `json:"contenttype,omitempty"`
	Hash        string     `json:"hash,omitempty"`
	Category    int        `json:"category,omitempty"`
	Thumb       bool       `json:"thumb,omitempty"`
	Contents    []Metadata `json:"contents,omitempty"`
}

type Revision struct {
	RevisionID uint64    `json:"revisionid"`
	Size       uint64    `json:"size"`
	Hash       string    `json:"hash"`
	Created    time.Time `json:"created"`
}

type UserInfo struct {
	Error
	UserID         uint64    `json:"userid"`
	Email          string    `json:"email"`
	EmailVerified  bool      `json:"emailverified"`
	Registered     time.Time `json:"registered"`
	Language       string    `json:"language"`
	Premium        bool      `json:"premium"`
	PremiumExpires time.Time `json:"premiumexpires,omitempty"`
	Quota          uint64    `json:"quota"`
	UsedQuota      uint64    `json:"usedquota"`
}

type FileLink struct {
	Error
	Path    string   `json:"path"`
	Expires string   `json:"expires"`
	Hosts   []string `json:"hosts"`
}

func (f *FileLink) URL() string {
	if len(f.Hosts) == 0 {
		return ""
	}
	return "https://" + f.Hosts[0] + f.Path
}
