package util

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStorage(endpoint, accessKey, secretKey, bucketName string, secure bool) (ObjectStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Minio client: %w", err)
	}

	// Check if the bucket exists. If not, create it.
	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket exists: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &minioStorage{client: client, bucketName: bucketName}, nil
}

// DeleteObject implements storage.ObjectStorage.
func (m *minioStorage) DeleteObject(ctx context.Context, objectName string) error {
	// Delete an object from the bucket
	err := m.client.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

// Download implements storage.ObjectStorage.
func (m *minioStorage) Download(ctx context.Context, objectName string) (io.ReadCloser, error) {
	// Download an object from the bucket
	obj, err := m.client.GetObject(ctx, m.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// ListObjects implements storage.ObjectStorage.
func (m *minioStorage) ListObjects(ctx context.Context, prefix string) ([]string, error) {
	// List objects in the bucket with the specified prefix
	objectCh := m.client.ListObjects(ctx, m.bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	var objects []string
	for obj := range objectCh {
		if obj.Err != nil {
			return nil, obj.Err
		}
		objects = append(objects, obj.Key)
	}
	return objects, nil
}

// ListObjectStream streams object names to a channel to avoid high memory usage.
func (m *minioStorage) ListObjectStream(ctx context.Context, prefix string) (<-chan string, error) {
	objectCh := m.client.ListObjects(ctx, m.bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	resultCh := make(chan string)

	go func() {
		defer close(resultCh)
		for obj := range objectCh {
			if obj.Err != nil {
				// Handle error
				return
			}
			select {
			case resultCh <- obj.Key:
			case <-ctx.Done():
				return
			}
		}
	}()

	return resultCh, nil
}

// PresignURL implements storage.ObjectStorage.
func (m *minioStorage) PresignURL(
	ctx context.Context,
	objectName string,
	expiry time.Duration,
	method string,
) (string, error) {
	switch method {
	case "GET":
		// 下载类型预签名 URL
		reqParams := make(url.Values)
		reqParams.Set("response-content-disposition", "attachment")
		presignedURL, err := m.client.PresignedGetObject(ctx, m.bucketName, objectName, expiry, reqParams)
		if err != nil {
			return "", err
		}
		return presignedURL.String(), nil

	case "PUT":
		// 上传类型预签名 URL
		presignedURL, err := m.client.PresignedPutObject(ctx, m.bucketName, objectName, expiry)
		if err != nil {
			return "", err
		}
		return presignedURL.String(), nil

	default:
		return "", fmt.Errorf("unsupported method: %s, must be GET or PUT", method)
	}
}

// Upload implements storage.ObjectStorage.
func (m *minioStorage) Upload(ctx context.Context, objectName string, r io.Reader, contentType string) error {
	var objectSize int64 = -1
	// Try to determine the size of the reader
	if s, ok := r.(interface {
		Size() int64
	}); ok {
		objectSize = s.Size()
	} else if f, ok := r.(*os.File); ok {
		if stat, err := f.Stat(); err == nil {
			objectSize = stat.Size()
		}
	}

	_, err := m.client.PutObject(ctx, m.bucketName, objectName, r, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}
	return nil
}
