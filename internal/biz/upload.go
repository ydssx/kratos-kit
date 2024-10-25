package biz

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"strings"

	// "math"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/ydssx/kratos-kit/common/conf"
	"github.com/ydssx/kratos-kit/models"
	"github.com/ydssx/kratos-kit/pkg/errors"
	"github.com/ydssx/kratos-kit/pkg/logger"
	"github.com/ydssx/kratos-kit/pkg/storage"
	"github.com/ydssx/kratos-kit/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

type UploadUseCase struct {
	store    storage.Storage // 存储接口
	c        *conf.Bootstrap
	commonUc *CommonUseCase
}

func NewUploadUseCase(store storage.Storage, conf *conf.Bootstrap, commonUc *CommonUseCase) *UploadUseCase {
	return &UploadUseCase{store: store, c: conf, commonUc: commonUc}
}

// UploadFile 上传文件
func (uc *UploadUseCase) UploadFile(c *gin.Context, userID int, file *multipart.FileHeader) (data interface{}, errInfo error) {
	result := new(UploadResult)

	userType := c.GetInt("user_type")

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("Failed to open file: %v", err)
	}
	defer src.Close()

	// 读取文件内容
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file: %v", err)
	}

	if len(fileBytes) > 8*1024*1024 {
		return nil, errors.NewUserError("file size exceeds 8MB")
	}

	file_md5 := util.MD5Bytes(fileBytes)

	fileInfo, err := models.NewFileMetadataModel().SetMd5(file_md5).SetUserId(int64(userID)).FirstOne()
	if err == nil {
		result.FileId = int(fileInfo.ID)
		result.FileUrl = fileInfo.FileUrl
		result.ThumbnailURL = fileInfo.CoverUrl
		return result, nil
	}

	// 构建文件名
	originalName := filepath.Base(file.Filename)
	ext := filepath.Ext(originalName)
	storedName := time.Now().Format("20060102150405") + util.GenerateCode(2) + ext

	if !util.IsImageFile(originalName) {
		return nil, errors.NewUserError("file type is not supported")
	}

	// TODO: 本地存储路径需要配置
	// 创建本地存储路径
	uploadDir := "C:\\uploads\\"
	if runtime.GOOS == "linux" {
		uploadDir = "/tmp/uploads/"
	}
	filePath := filepath.Join(uploadDir, storedName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Error("Failed to create upload directory:", err)
		return nil, err
	}
	defer os.RemoveAll(uploadDir) // 删除本地存储路径

	saveToStorage := true
	defer func() {
		if saveToStorage {
			// 保存文件
			fileInfo, err := uc.commonUc.UploadFile(context.Background(), userID, userType, file)
			if err != nil {
				log.Error("Failed to save file to storage:", err)
				errInfo = err
				return
			}

			result.FileId = int(fileInfo.ID)
			result.FileUrl = fileInfo.FileUrl
			result.ThumbnailURL = fileInfo.CoverUrl
		}
	}()

	if util.IsVideoFile(originalName) {
		// 获取视频时长
		duration, err := util.GetVideoDuration(filePath)
		if err != nil {
			log.Error("Failed to get video duration:", err)
			return nil, err
		}

		if duration > float64(time.Minute*10) {
			saveToStorage = false
			return nil, errors.NewUserError("video length exceeds 10 minutes")
		}
	}
	return result, errInfo
}

// 清理用户上传文件
func (uc *UploadUseCase) CleanUploadFile(ctx context.Context) error {
	users, err := models.NewUserModel().SetUserType(models.UserTypeNormal).PluckIds()
	if err != nil {
		return errors.Wrapf(err, "Failed to get user ids")
	}
	files := models.NewFileMetadataModel().SetUserId(users...).CreatedAtLT(time.Now().AddDate(0, -1, 0)).List()
	for _, file := range files {
		if file.FileType == models.FileTypeImage {
			continue
		}
		err = uc.store.DeleteFile(ctx, file.Filename)
		if err != nil {
			logger.Errorf(ctx, "Failed to delete file %s: %s", file.Filename, err.Error())
		}
		if file.CoverUrl != "" {
			err = uc.store.DeleteFile(ctx, strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))+"_thumbnail.jpg")
			if err != nil {
				logger.Errorf(ctx, "Failed to delete file %s: %s", file.CoverUrl, err.Error())
			}
		}
		if err == nil || strings.Contains(err.Error(), "object doesn't exist") {
			models.NewFileMetadataModel().SetIds(int64(file.ID)).Delete()
		}
	}
	return nil
}
