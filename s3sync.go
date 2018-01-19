package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/songjiayang/s3sync/config"
	"github.com/songjiayang/s3sync/scanner"
	"github.com/songjiayang/s3sync/storage"
	"github.com/songjiayang/s3sync/syncer"
)

var (
	conf                string
	rootPath            string
	upload, download, d bool

	lu, luf bool

	defaultRootpath = "./data/test"
)

func init() {
	flag.StringVar(&conf, "config", "./config.json", "the config file")
	flag.StringVar(&rootPath, "r", defaultRootpath, "root path")

	flag.BoolVar(&upload, "upload", false, "upload files to storage")
	flag.BoolVar(&download, "download", false, "download files from storage")
	flag.BoolVar(&lu, "lu", false, "list upload status")
	flag.BoolVar(&luf, "luf", false, "list all upload failed files")
	flag.BoolVar(&d, "d", false, "run sync task backgound with interval time, default: 30s")
}

func main() {
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

	syncStatus()

	interval := cfg.Duration()
	for {
		if upload {
			push(cfg)
		}

		if download {
			pull(cfg)
		}

		if !d {
			break
		}

		time.Sleep(interval)
		log.Println("---------- interval loop ----------")
	}
}

func push(cfg *config.Config) {
	log.Println("start scan job .......")
	s := scanner.NewScanner(cfg.Root, cfg.ScanWorker)
	go s.Status()
	s.Scan()
	log.Println("total files is", len(s.Files()))
	log.Println("end scan job .......")

	stg := storage.NewStorage(cfg.DB)
	stg.Recover()
	stg.UpdateMapFiles(s.Files())

	log.Println("start upload job .......")
	log.Println("left upload files", stg.UploadFilesCount)
	go stg.Status()
	syncResult := syncer.NewSyncResult("upload")
	syncer.Run(stg, cfg.SyncConfig, cfg.Trim, cfg.Root, syncResult)
	log.Println("end upload job .......")

	syncResult.End()
	syncResult.Persist()

	log.Println("start touch db .......")
	stg.PersistDB()
	log.Println("end touch db .......")
}

func pull(cfg *config.Config) {
	log.Println("start download job .......")
	syncer.RunDownloader(cfg.SyncConfig, cfg.Root)
	log.Println("end download job .......")
}

func syncStatus() {
	if lu || luf {
		rs := syncer.NewSyncResult("upload")
		err := rs.Recover()
		if err != nil {
			log.Println("can not get upload result data")
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
