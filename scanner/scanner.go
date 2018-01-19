package scanner

import (
	"log"
	"sync"
	"time"

	ff "github.com/songjiayang/s3sync/file"
	"github.com/songjiayang/s3sync/lib"
)

type Scanner struct {
	c       chan string
	folders []string
	files   []*ff.FileModel
	leftJob int // left scan folder number
	dmux    sync.RWMutex
	fmux    sync.RWMutex
}

func NewScanner(root string, w int) *Scanner {
	s := &Scanner{
		c:       make(chan string, w),
		folders: make([]string, 1, 128),
		files:   make([]*ff.FileModel, 0, 128),
		leftJob: 1,
	}

	s.folders[0] = root

	return s
}

func (this *Scanner) Scan() {
	var wg sync.WaitGroup
	wg.Add(1)

	// make the worker running
	go this.run(&wg)
	go this.pushJob(&wg)

	wg.Wait()
	close(this.c)
}

func (this *Scanner) Files() []*ff.FileModel {
	return this.files
}

func (this *Scanner) Status() {
	for {
		if this.leftJob == 0 {
			break
		}
		time.Sleep(time.Second)
		log.Println("scanned files", len(this.Files()))
	}
}

func (this *Scanner) pushJob(wg *sync.WaitGroup) {
	for {
		d := this.pop()

		this.dmux.Lock()
		leftJob := this.leftJob
		this.dmux.Unlock()

		if leftJob == 0 {
			wg.Done()
			break
		}

		if d != "" {
			this.c <- d
		}
	}
}

func (this *Scanner) run(wg *sync.WaitGroup) {
	for i := 0; i < cap(this.c); i++ {
		go this.list(wg)
	}
}

func (this *Scanner) list(wg *sync.WaitGroup) {
	for dir := range this.c {
		_folders, _files := lib.ListDir(dir)
		l := len(_folders) - 1

		this.dmux.Lock()
		this.folders = append(this.folders, _folders...)
		this.leftJob += l
		this.dmux.Unlock()

		this.fmux.Lock()
		this.files = append(this.files, _files...)
		this.fmux.Unlock()
	}
}

func (this *Scanner) pop() (ret string) {
	this.dmux.Lock()
	l := len(this.folders)
	if l > 0 {
		ret = this.folders[l-1]
		this.folders = this.folders[0 : l-1]
	}
	this.dmux.Unlock()
	return
}
