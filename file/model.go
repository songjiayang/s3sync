package file

import "time"

type Model struct {
	Name    string    `json:'name'`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modtime"`
	Status  Status    `json:"status"`
}

func NewModel(name string, size int64, modTime time.Time) *Model {
	return &Model{
		Name:    name,
		Size:    size,
		ModTime: modTime,
		Status:  NewStatus,
	}
}

func (m *Model) Uploaded() {
	m.Status = UploadedStatus
}

func (m *Model) Modify() {
	m.Status = ModifyStatus
}

func (m *Model) UploadedFailed() {
	m.Status = UploadedFailedStatus
}

func (m *Model) IsNeedUpload() bool {
	return m.Status.IsNew() || m.Status.IsModify() || m.Status.IsUploadedFailed()
}

func (m *Model) DiffWith(target *Model) bool {
	return m.ModTime != target.ModTime || m.Size != target.Size
}

func (m *Model) CopyWith(target *Model) {
	m.ModTime = target.ModTime
	m.Size = target.Size
}
