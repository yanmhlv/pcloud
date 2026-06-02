package pcloud

import "strconv"

const (
	paramFileID         = "fileid"
	paramFolderID       = "folderid"
	paramToFolderID     = "tofolderid"
	paramPath           = "path"
	paramName           = "name"
	paramToName         = "toname"
	paramFileName       = "filename"
	paramRevisionID     = "revisionid"
	paramLinkID         = "linkid"
	paramShareID        = "shareid"
	paramShareRequestID = "sharerequestid"
	paramMail           = "mail"
)

func formatUint(v uint64) string {
	return strconv.FormatUint(v, 10)
}
