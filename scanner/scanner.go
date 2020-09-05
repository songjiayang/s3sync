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
	files   []*ff.Model
	leftJob int // left scan folder number
	dmux    sync.RWMutex
	fmux    sync.RWMutex
}

func NewScanner(root string, w int) *Scanner {
	s := &Scanner{
		c:       make(chan string, w),
		folders: make([]string, 1, 128),
		files:   make([]*ff.Model, 0, 128),
		leftJob: 1,
	}

	s.folders[0] = root

	return s
}

func (s *Scanner) Scan() {
	var wg sync.WaitGroup
	wg.Add(1)

	// make the worker running
	go s.run(&wg)
	go s.pushJob(&wg)

	wg.Wait()
	close(s.c)
}

func (s *Scanner) Files() []*ff.Model {
	return s.files
}

func (s *Scanner) Status() {
	for {
		if s.leftJob == 0 {
			break
		}
		time.Sleep(time.Second)
		log.Println("scanned files", len(s.Files()))
	}
}

func (s *Scanner) pushJob(wg *sync.WaitGroup) {
	for {
		d := s.pop()

		s.dmux.Lock()
		leftJob := s.leftJob
		s.dmux.Unlock()

		if leftJob == 0 {
			wg.Done()
			break
		}

		if d != "" {
			s.c <- d
		}
	}
}

func (s *Scanner) run(wg *sync.WaitGroup) {
	for i := 0; i < cap(s.c); i++ {
		go s.list(wg)
	}
}

func (s *Scanner) list(wg *sync.WaitGroup) {
	for dir := range s.c {
		_folders, _files := lib.ListDir(dir)
		l := len(_folders) - 1

		s.dmux.Lock()
		s.folders = append(s.folders, _folders...)
		s.leftJob += l
		s.dmux.Unlock()

		s.fmux.Lock()
		s.files = append(s.files, _files...)
		s.fmux.Unlock()
	}
}

func (s *Scanner) pop() (ret string) {
	s.dmux.Lock()
	l := len(s.folders)
	if l > 0 {
		ret = s.folders[l-1]
		s.folders = s.folders[0 : l-1]
	}
	s.dmux.Unlock()
	return
}
