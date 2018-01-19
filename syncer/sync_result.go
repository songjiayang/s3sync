package syncer

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/songjiayang/s3sync/util"
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
	log.Println(this.printTitle())
	log.Println("Start at:", time.Unix(this.StartAt, 0))
	log.Println("End at:", time.Unix(this.EndAt, 0))
	log.Println("Succeed count:", this.SucceedCount)
	log.Println("Succeed capacity:", this.SucceedCapacity)
	log.Println("Failed count:", this.FailedCount)
	log.Println("Failed capacity:", this.FailedCapacity)
	log.Println(this.printTitle())
}

func (this *SyncResult) PrintFailedFiles() {
	for _, f := range this.FailedFiles {
		log.Println(f)
	}
}

func (this *SyncResult) path() string {
	return "./data/" + this.name + "_result"
}

func (this *SyncResult) printTitle() string {
	return "-------- " + this.name + " status ---------"
}
