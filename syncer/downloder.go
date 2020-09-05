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
	downloader := &Downloader{
		root: root,
		mux:  &sync.Mutex{},
		cfg:  cfg,
	}

	go downloader.status()
	downloader.run()
}

func (d *Downloader) run() {
	queue := make(chan *string, d.cfg.Workers)
	var wg sync.WaitGroup
	wg.Add(d.cfg.Workers)

	s3client := newS3(d.cfg.S3Config)

	for i := 0; i < d.cfg.Workers; i++ {
		go func(wg *sync.WaitGroup, queue chan *string) {
			defer wg.Done()
			for key := range queue {
				err := s3client.downFile(d.root, key)
				if err == nil {
					d.mux.Lock()
					d.finish += 1
					d.mux.Unlock()
				} else {
					fmt.Printf("download file with err: <%s>, <%v> \n", *key, err)
				}
			}
		}(&wg, queue)
	}

	var marker *string

	for !d.isFinished {
		result, err := s3client.listObjects(marker)
		if err != nil {
			fmt.Printf("list objects with error %v\n", err)
			d.isFinished = true
			return
		}

		marker = result.NextMarker
		for _, obj := range result.Contents {
			queue <- obj.Key
		}
		d.isFinished = !*result.IsTruncated
	}

	close(queue)
	wg.Wait()

	log.Println("total finish download", d.finish)
}

func (d *Downloader) status() {
	for !d.isFinished {
		time.Sleep(3 * time.Second)
		log.Println("finish download", d.finish)
	}
}
