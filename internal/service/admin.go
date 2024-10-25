package service

import (
	"context"

	adminv1 "github.com/ydssx/kratos-kit/api/admin/v1"
	"github.com/ydssx/kratos-kit/internal/biz"
	"github.com/ydssx/kratos-kit/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
)

var _ = context.Background

type AdminService struct {
	uc *biz.AdminUseCase

	adminv1.UnimplementedAdminServiceServer
}

func NewAdminService(uc *biz.AdminUseCase) *AdminService {
	return &AdminService{uc: uc}
}

func (s *AdminService) Upload(c *gin.Context) {
	userID := c.GetInt("user_id") // 从gin上下文中获取用户ID

	file, err := c.FormFile("file")
	if err != nil {
		util.FailWithError(c, err)
		return
	}

	// 调用业务逻辑处理文件上传
	data, err := s.uc.UploadFile(c, userID, file)
	if err != nil {
		if util.IsImageFile(file.Filename) {
			err = errors.New(util.ERROR, err.Error(), "The image is invalid. Please upload again.")
		}
		util.FailWithError(c, err)
		return
	}

	// 返回文件URL
	util.OKWithData(c, data)
}
