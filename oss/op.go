package oss

import "io"

func UploadFromFile(ossPath, localFilePath string) error {
	return bucket.PutObjectFromFile(ossPath, localFilePath)
}

func UploadFromReader(ossPath string, srcReader io.Reader) error {
	return bucket.PutObject(ossPath, srcReader)
}
