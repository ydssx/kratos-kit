package biz

import (
	"context"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	apiai "github.com/ydssx/kratos-kit/api/ai/v1"
	"github.com/ydssx/kratos-kit/pkg/logger"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
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

type AiUseCase struct {
	llmModels map[string]llms.Model
	mu        sync.Mutex
	llmChains map[string]*chains.LLMChain
}

func NewAiUseCase() *AiUseCase {
	return &AiUseCase{llmModels: map[string]llms.Model{}, llmChains: map[string]*chains.LLMChain{}}
}

// Chat Chat 与AI助手对话
func (uc *AiUseCase) Chat(ctx *gin.Context, req *apiai.ChatRequest) (res *apiai.ChatResponse, err error) {
	// 设置SSE响应头
	ctx.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")

	// 创建一个channel用于检测客户端断开连接
	clientGone := ctx.Writer.CloseNotify()

	err = uc.initLLMModels(ctx, req.Model)
	if err != nil {
		return nil, err
	}

	// content := []llms.MessageContent{
	// 	llms.TextParts(llms.ChatMessageTypeHuman, req.Content),
	// }
	// memoryVariables, _ := uc.llmChains[req.Model].GetMemory().LoadMemoryVariables(ctx, map[string]any{})
	// logger.Infof(ctx, "memoryVariables: %v", memoryVariables)

	completion, err := chains.Run(ctx, uc.llmChains[req.Model], req.Content, chains.WithStreamingFunc(func(_ context.Context, chunk []byte) error {
		select {
		case <-clientGone:
			return fmt.Errorf("client disconnected")
		default:
			// 发送SSE格式的数据
			_, err := fmt.Fprintf(ctx.Writer, "data: %s\n\n", chunk)
			ctx.Writer.Flush()
			return err
		}
	}))
	if err != nil {
		return nil, err
	}

	res = &apiai.ChatResponse{
		Message: &apiai.Message{
			Content: completion,
		},
	}

	return res, nil
}

// initLLMModels 初始化LLM模型
func (uc *AiUseCase) initLLMModels(ctx context.Context, model string) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if _, ok := uc.llmModels[model]; ok {
		return nil
	}
	if _, ok := uc.llmChains[model]; ok {
		return nil
	}

	llm, err := ollama.New(ollama.WithModel(model), ollama.WithSystemPrompt("You are a helpful assistant and answer questions in Chinese."))
	if err != nil {
		logger.Errorf(ctx, "Failed to create LLM model: %v", err)
		return err
	}
	uc.llmModels[model] = llm

	// buffer := memory.NewConversationBuffer()
	prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate("Translate {{.input}} to Chinese. Only give the translation, no other information. If can't translate, just return the original text.", []string{"input"}),
	})
	llmchain := chains.NewLLMChain(uc.llmModels[model], prompt)
	uc.llmChains[model] = llmchain

	return nil
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
