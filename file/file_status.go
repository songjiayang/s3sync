package file

type fileStatus int

const (
	fileMinStatus fileStatus = iota
	fileNewStatus
	fileUploadedStatus
	fileModifyStatus
	fileDeletedStatus
	fileUploadedFailedStatus
	fileMaxStatus
)

func (this fileStatus) IsNew() bool {
	return this == fileNewStatus
}

func (this fileStatus) IsModify() bool {
	return this == fileModifyStatus
}

func (this fileStatus) IsUploaded() bool {
	return this == fileUploadedStatus
}

func (this fileStatus) IsUploadedFailed() bool {
	return this == fileUploadedFailedStatus
}
