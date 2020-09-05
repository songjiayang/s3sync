package file

type Status int

const (
	_minStatus Status = iota
	NewStatus
	UploadedStatus
	ModifyStatus
	DeletedStatus
	UploadedFailedStatus
	_maxStatus
)

func (s Status) IsNew() bool {
	return s == NewStatus
}

func (s Status) IsModify() bool {
	return s == ModifyStatus
}

func (s Status) IsUploaded() bool {
	return s == UploadedStatus
}

func (s Status) IsUploadedFailed() bool {
	return s == UploadedFailedStatus
}
