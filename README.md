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
    "worker": 20 // the number of works to upload files.
  }
}

```


