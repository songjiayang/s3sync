package main

import (
	"flag"
	"fmt"

	"github.com/file_scan/config"
	"github.com/file_scan/s3sync"
	"github.com/file_scan/scanner"
	"github.com/file_scan/storage"
)

var (
	conf     string
	rootPath string
	upload   bool
	download bool
	lu       bool
	luf      bool

	defaultRootpath = "./data/test"
)

func main() {
	flag.StringVar(&conf, "config", "./config.json", "the config file")
	flag.StringVar(&rootPath, "r", defaultRootpath, "root path")
	flag.BoolVar(&upload, "upload", false, "upload files to storage")
	flag.BoolVar(&download, "download", false, "download files from storage")
	flag.BoolVar(&lu, "lu", false, "list upload status")
	flag.BoolVar(&luf, "luf", false, "list all upload failed files")

	flag.Parse()

	cfg, err := config.LoadFile(conf)
	if err != nil {
		fmt.Printf("config.LoadFile(%s), %v\n", conf, err)
		return
	}

	// fix rootPath
	if rootPath != defaultRootpath {
		cfg.Root = rootPath
	} else if cfg.Root == "" {
		cfg.Root = defaultRootpath
	}

	printUploadStatus()

	if upload {
		syncWithPush(cfg)
	}

	if download {
		syncWithPull(cfg)
	}
}

func syncWithPush(cfg *config.Config) {
	fmt.Println("start scan job .......")
	s := scanner.NewScanner(cfg.Root, cfg.ScanWorker)
	go s.Status()
	s.Scan()
	fmt.Println("total files is", len(s.Files()))
	fmt.Println("end scan job .......")

	stg := storage.NewStorage(cfg.DB)
	stg.Recover()
	stg.UpdateMapFiles(s.Files())

	fmt.Println("start upload job .......")
	fmt.Println("left upload files", stg.UploadFilesCount)
	go stg.Status()
	syncResult := s3sync.NewSyncResult("upload")
	s3sync.Run(stg, cfg.SyncConfig, cfg.Trim, cfg.Root, syncResult)
	fmt.Println("end upload job .......")

	syncResult.End()
	syncResult.Persist()

	fmt.Println("start touch db .......")
	stg.PersistDB()
	fmt.Println("end touch db .......")
}

func syncWithPull(cfg *config.Config) {
	fmt.Println("start download job .......")
	s3sync.RunDownloader(cfg.SyncConfig, cfg.Root)
	fmt.Println("end download job .......")
}

func printUploadStatus() {
	if lu || luf {
		rs := s3sync.NewSyncResult("upload")
		err := rs.Recover()
		if err != nil {
			fmt.Println("can not get upload result data")
			return
		}

		if lu {
			rs.PrintStat()
		}

		if luf {
			rs.PrintFailedFiles()
		}
	}
}
