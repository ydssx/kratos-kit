package biz

import (
    "context"
    "time"
)

// Message 消息模型
type Message struct {
    ID        uint      `json:"id"`
    Role      string    `json:"role"`      // user 或 assistant
    Content   string    `json:"content"`   // 消息内容
    CreatedAt time.Time `json:"createdAt"` 
}

// Conversation 对话模型
type Conversation struct {
    ID        uint      `json:"id"`
    Title     string    `json:"title"`     // 对话标题
    UserID    uint      `json:"userId"`    // 所属用户
    Model     string    `json:"model"`     // 使用的模型
    Messages  []Message `json:"messages"`   // 消息列表
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

// ChatService 聊天服务接口
type ChatService interface {
    CreateConversation(ctx context.Context, userID uint, model string) (*Conversation, error)
    GetConversation(ctx context.Context, id uint) (*Conversation, error)
    ListConversations(ctx context.Context, userID uint) ([]*Conversation, error)
    DeleteConversation(ctx context.Context, id uint) error
    SendMessage(ctx context.Context, convID uint, content string) (*Message, error)
}

// chatService 聊天服务实现
type chatService struct {
    repo ChatRepository
    ai   AIService
}

// NewChatService 创建聊天服务
func NewChatService(repo ChatRepository, ai AIService) *chatService {
    return &chatService{
        repo: repo,
        ai:   ai,
    }
}

// ChatRepository 聊天数据仓库接口
type ChatRepository interface {
    CreateConversation(ctx context.Context, conv *Conversation) error
    GetConversation(ctx context.Context, id uint) (*Conversation, error)
    ListConversations(ctx context.Context, userID uint) ([]*Conversation, error)
    DeleteConversation(ctx context.Context, id uint) error
    AddMessage(ctx context.Context, convID uint, msg *Message) error
}

// AIService AI服务接口
type AIService interface {
    GenerateResponse(ctx context.Context, messages []Message, model string) (string, error)
}

// 实现 ChatService 接口方法
func (s *chatService) CreateConversation(ctx context.Context, userID uint, model string) (*Conversation, error) {
    conv := &Conversation{
        UserID:    userID,
        Model:     model,
        Title:     "新对话",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    if err := s.repo.CreateConversation(ctx, conv); err != nil {
        return nil, err
    }
    
    return conv, nil
}

func (s *chatService) SendMessage(ctx context.Context, convID uint, content string) (*Message, error) {
    // 获取对话
    conv, err := s.repo.GetConversation(ctx, convID)
    if err != nil {
        return nil, err
    }

    // 创建用户消息
    userMsg := &Message{
        Role:      "user",
        Content:   content,
        CreatedAt: time.Now(),
    }
    
    // 保存用户消息
    if err := s.repo.AddMessage(ctx, convID, userMsg); err != nil {
        return nil, err
    }

    // 生成 AI 回复
    aiResponse, err := s.ai.GenerateResponse(ctx, append(conv.Messages, *userMsg), conv.Model)
    if err != nil {
        return nil, err
    }

    // 创建 AI 消息
    aiMsg := &Message{
        Role:      "assistant",
        Content:   aiResponse,
        CreatedAt: time.Now(),
    }
    
    // 保存 AI 消息
    if err := s.repo.AddMessage(ctx, convID, aiMsg); err != nil {
        return nil, err
    }

    return aiMsg, nil
}

// ... 实现其他接口方法 