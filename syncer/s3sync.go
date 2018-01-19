package syncer

import (
	"sync"

	"github.com/songjiayang/s3sync/file"
	"github.com/songjiayang/s3sync/storage"
)

type Config struct {
	Workers  int       `json:"worker"`
	S3Config *S3Config `json:"s3"`
}

func Run(stg *storage.DiskStorage, cfg *Config, trim bool, root string, sResult *SyncResult) {
	queue := make(chan *file.FileModel, cfg.Workers)

	var wg sync.WaitGroup
	wg.Add(cfg.Workers)

	s3client := newS3(cfg.S3Config)

	for i := 0; i < cfg.Workers; i++ {
		go func(wg *sync.WaitGroup, queue chan *file.FileModel, sr *SyncResult) {
			defer wg.Done()
			for fileModel := range queue {
				err := s3client.uploadFile(fileModel.Name, trim, root)
				// maybe there need write lock
				if err == nil {
					fileModel.Uploaded()
				} else {
					fileModel.UploadedFailed()
				}

				// update sync result
				go func(sr *SyncResult, fileModel *file.FileModel, succeed bool) {
					sr.mux.Lock()
					if succeed {
						sr.SucceedCount++
						sr.SucceedCapacity += fileModel.Size
					} else {
						sr.FailedCount++
						sr.FailedCapacity += fileModel.Size
						sr.FailedFiles = append(sr.FailedFiles, fileModel.Name)
					}
					sr.mux.Unlock()
				}(sr, fileModel, err == nil)

			}
		}(&wg, queue, sResult)
	}

	for _, fileModel := range stg.MapFiles() {
		if fileModel.IsNeedUpload() {
			stg.UploadFilesCount--
			queue <- fileModel
		}
	}

	close(queue)
	wg.Wait()
}
