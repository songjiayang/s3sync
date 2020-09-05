package lib

import (
	"io/ioutil"

	ff "github.com/songjiayang/s3sync/file"
)

func ListDir(dir string) (folders []string, files []*ff.Model) {
	folders = make([]string, 0)
	files = make([]*ff.Model, 0)

	dirFiles, _ := ioutil.ReadDir(dir)
	for _, file := range dirFiles {
		if file.IsDir() {
			folders = append(folders, dir+"/"+file.Name())
		} else {
			fileModel := ff.NewModel(dir+"/"+file.Name(), file.Size(), file.ModTime())
			files = append(files, fileModel)
		}
	}

	return
}
