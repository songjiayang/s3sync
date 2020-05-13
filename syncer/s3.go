package syncer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Config struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Host            string `json:"host"`
	Region          string `json:"region"`
	Bucket          string `json:"bucket"`
}

type s3client struct {
	s3service  *s3.S3
	cfg        *S3Config
	downloader *s3manager.Downloader
}

func newS3(cfg *S3Config) *s3client {
	s3provider := credentials.StaticProvider{
		Value: credentials.Value{
			AccessKeyID:     cfg.AccessKeyID,
			SecretAccessKey: cfg.SecretAccessKey,
		},
	}
	s3credential := credentials.NewCredentials(&s3provider)
	s3session := session.New()

	// s3 config
	s3session.Config.WithEndpoint(cfg.Host)
	s3session.Config.WithRegion(cfg.Region)
	s3session.Config.WithCredentials(s3credential)
	s3session.Config.WithMaxRetries(1)
	// s3session.Config.WithS3ForcePathStyle(true)

	s3service := s3.New(s3session, nil)

	return &s3client{
		s3service: s3service,
		cfg:       cfg,
		downloader: s3manager.NewDownloader(s3session, func(d *s3manager.Downloader) {
			d.PartSize = 100 * 1024 * 1024
			d.Concurrency = 20
		}),
	}
}

func (this *s3client) uploadFile(file string, trim bool, root string) (err error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("read file with ERROR: %v\n", err)
		return
	}

	payload := &s3.PutObjectInput{
		Bucket: aws.String(this.cfg.Bucket),
		Body:   bytes.NewReader(buf),
	}

	// if prefix
	if trim {
		payload.Key = aws.String(strings.Replace(file, root, "", 1))
	} else {
		payload.Key = aws.String(file)
	}

	_, err = this.s3service.PutObject(payload)
	if err != nil {
		fmt.Printf("s3service.PutObject(%v): %v\n", payload, err)
	}

	return
}

func (this *s3client) listObjects(marker *string) (*s3.ListObjectsOutput, error) {
	return this.s3service.ListObjects(&s3.ListObjectsInput{
		Bucket:  aws.String(this.cfg.Bucket),
		MaxKeys: aws.Int64(100),
		Marker:  marker,
	})
}

func (this *s3client) downFile(root string, key *string) (err error) {
	fpath := filepath.Join(root, *key)
	os.MkdirAll(path.Dir(fpath), os.ModePerm)

	file, err := os.Create(fpath)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = this.downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(this.cfg.Bucket),
		Key:    key,
	})

	return
}
