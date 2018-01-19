# s3sync
Local files auto sync to cloud storages with S3API.

It cached file's `mtime` and `fsize` for performance, so sync the difference only.

### Installation 

Use `go get github.com/songjiayang/s3sync` or download the binary [relelase](https://github.com/songjiayang/s3sync/releases).

### Configuration

`s3sync` run with config file `config.json`, of course you can change it with `-config` option. 

The config details are:

```
{
  "root": "/home/user/example", // target folder
  "scan_worker": 20,  // the number of works to scan the files changes.
  "db": "./data/db", // cache files
  "s3sync": {
    "s3": {  //example with qiniu
      "access_key_id": "",
      "secret_access_key":"",
      "host":"https://api-s3.qiniu.com/",
      "region":"cn-east-1",
      "bucket":"test"
    },
    "worker": 20 // the number of works for s3sync.
  },
  "trim": true, // for upload file prefix trim
  "interval": 30
}

```

### Usage

You can run `s3sync -h` to check all options：

```
-config string
    the config file (default "./config.json")
-download
    download files from storage
-lu
    list upload status
-luf
    list all upload failed files
-r string
    root path (default "./data/test")
-upload
    upload files to storage
    
-d	run sync task backgound with interval time, default: 30s
```

Use case example:

- sync files upload only

```
$ s3sync -upload

2018/01/19 19:33:19 start scan job .......
2018/01/19 19:33:19 total files is 5
2018/01/19 19:33:19 end scan job .......
2018/01/19 19:33:19 start upload job .......
2018/01/19 19:33:19 left upload files 5
2018/01/19 19:33:20 end upload job .......
2018/01/19 19:33:20 start touch db .......
2018/01/19 19:33:20 end touch db .......
```

- sync files download only

```
$ s3sync -download

2018/01/19 19:33:44 start download job .......
2018/01/19 19:33:45 total finish download 5
2018/01/19 19:33:45 end download job .......
```

- sync files upload and download 

```
$ s3sync -upload -download

2018/01/19 19:34:57 start scan job .......
2018/01/19 19:34:57 total files is 5
2018/01/19 19:34:57 end scan job .......
2018/01/19 19:34:57 start upload job .......
2018/01/19 19:34:57 left upload files 5
2018/01/19 19:34:58 end upload job .......
2018/01/19 19:34:58 start touch db .......
2018/01/19 19:34:58 end touch db .......
2018/01/19 19:34:58 start download job .......
2018/01/19 19:34:58 scanned files 5
2018/01/19 19:34:58 left upload files 0
2018/01/19 19:34:58 total finish download 5
2018/01/19 19:34:58 end download job .......
```

Tips: If you want auto sync, use `-d` option please.


### Supported Cloud Storage 

- AWS
- Qiniu
- S3 Protocol Compatible Storage
