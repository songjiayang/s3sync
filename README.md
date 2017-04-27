# s3_sync
Local files auto sync to cloud storages with S3API.
It cached file's mtime and fsize, so can upload the changes only.

### Config

```
{
  "root": "~",
  "scan_worker": 20,  // the number of works to scan the files changes.
  "db": "./data/db", // cache files
  "s3sync": {
    "s3": {
      "access_key_id": "",
      "secret_access_key":"",
      "host":"",
      "region":"",
      "bucket":""
    },
    "worker": 20 // the number of works for s3sync.
  },
  "trim": true // for upload file prefix trim
}

```

### Command

you can run `file_scan -h` to get all the options：

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
```
