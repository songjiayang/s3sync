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
	files            map[string]*ff.Model
	UploadFilesCount int
}

func NewStorage(path string) *DiskStorage {
	return &DiskStorage{
		path:  path,
		files: make(map[string]*ff.Model),
	}
}

func (s *DiskStorage) InitDB(files []*ff.Model) {
	util.GzWrite(s.path, files)
}

func (s *DiskStorage) Recover() {
	b, err := util.GzRead(s.path)
	if err != nil {
		fmt.Printf("ioutil.ReadAll(): %v\n", err)
		return
	}

	var files []*ff.Model
	err = json.Unmarshal(b, &files)
	if err != nil {
		fmt.Printf("json.Unmarshal(): %v\n", err)
		return
	}

	for _, fileModel := range files {
		s.files[fileModel.Name] = fileModel
	}
}

func (s *DiskStorage) MapFiles() map[string]*ff.Model {
	return s.files
}

func (s *DiskStorage) UpdateMapFiles(files []*ff.Model) {
	for _, file := range files {
		old, ok := s.files[file.Name]
		if !ok {
			s.files[file.Name] = file
			s.UploadFilesCount++
		} else if old.DiffWith(file) {
			old.Modify()
			old.CopyWith(file)
			s.UploadFilesCount++
		} else if old.Status.IsUploadedFailed() {
			s.UploadFilesCount++
		}
	}
}

func (s *DiskStorage) PersistDB() {
	s.mux.Lock()
	files := make([]*ff.Model, len(s.files))
	i := 0
	for _, fileModel := range s.files {
		files[i] = fileModel
		i++
	}
	s.InitDB(files)
	s.mux.Unlock()
}

func (s *DiskStorage) Status() {
	for {
		if s.UploadFilesCount == 0 {
			break
		}
		time.Sleep(time.Second)
		log.Println("left upload files", s.UploadFilesCount)
	}
}
