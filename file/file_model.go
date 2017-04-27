package file

import "time"

type FileModel struct {
	Name    string     `json:'name'`
	Size    int64      `json:"size"`
	ModTime time.Time  `json:"modtime"`
	Status  fileStatus `json:"status"`
}

func NewFileModel(name string, size int64, modTime time.Time) *FileModel {
	return &FileModel{
		Name:    name,
		Size:    size,
		ModTime: modTime,
		Status:  fileNewStatus,
	}
}

func (this *FileModel) Uploaded() {
	this.Status = fileUploadedStatus
}

func (this *FileModel) Modify() {
	this.Status = fileModifyStatus
}

func (this *FileModel) UploadedFailed() {
	this.Status = fileUploadedFailedStatus
}

func (this *FileModel) IsNeedUpload() bool {
	return this.Status.IsNew() || this.Status.IsModify() || this.Status.IsUploadedFailed()
}

func (this *FileModel) DiffWith(target *FileModel) bool {
	return this.ModTime != target.ModTime || this.Size != target.Size
}

func (this *FileModel) CopyWith(target *FileModel) {
	this.ModTime = target.ModTime
	this.Size = target.Size
}
