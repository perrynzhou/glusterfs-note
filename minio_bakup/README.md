
## split tar 

```
split -b 30m  minio-RELEASE.2020-10-18T21-54-12Z.tar.gz  minio-RELEASE.2020-10-18T21-54-12Z.tar.gz.  –verbose
```

## merge tar 

```
cat  minio-RELEASE.2020-10-18T21-54-12Z.tar.gz.a* minio-RELEASE.2020-10-18T21-54-12Z.tar.gz 
```
