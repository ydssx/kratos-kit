package storage

import (
	"context"
	"path/filepath"

	"github.com/ydssx/kratos-kit/pkg/errors"
	"github.com/ydssx/kratos-kit/pkg/util"

	"cloud.google.com/go/storage"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/api/option"
)

type GoogleCloudStorage struct {
	bucket *storage.BucketHandle
	client *storage.Client
}

func NewGoogleCloudStorage(bucketName, projectID, credentialsFile string) (s *GoogleCloudStorage, cleanup func()) {
	// 创建 Google Cloud Storage 客户端
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(credentialsFile)))
	if err != nil {
		log.Fatalf("Failed to create Google Cloud Storage client: %v", err)
	}
	log.Info("Google Cloud Storage client created.")
	cleanup = func() { client.Close() }

	return &GoogleCloudStorage{bucket: client.Bucket(bucketName), client: client}, cleanup
}

func (g *GoogleCloudStorage) SaveFile(ctx context.Context, path, filename, contentType string, fileBytes []byte) (string, error) {
	objectName := util.MD5Bytes(fileBytes) + filepath.Ext(filename)
	if path != "" {
		objectName = filepath.Join(path, objectName)
	}
	// 在存储桶中创建一个新的对象，并设置对象属性
	object := g.bucket.Object(objectName)

	// 检查对象是否存在
	attrs, err := object.Attrs(ctx)
	if err == nil {
		return attrs.MediaLink, nil
	}

	// 打开对象的写入器
	wc := object.NewWriter(ctx)
	if contentType != "" {
		wc.ContentType = contentType
	}

	// 将文件内容复制到对象的写入器中
	if _, err := wc.Write(fileBytes); err != nil {
		return "", errors.Wrap(err, "上传文件失败")
	}

	// 关闭写入器以完成上传
	if err := wc.Close(); err != nil {
		return "", errors.Wrap(err, "关闭写入器失败")
	}

	attrs, err = object.Attrs(ctx)
	if err != nil {
		return "", errors.Wrap(err, "获取对象属性失败")
	}

	log.Info("Object uploaded successfully. ", "object: ", attrs.Name, "mediaLink: ", attrs.MediaLink)

	return attrs.MediaLink, nil
}

func (g *GoogleCloudStorage) DeleteFile(ctx context.Context, filename string) error {
	object := g.bucket.Object(filename)
	err := object.Delete(ctx)
	if err != nil {
		return errors.Wrap(err, "删除文件失败")
	}

	log.Info("Object deleted successfully. ", "object: ", filename)
	return nil
}
