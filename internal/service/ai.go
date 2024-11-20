package service

import (
	"context"

	"github.com/gin-gonic/gin"
	apiai "github.com/ydssx/kratos-kit/api/ai/v1"
	"github.com/ydssx/kratos-kit/internal/biz"
	"github.com/ydssx/kratos-kit/pkg/util"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	_ = context.Background
	_ = emptypb.Empty{}
	_ = timestamppb.Timestamp{}
	_ = durationpb.Duration{}
)

type AIService struct {
	uc *biz.AiUseCase

	apiai.UnimplementedAIServiceServer
}

func NewAIService(uc *biz.AiUseCase) *AIService {
	return &AIService{uc: uc}
}

// Chat Chat 与AI助手对话
func (s *AIService) Chat(ctx *gin.Context) {
	req := new(apiai.ChatRequest)
	// if err := ctx.ShouldBindJSON(req); err != nil {
	// 	util.FailWithError(ctx, err)
	// 	return
	// }
	req.Model = "llama3.2:3b"
	req.Content = ctx.Query("message")
	_, err := s.uc.Chat(ctx, req)
	if err != nil {
		util.FailWithError(ctx, err)
		return
	}
	// util.OKWithData(ctx, res)
}

// CreateConversation CreateConversation 创建新的对话
func (s *AIService) CreateConversation(ctx context.Context, req *apiai.CreateConversationRequest) (res *apiai.CreateConversationResponse, err error) {
	return s.uc.CreateConversation(ctx, req)
}

// ListConversations ListConversations 获取对话列表
func (s *AIService) ListConversations(ctx context.Context, req *apiai.ListConversationsRequest) (res *apiai.ListConversationsResponse, err error) {
	return s.uc.ListConversations(ctx, req)
}

// GetConversation GetConversation 获取对话详情
func (s *AIService) GetConversation(ctx context.Context, req *apiai.GetConversationRequest) (res *apiai.GetConversationResponse, err error) {
	return s.uc.GetConversation(ctx, req)
}

// DeleteConversation DeleteConversation 删除对话
func (s *AIService) DeleteConversation(ctx context.Context, req *apiai.DeleteConversationRequest) (res *emptypb.Empty, err error) {
	return s.uc.DeleteConversation(ctx, req)
}

// UpdateConversation UpdateConversation 更新对话信息
func (s *AIService) UpdateConversation(ctx context.Context, req *apiai.UpdateConversationRequest) (res *emptypb.Empty, err error) {
	return s.uc.UpdateConversation(ctx, req)
}

// GenerateImage GenerateImage 生成图片
func (s *AIService) GenerateImage(ctx context.Context, req *apiai.GenerateImageRequest) (res *apiai.GenerateImageResponse, err error) {
	return s.uc.GenerateImage(ctx, req)
}

// EditImage EditImage 编辑/变体图片
func (s *AIService) EditImage(ctx context.Context, req *apiai.EditImageRequest) (res *apiai.GenerateImageResponse, err error) {
	return s.uc.EditImage(ctx, req)
}

// ListGeneratedImages ListGeneratedImages 获取生成的图片列表
func (s *AIService) ListGeneratedImages(ctx context.Context, req *apiai.ListGeneratedImagesRequest) (res *apiai.ListGeneratedImagesResponse, err error) {
	return s.uc.ListGeneratedImages(ctx, req)
}
