package storage

import "io"

type Storage interface {
	Upload(bucket, objectName string, data []byte, contentType string) error
	Download(bucket, objectName string) ([]byte, error)
	Delete(bucket, objectName string) error
	Exists(bucket, objectName string) bool
	Serve(bucket, objectName string, w io.Writer) error
	GetURL(bucket, objectName string) string
}
