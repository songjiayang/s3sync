package util

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func GzWrite(fpath string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("json.Marshal(): %v\n", err)
		return
	}

	var fileGZ bytes.Buffer
	zipper := gzip.NewWriter(&fileGZ)
	_, err = zipper.Write(b)
	if err != nil {
		fmt.Printf("zipper.Write ERROR: %+v", err)
		return
	}
	zipper.Close()

	os.MkdirAll(path.Dir(fpath), os.ModePerm)

	err = ioutil.WriteFile(fpath, fileGZ.Bytes(), 0644)
	if err != nil {
		fmt.Printf("ioutil.WriteFile(%s): %v\n", fpath, err)
		return
	}
}

func GzRead(fpath string) (b []byte, err error) {
	rzip, err := ioutil.ReadFile(fpath)
	if err != nil {
		fmt.Printf("ioutil.ReadFile(): %v\n", err)
		return
	}
	r, err := gzip.NewReader(bytes.NewBuffer(rzip))
	if err != nil {
		fmt.Printf("gzip.NewReader(): %v\n", err)
		return
	}

	b, err = ioutil.ReadAll(r)
	if err != nil {
		fmt.Printf("ioutil.ReadAll(): %v\n", err)
	}

	return
}
