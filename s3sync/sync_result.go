package s3sync

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/file_scan/util"
)

type SyncResult struct {
	name            string
	SucceedCount    int64    `json:"succeedCount"`
	SucceedCapacity int64    `json:"succeedCapacity"`
	FailedCount     int64    `json:"failedCount"`
	FailedCapacity  int64    `json:"failedCapacity"`
	FailedFiles     []string `json:"failedFiles"`
	StartAt         int64    `json:"startAt"`
	EndAt           int64    `json:"endAt"`

	mux sync.RWMutex
}

func NewSyncResult(name string) *SyncResult {
	return &SyncResult{
		name:        name,
		StartAt:     time.Now().Unix(),
		FailedFiles: make([]string, 0),
	}
}

func (this *SyncResult) End() {
	this.EndAt = time.Now().Unix()
}

func (this *SyncResult) Persist() {
	util.GzWrite(this.path(), this)
}

func (this *SyncResult) Recover() (err error) {
	b, err := util.GzRead(this.path())
	if err != nil {
		fmt.Printf("ioutil.ReadAll(): %v\n", err)
		return
	}

	err = json.Unmarshal(b, this)
	if err != nil {
		fmt.Printf("json.Unmarshal(): %v\n", err)
	}

	return
}

func (this *SyncResult) PrintStat() {
	fmt.Println(this.printTitle())
	fmt.Println("Start at:", time.Unix(this.StartAt, 0))
	fmt.Println("End at:", time.Unix(this.EndAt, 0))
	fmt.Println("Succeed count:", this.SucceedCount)
	fmt.Println("Succeed capacity:", this.SucceedCapacity)
	fmt.Println("Failed count:", this.FailedCount)
	fmt.Println("Failed capacity:", this.FailedCapacity)
	fmt.Println(this.printTitle())
}

func (this *SyncResult) PrintFailedFiles() {
	for _, f := range this.FailedFiles {
		fmt.Println(f)
	}
}

func (this *SyncResult) path() string {
	return "./data/" + this.name + "_result"
}

func (this *SyncResult) printTitle() string {
	return "-------- " + this.name + " status ---------"
}
