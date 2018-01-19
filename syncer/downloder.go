package syncer

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Downloader struct {
	cfg        *Config
	root       string
	finish     int64
	isFinished bool
	marker     string

	mux *sync.Mutex
}

func RunDownloader(cfg *Config, root string) {
	downloder := &Downloader{
		root: root,
		mux:  &sync.Mutex{},
		cfg:  cfg,
	}

	go downloder.status()
	downloder.run()
}

func (this *Downloader) run() {
	queue := make(chan *string, this.cfg.Workers)
	var wg sync.WaitGroup
	wg.Add(this.cfg.Workers)

	s3client := newS3(this.cfg.S3Config)

	for i := 0; i < this.cfg.Workers; i++ {
		go func(wg *sync.WaitGroup, queue chan *string) {
			defer wg.Done()
			for key := range queue {
				err := s3client.downFile(this.root, key)
				if err == nil {
					this.mux.Lock()
					this.finish += 1
					this.mux.Unlock()
				} else {
					fmt.Printf("download file with err: <%s>, <%v> \n", *key, err)
				}
			}
		}(&wg, queue)
	}

	var marker *string

	for !this.isFinished {
		result, err := s3client.listObjects(marker)
		if err != nil {
			fmt.Printf("list objects with error %v\n", err)
			this.isFinished = true
			return
		}

		marker = result.NextMarker
		for _, obj := range result.Contents {
			queue <- obj.Key
		}
		this.isFinished = !*result.IsTruncated
	}

	close(queue)
	wg.Wait()

	log.Println("total finish download", this.finish)
}

func (this *Downloader) status() {
	for !this.isFinished {
		time.Sleep(3 * time.Second)
		log.Println("finish download", this.finish)
	}
}
