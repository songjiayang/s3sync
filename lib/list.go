package lib

import (
	"io/ioutil"

	ff "github.com/file_scan/file"
)

func ListDir(dir string) (folders []string, files []*ff.FileModel) {
	folders = make([]string, 0)
	files = make([]*ff.FileModel, 0)

	dirfiles, _ := ioutil.ReadDir(dir)

	for _, file := range dirfiles {
		if file.IsDir() {
			folders = append(folders, dir+"/"+file.Name())
		} else {
			fileModel := ff.NewFileModel(dir+"/"+file.Name(), file.Size(), file.ModTime())
			files = append(files, fileModel)
		}
	}

	return
}
