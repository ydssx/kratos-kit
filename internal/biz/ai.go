package biz

import (
	"context"

	apiai "github.com/ydssx/kratos-kit/api/ai/v1"

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

type AiUseCase struct{}

func NewAiUseCase() *AiUseCase {
	return &AiUseCase{}
}

// Chat Chat 与AI助手对话
func (uc *AiUseCase) Chat(ctx context.Context, req *apiai.ChatRequest) (res *apiai.ChatResponse, err error) {
	res = new(apiai.ChatResponse)

	// TODO:ADD logic here and delete this line.

	return
}

// CreateConversation CreateConversation 创建新的对话
func (uc *AiUseCase) CreateConversation(ctx context.Context, req *apiai.CreateConversationRequest) (res *apiai.CreateConversationResponse, err error) {
	res = new(apiai.CreateConversationResponse)

	// TODO:ADD logic here and delete this line.

	return
}

// ListConversations ListConversations 获取对话列表
func (uc *AiUseCase) ListConversations(ctx context.Context, req *apiai.ListConversationsRequest) (res *apiai.ListConversationsResponse, err error) {
	res = new(apiai.ListConversationsResponse)

	// TODO:ADD logic here and delete this line.

	return
}

// GetConversation GetConversation 获取对话详情
func (uc *AiUseCase) GetConversation(ctx context.Context, req *apiai.GetConversationRequest) (res *apiai.GetConversationResponse, err error) {
	res = new(apiai.GetConversationResponse)

	// TODO:ADD logic here and delete this line.

	return
}

// DeleteConversation DeleteConversation 删除对话
func (uc *AiUseCase) DeleteConversation(ctx context.Context, req *apiai.DeleteConversationRequest) (res *emptypb.Empty, err error) {
	res = new(emptypb.Empty)

	// TODO:ADD logic here and delete this line.

	return
}

// UpdateConversation UpdateConversation 更新对话信息
func (uc *AiUseCase) UpdateConversation(ctx context.Context, req *apiai.UpdateConversationRequest) (res *emptypb.Empty, err error) {
	res = new(emptypb.Empty)

	// TODO:ADD logic here and delete this line.

	return
}

// GenerateImage GenerateImage 生成图片
func (uc *AiUseCase) GenerateImage(ctx context.Context, req *apiai.GenerateImageRequest) (res *apiai.GenerateImageResponse, err error) {
	res = new(apiai.GenerateImageResponse)

	// TODO:ADD logic here and delete this line.

	return
}

// EditImage EditImage 编辑/变体图片
func (uc *AiUseCase) EditImage(ctx context.Context, req *apiai.EditImageRequest) (res *apiai.GenerateImageResponse, err error) {
	res = new(apiai.GenerateImageResponse)

	// TODO:ADD logic here and delete this line.

	return
}

// ListGeneratedImages ListGeneratedImages 获取生成的图片列表
func (uc *AiUseCase) ListGeneratedImages(ctx context.Context, req *apiai.ListGeneratedImagesRequest) (res *apiai.ListGeneratedImagesResponse, err error) {
	res = new(apiai.ListGeneratedImagesResponse)

	// TODO:ADD logic here and delete this line.

	return
}
