// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: api/ai/v1/ai.proto

package aiv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	AI_Chat_FullMethodName                = "/api.ai.AI/Chat"
	AI_CreateConversation_FullMethodName  = "/api.ai.AI/CreateConversation"
	AI_ListConversations_FullMethodName   = "/api.ai.AI/ListConversations"
	AI_GetConversation_FullMethodName     = "/api.ai.AI/GetConversation"
	AI_DeleteConversation_FullMethodName  = "/api.ai.AI/DeleteConversation"
	AI_UpdateConversation_FullMethodName  = "/api.ai.AI/UpdateConversation"
	AI_GenerateImage_FullMethodName       = "/api.ai.AI/GenerateImage"
	AI_EditImage_FullMethodName           = "/api.ai.AI/EditImage"
	AI_ListGeneratedImages_FullMethodName = "/api.ai.AI/ListGeneratedImages"
)

// AIClient is the client API for AI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// AI服务接口定义
type AIClient interface {
	// Chat 与AI助手对话
	Chat(ctx context.Context, in *ChatRequest, opts ...grpc.CallOption) (*ChatResponse, error)
	// CreateConversation 创建新的对话
	CreateConversation(ctx context.Context, in *CreateConversationRequest, opts ...grpc.CallOption) (*CreateConversationResponse, error)
	// ListConversations 获取对话列表
	ListConversations(ctx context.Context, in *ListConversationsRequest, opts ...grpc.CallOption) (*ListConversationsResponse, error)
	// GetConversation 获取对话详情
	GetConversation(ctx context.Context, in *GetConversationRequest, opts ...grpc.CallOption) (*GetConversationResponse, error)
	// DeleteConversation 删除对话
	DeleteConversation(ctx context.Context, in *DeleteConversationRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// UpdateConversation 更新对话信息
	UpdateConversation(ctx context.Context, in *UpdateConversationRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// GenerateImage 生成图片
	GenerateImage(ctx context.Context, in *GenerateImageRequest, opts ...grpc.CallOption) (*GenerateImageResponse, error)
	// EditImage 编辑/变体图片
	EditImage(ctx context.Context, in *EditImageRequest, opts ...grpc.CallOption) (*GenerateImageResponse, error)
	// ListGeneratedImages 获取生成的图片列表
	ListGeneratedImages(ctx context.Context, in *ListGeneratedImagesRequest, opts ...grpc.CallOption) (*ListGeneratedImagesResponse, error)
}

type aIClient struct {
	cc grpc.ClientConnInterface
}

func NewAIClient(cc grpc.ClientConnInterface) AIClient {
	return &aIClient{cc}
}

func (c *aIClient) Chat(ctx context.Context, in *ChatRequest, opts ...grpc.CallOption) (*ChatResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ChatResponse)
	err := c.cc.Invoke(ctx, AI_Chat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aIClient) CreateConversation(ctx context.Context, in *CreateConversationRequest, opts ...grpc.CallOption) (*CreateConversationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateConversationResponse)
	err := c.cc.Invoke(ctx, AI_CreateConversation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aIClient) ListConversations(ctx context.Context, in *ListConversationsRequest, opts ...grpc.CallOption) (*ListConversationsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListConversationsResponse)
	err := c.cc.Invoke(ctx, AI_ListConversations_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aIClient) GetConversation(ctx context.Context, in *GetConversationRequest, opts ...grpc.CallOption) (*GetConversationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetConversationResponse)
	err := c.cc.Invoke(ctx, AI_GetConversation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aIClient) DeleteConversation(ctx context.Context, in *DeleteConversationRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AI_DeleteConversation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aIClient) UpdateConversation(ctx context.Context, in *UpdateConversationRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AI_UpdateConversation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aIClient) GenerateImage(ctx context.Context, in *GenerateImageRequest, opts ...grpc.CallOption) (*GenerateImageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GenerateImageResponse)
	err := c.cc.Invoke(ctx, AI_GenerateImage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aIClient) EditImage(ctx context.Context, in *EditImageRequest, opts ...grpc.CallOption) (*GenerateImageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GenerateImageResponse)
	err := c.cc.Invoke(ctx, AI_EditImage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aIClient) ListGeneratedImages(ctx context.Context, in *ListGeneratedImagesRequest, opts ...grpc.CallOption) (*ListGeneratedImagesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListGeneratedImagesResponse)
	err := c.cc.Invoke(ctx, AI_ListGeneratedImages_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AIServer is the server API for AI service.
// All implementations should embed UnimplementedAIServer
// for forward compatibility.
//
// AI服务接口定义
type AIServer interface {
	// Chat 与AI助手对话
	Chat(context.Context, *ChatRequest) (*ChatResponse, error)
	// CreateConversation 创建新的对话
	CreateConversation(context.Context, *CreateConversationRequest) (*CreateConversationResponse, error)
	// ListConversations 获取对话列表
	ListConversations(context.Context, *ListConversationsRequest) (*ListConversationsResponse, error)
	// GetConversation 获取对话详情
	GetConversation(context.Context, *GetConversationRequest) (*GetConversationResponse, error)
	// DeleteConversation 删除对话
	DeleteConversation(context.Context, *DeleteConversationRequest) (*emptypb.Empty, error)
	// UpdateConversation 更新对话信息
	UpdateConversation(context.Context, *UpdateConversationRequest) (*emptypb.Empty, error)
	// GenerateImage 生成图片
	GenerateImage(context.Context, *GenerateImageRequest) (*GenerateImageResponse, error)
	// EditImage 编辑/变体图片
	EditImage(context.Context, *EditImageRequest) (*GenerateImageResponse, error)
	// ListGeneratedImages 获取生成的图片列表
	ListGeneratedImages(context.Context, *ListGeneratedImagesRequest) (*ListGeneratedImagesResponse, error)
}

// UnimplementedAIServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAIServer struct{}

func (UnimplementedAIServer) Chat(context.Context, *ChatRequest) (*ChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Chat not implemented")
}
func (UnimplementedAIServer) CreateConversation(context.Context, *CreateConversationRequest) (*CreateConversationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateConversation not implemented")
}
func (UnimplementedAIServer) ListConversations(context.Context, *ListConversationsRequest) (*ListConversationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListConversations not implemented")
}
func (UnimplementedAIServer) GetConversation(context.Context, *GetConversationRequest) (*GetConversationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConversation not implemented")
}
func (UnimplementedAIServer) DeleteConversation(context.Context, *DeleteConversationRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteConversation not implemented")
}
func (UnimplementedAIServer) UpdateConversation(context.Context, *UpdateConversationRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateConversation not implemented")
}
func (UnimplementedAIServer) GenerateImage(context.Context, *GenerateImageRequest) (*GenerateImageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateImage not implemented")
}
func (UnimplementedAIServer) EditImage(context.Context, *EditImageRequest) (*GenerateImageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditImage not implemented")
}
func (UnimplementedAIServer) ListGeneratedImages(context.Context, *ListGeneratedImagesRequest) (*ListGeneratedImagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListGeneratedImages not implemented")
}
func (UnimplementedAIServer) testEmbeddedByValue() {}

// UnsafeAIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AIServer will
// result in compilation errors.
type UnsafeAIServer interface {
	mustEmbedUnimplementedAIServer()
}

func RegisterAIServer(s grpc.ServiceRegistrar, srv AIServer) {
	// If the following call pancis, it indicates UnimplementedAIServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AI_ServiceDesc, srv)
}

func _AI_Chat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AIServer).Chat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AI_Chat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AIServer).Chat(ctx, req.(*ChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AI_CreateConversation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateConversationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AIServer).CreateConversation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AI_CreateConversation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AIServer).CreateConversation(ctx, req.(*CreateConversationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AI_ListConversations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListConversationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AIServer).ListConversations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AI_ListConversations_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AIServer).ListConversations(ctx, req.(*ListConversationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AI_GetConversation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConversationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AIServer).GetConversation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AI_GetConversation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AIServer).GetConversation(ctx, req.(*GetConversationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AI_DeleteConversation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteConversationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AIServer).DeleteConversation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AI_DeleteConversation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AIServer).DeleteConversation(ctx, req.(*DeleteConversationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AI_UpdateConversation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateConversationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AIServer).UpdateConversation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AI_UpdateConversation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AIServer).UpdateConversation(ctx, req.(*UpdateConversationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AI_GenerateImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AIServer).GenerateImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AI_GenerateImage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AIServer).GenerateImage(ctx, req.(*GenerateImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AI_EditImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AIServer).EditImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AI_EditImage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AIServer).EditImage(ctx, req.(*EditImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AI_ListGeneratedImages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListGeneratedImagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AIServer).ListGeneratedImages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AI_ListGeneratedImages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AIServer).ListGeneratedImages(ctx, req.(*ListGeneratedImagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AI_ServiceDesc is the grpc.ServiceDesc for AI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.ai.AI",
	HandlerType: (*AIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Chat",
			Handler:    _AI_Chat_Handler,
		},
		{
			MethodName: "CreateConversation",
			Handler:    _AI_CreateConversation_Handler,
		},
		{
			MethodName: "ListConversations",
			Handler:    _AI_ListConversations_Handler,
		},
		{
			MethodName: "GetConversation",
			Handler:    _AI_GetConversation_Handler,
		},
		{
			MethodName: "DeleteConversation",
			Handler:    _AI_DeleteConversation_Handler,
		},
		{
			MethodName: "UpdateConversation",
			Handler:    _AI_UpdateConversation_Handler,
		},
		{
			MethodName: "GenerateImage",
			Handler:    _AI_GenerateImage_Handler,
		},
		{
			MethodName: "EditImage",
			Handler:    _AI_EditImage_Handler,
		},
		{
			MethodName: "ListGeneratedImages",
			Handler:    _AI_ListGeneratedImages_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/ai/v1/ai.proto",
}