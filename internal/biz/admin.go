package biz

import (
	"context"
	"mime/multipart"

	"github.com/ydssx/kratos-kit/models"
	"github.com/ydssx/kratos-kit/pkg/errors"

	"github.com/gin-gonic/gin"
)

type AdminUseCase struct {
	tx       Transaction
	commonUc *CommonUseCase
	userRepo UserRepo
}

type adminUseCaseKey struct{}

func AdminUseCaseFromContext(ctx context.Context) *AdminUseCase {
	return ctx.Value(adminUseCaseKey{}).(*AdminUseCase)
}

func WithAdminUseCase(ctx context.Context, uc *AdminUseCase) context.Context {
	return context.WithValue(ctx, adminUseCaseKey{}, uc)
}

func NewAdminUseCase(commonUc *CommonUseCase, userRepo UserRepo, tx Transaction) *AdminUseCase {
	return &AdminUseCase{commonUc: commonUc, userRepo: userRepo, tx: tx}
}

func (uc *AdminUseCase) UploadFile(ctx *gin.Context, userID int, file *multipart.FileHeader) (result *UploadResult, err error) {
	result = new(UploadResult)
	userType := int(models.UserTypeAdmin)
	fileInfo, err := uc.commonUc.UploadFile(context.Background(), userID, userType, file)
	if err != nil {
		return nil, errors.Wrap(err, "上传文件失败")
	}

	result.FileId = int(fileInfo.ID)
	result.FileUrl = fileInfo.FileUrl
	result.ThumbnailURL = fileInfo.CoverUrl
	result.FileName = file.Filename

	return
}
