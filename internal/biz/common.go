// Package biz contains common business logic used across the application.
package biz

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ydssx/kratos-kit/models"
	"github.com/ydssx/kratos-kit/pkg/errors"
	"github.com/ydssx/kratos-kit/pkg/logger"
	"github.com/ydssx/kratos-kit/pkg/storage"
	"github.com/ydssx/kratos-kit/pkg/util"
)

// CommonUseCase 通用业务逻辑, 用于业务上的公共方法
type CommonUseCase struct {
	tx       Transaction
	store    storage.Storage
	userRepo UserRepo
}

func NewCommonUseCase(
	tx Transaction,
	store storage.Storage,
	userRepo UserRepo,
) *CommonUseCase {
	return &CommonUseCase{
		tx:       tx,
		store:    store,
		userRepo: userRepo,
	}
}

// UploadFile 上传文件
func (uc *CommonUseCase) UploadFile(ctx context.Context, userID, userType int, file *multipart.FileHeader) (fileMetadata *models.FileMetadata, err error) {
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, errors.Errorf("Failed to open file: %v", err)
	}
	defer src.Close()

	contentType := file.Header.Get("Content-Type")

	// 读取文件内容
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return nil, errors.Errorf("Failed to read file: %v", err)
	}

	file_md5 := util.MD5Bytes(fileBytes)
	fileInfo, err := models.NewFileMetadataModel().SetMd5(file_md5).SetUserId(int64(userID)).FirstOne()
	if err == nil {
		return &fileInfo, nil
	}

	fileMetadata = &models.FileMetadata{
		UserId:   userID,
		FileSize: len(fileBytes),
		FileMd5:  file_md5,
		FileType: models.FileTypeImage,
	}

	if util.IsAudioFile(file.Filename) {
		fileMetadata.FileType = models.FileTypeAudio
	}

	// 构建文件名
	originalName := filepath.Base(file.Filename)
	ext := filepath.Ext(originalName)
	storedName := time.Now().Format("20060102150405") + util.GenerateCode(2) + ext
	var folder string
	if userType == int(models.UserTypeNormal) {
		folder = fmt.Sprintf("user_uploads/%d", userID)
	} else {
		folder = "system_uploads"
	}
	objectPath := fmt.Sprintf("%s/%s", folder, storedName)
	fileMetadata.Filename = file_md5 + ext

	// 创建本地存储路径
	uploadDir := "D:\\uploads\\"
	if runtime.GOOS == "linux" {
		uploadDir = "/tmp/uploads/"
	}
	filePath := filepath.Join(uploadDir, storedName)
	if err := util.SaveUploadedFile(file, filePath); err != nil {
		logger.Errorf(ctx, "Failed to save uploaded file: %v", err)
		return nil, err
	}
	defer os.RemoveAll(uploadDir) // 删除本地存储路径

	if util.IsVideoFile(originalName) {
		fileMetadata.FileType = models.FileTypeVideo

		// 获取视频元数据
		videoMetadata, err := util.GetVideoMetadata(filePath)
		if err != nil {
			logger.Error(ctx, "Failed to get video metadata:", err)
			return nil, err
		}
		fileMetadata.Width = videoMetadata.Width
		fileMetadata.Height = videoMetadata.Height
		fileMetadata.Fps = videoMetadata.FPS
		fileMetadata.VideoDuration = videoMetadata.Duration
		fileMetadata.Encoding = videoMetadata.CodecName

		// 如果不是H264视频，转码成H264
		if !videoMetadata.IsH264 {
			outputPath := filepath.Join(uploadDir, filepath.Base(filePath)+"_h264.mp4")
			err := util.ConvertToH264(filePath, outputPath)
			if err != nil {
				logger.Error(ctx, "Failed to convert video to h264:", err)
				return nil, err
			}

			fileBytes, err = os.ReadFile(outputPath)
			if err != nil {
				logger.Error(ctx, "Failed to read converted video:", err)
				return nil, err
			}
			fileMetadata.FileSize = len(fileBytes)
		}

		// 生成视频缩略图
		thumbnailPath, err := util.GenerateThumbnail(filePath)
		if err != nil {
			logger.Error(ctx, "Failed to generate thumbnail:", err)
			return nil, err
		}

		// 保存缩略图到存储服务
		thumbnailBytes, err := os.ReadFile(thumbnailPath)
		if err != nil {
			logger.Error(ctx, "Failed to read thumbnail:", err)
			return nil, err
		}

		thumbnailSavePath := fmt.Sprintf("%s/%s", folder, filepath.Base(thumbnailPath))
		thumbnailURL, err := uc.store.SaveFile(context.Background(), folder, thumbnailSavePath, "", thumbnailBytes)
		if err != nil {
			logger.Error(ctx, "Failed to save thumbnail to storage:", err)
			return nil, err
		}
		fileMetadata.CoverUrl = thumbnailURL
	}
	// 保存文件到存储服务
	fileURL, err := uc.store.SaveFile(context.Background(), folder, objectPath, contentType, fileBytes)
	if err != nil {
		logger.Error(ctx, "Failed to save file to storage:", err)
		return nil, err
	}
	fileMetadata.FileUrl = fileURL

	// 保存文件元数据到数据库
	_, err = models.NewFileMetadataModel().Create(fileMetadata)
	if err != nil {
		logger.Error(ctx, "Failed to save file metadata to database:", err)
		return nil, err
	}

	return
}
