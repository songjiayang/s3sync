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

func (r *SyncResult) End() {
	r.EndAt = time.Now().Unix()
}

func (r *SyncResult) Persist() {
	util.GzWrite(r.path(), r)
}

func (r *SyncResult) Recover() (err error) {
	b, err := util.GzRead(r.path())
	if err != nil {
		fmt.Printf("ioutil.ReadAll(): %v\n", err)
		return
	}

	err = json.Unmarshal(b, r)
	if err != nil {
		fmt.Printf("json.Unmarshal(): %v\n", err)
	}

	return
}

func (r *SyncResult) PrintStat() {
	log.Println(r.printTitle())
	log.Println("Start at:", time.Unix(r.StartAt, 0))
	log.Println("End at:", time.Unix(r.EndAt, 0))
	log.Println("Succeed count:", r.SucceedCount)
	log.Println("Succeed capacity:", r.SucceedCapacity)
	log.Println("Failed count:", r.FailedCount)
	log.Println("Failed capacity:", r.FailedCapacity)
	log.Println(r.printTitle())
}

func (r *SyncResult) PrintFailedFiles() {
	for _, f := range r.FailedFiles {
		log.Println(f)
	}
}

func (r *SyncResult) path() string {
	return "./data/" + r.name + "_result"
}

func (r *SyncResult) printTitle() string {
	return "-------- " + r.name + " status ---------"
}
