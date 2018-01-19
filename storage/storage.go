package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	ff "github.com/songjiayang/s3sync/file"
	"github.com/songjiayang/s3sync/util"
)

type DiskStorage struct {
	path             string
	mux              sync.RWMutex
	files            map[string]*ff.FileModel
	UploadFilesCount int
}

func NewStorage(path string) *DiskStorage {
	return &DiskStorage{
		path:  path,
		files: make(map[string]*ff.FileModel),
	}
}

func (this *DiskStorage) InitDB(files []*ff.FileModel) {
	util.GzWrite(this.path, files)
}

func (this *DiskStorage) Recover() {
	b, err := util.GzRead(this.path)
	if err != nil {
		fmt.Printf("ioutil.ReadAll(): %v\n", err)
		return
	}

	var files []*ff.FileModel
	err = json.Unmarshal(b, &files)
	if err != nil {
		fmt.Printf("json.Unmarshal(): %v\n", err)
		return
	}

	for _, fileModel := range files {
		this.files[fileModel.Name] = fileModel
	}
}

func (this *DiskStorage) MapFiles() map[string]*ff.FileModel {
	return this.files
}

func (this *DiskStorage) UpdateMapFiles(files []*ff.FileModel) {
	for _, file := range files {
		old, ok := this.files[file.Name]
		if !ok {
			this.files[file.Name] = file
			this.UploadFilesCount++
		} else if old.DiffWith(file) {
			old.Modify()
			old.CopyWith(file)
			this.UploadFilesCount++
		} else if old.Status.IsUploadedFailed() {
			this.UploadFilesCount++
		}
	}
}

func (this *DiskStorage) PersistDB() {
	this.mux.Lock()
	files := make([]*ff.FileModel, len(this.files))
	i := 0
	for _, fileModel := range this.files {
		files[i] = fileModel
		i++
	}
	this.InitDB(files)
	this.mux.Unlock()
}

func (this *DiskStorage) Status() {
	for {
		if this.UploadFilesCount == 0 {
			break
		}
		time.Sleep(time.Second)
		log.Println("left upload files", this.UploadFilesCount)
	}
}
