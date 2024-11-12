// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: api/job/v1/job.proto

package jobv1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	_ "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 任务类型定义
type JobType int32

const (
	JobType_TEST_JOB               JobType = 0
	JobType_TEST_CRON_JOB          JobType = 1
	JobType_GOOGLE_INSTANCE_ADJUST JobType = 2
	// 订阅自动续费
	JobType_SUBSCRIPTION_RENEWAL JobType = 3
	// 任务超时积分退还
	JobType_TASK_TIMEOUT_REFUND JobType = 4
	// 更新视频排序
	JobType_UPDATE_VIDEO_SORT JobType = 5
	// 清理用户上传文件
	JobType_CLEAN_USER_UPLOAD_FILES JobType = 6
	// 清理旧日志文件
	JobType_CLEAN_OLD_LOG_FILES JobType = 7
	// 清理用户过期积分
	JobType_CLEAN_USER_EXPIRED_POINTS JobType = 8
	// 重置每日字符数
	JobType_RESET_DAILY_CHARACTERS JobType = 9
)

// Enum value maps for JobType.
var (
	JobType_name = map[int32]string{
		0: "TEST_JOB",
		1: "TEST_CRON_JOB",
		2: "GOOGLE_INSTANCE_ADJUST",
		3: "SUBSCRIPTION_RENEWAL",
		4: "TASK_TIMEOUT_REFUND",
		5: "UPDATE_VIDEO_SORT",
		6: "CLEAN_USER_UPLOAD_FILES",
		7: "CLEAN_OLD_LOG_FILES",
		8: "CLEAN_USER_EXPIRED_POINTS",
		9: "RESET_DAILY_CHARACTERS",
	}
	JobType_value = map[string]int32{
		"TEST_JOB":                  0,
		"TEST_CRON_JOB":             1,
		"GOOGLE_INSTANCE_ADJUST":    2,
		"SUBSCRIPTION_RENEWAL":      3,
		"TASK_TIMEOUT_REFUND":       4,
		"UPDATE_VIDEO_SORT":         5,
		"CLEAN_USER_UPLOAD_FILES":   6,
		"CLEAN_OLD_LOG_FILES":       7,
		"CLEAN_USER_EXPIRED_POINTS": 8,
		"RESET_DAILY_CHARACTERS":    9,
	}
)

func (x JobType) Enum() *JobType {
	p := new(JobType)
	*p = x
	return p
}

func (x JobType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (JobType) Descriptor() protoreflect.EnumDescriptor {
	return file_api_job_v1_job_proto_enumTypes[0].Descriptor()
}

func (JobType) Type() protoreflect.EnumType {
	return &file_api_job_v1_job_proto_enumTypes[0]
}

func (x JobType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use JobType.Descriptor instead.
func (JobType) EnumDescriptor() ([]byte, []int) {
	return file_api_job_v1_job_proto_rawDescGZIP(), []int{0}
}

type AdminJob int32

const (
	// 生成每日报表
	AdminJob_GENERATE_DAILY_REPORT AdminJob = 0
	// 生成当日报表
	AdminJob_GENERATE_TODAY_REPORT AdminJob = 1
)

// Enum value maps for AdminJob.
var (
	AdminJob_name = map[int32]string{
		0: "GENERATE_DAILY_REPORT",
		1: "GENERATE_TODAY_REPORT",
	}
	AdminJob_value = map[string]int32{
		"GENERATE_DAILY_REPORT": 0,
		"GENERATE_TODAY_REPORT": 1,
	}
)

func (x AdminJob) Enum() *AdminJob {
	p := new(AdminJob)
	*p = x
	return p
}

func (x AdminJob) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AdminJob) Descriptor() protoreflect.EnumDescriptor {
	return file_api_job_v1_job_proto_enumTypes[1].Descriptor()
}

func (AdminJob) Type() protoreflect.EnumType {
	return &file_api_job_v1_job_proto_enumTypes[1]
}

func (x AdminJob) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AdminJob.Descriptor instead.
func (AdminJob) EnumDescriptor() ([]byte, []int) {
	return file_api_job_v1_job_proto_rawDescGZIP(), []int{1}
}

type EnqueueRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JobType   JobType                `protobuf:"varint,1,opt,name=job_type,json=jobType,proto3,enum=job.v1.JobType" json:"job_type,omitempty"` // 任务类型
	Payload   []byte                 `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`                                     // 任务参数
	RetryTime int64                  `protobuf:"varint,3,opt,name=retry_time,json=retryTime,proto3" json:"retry_time,omitempty"`               // 重试次数
	ProcessAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=process_at,json=processAt,proto3" json:"process_at,omitempty"`                // 任务执行时间
	ProcessIn *durationpb.Duration   `protobuf:"bytes,5,opt,name=process_in,json=processIn,proto3" json:"process_in,omitempty"`
	Deadline  *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=deadline,proto3" json:"deadline,omitempty"`
	// 任务完成后保留时间
	Retention *durationpb.Duration `protobuf:"bytes,7,opt,name=retention,proto3" json:"retention,omitempty"`
}

func (x *EnqueueRequest) Reset() {
	*x = EnqueueRequest{}
	mi := &file_api_job_v1_job_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EnqueueRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnqueueRequest) ProtoMessage() {}

func (x *EnqueueRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1_job_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnqueueRequest.ProtoReflect.Descriptor instead.
func (*EnqueueRequest) Descriptor() ([]byte, []int) {
	return file_api_job_v1_job_proto_rawDescGZIP(), []int{0}
}

func (x *EnqueueRequest) GetJobType() JobType {
	if x != nil {
		return x.JobType
	}
	return JobType_TEST_JOB
}

func (x *EnqueueRequest) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *EnqueueRequest) GetRetryTime() int64 {
	if x != nil {
		return x.RetryTime
	}
	return 0
}

func (x *EnqueueRequest) GetProcessAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ProcessAt
	}
	return nil
}

func (x *EnqueueRequest) GetProcessIn() *durationpb.Duration {
	if x != nil {
		return x.ProcessIn
	}
	return nil
}

func (x *EnqueueRequest) GetDeadline() *timestamppb.Timestamp {
	if x != nil {
		return x.Deadline
	}
	return nil
}

func (x *EnqueueRequest) GetRetention() *durationpb.Duration {
	if x != nil {
		return x.Retention
	}
	return nil
}

type EnqueueResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskId string `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
}

func (x *EnqueueResponse) Reset() {
	*x = EnqueueResponse{}
	mi := &file_api_job_v1_job_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EnqueueResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnqueueResponse) ProtoMessage() {}

func (x *EnqueueResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1_job_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnqueueResponse.ProtoReflect.Descriptor instead.
func (*EnqueueResponse) Descriptor() ([]byte, []int) {
	return file_api_job_v1_job_proto_rawDescGZIP(), []int{1}
}

func (x *EnqueueResponse) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

type QueryTasksRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskIds []string `protobuf:"bytes,1,rep,name=task_ids,json=taskIds,proto3" json:"task_ids,omitempty"`
}

func (x *QueryTasksRequest) Reset() {
	*x = QueryTasksRequest{}
	mi := &file_api_job_v1_job_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryTasksRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryTasksRequest) ProtoMessage() {}

func (x *QueryTasksRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1_job_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryTasksRequest.ProtoReflect.Descriptor instead.
func (*QueryTasksRequest) Descriptor() ([]byte, []int) {
	return file_api_job_v1_job_proto_rawDescGZIP(), []int{2}
}

func (x *QueryTasksRequest) GetTaskIds() []string {
	if x != nil {
		return x.TaskIds
	}
	return nil
}

type QueryTasksResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tasks []*QueryTasksResponse_TaskInfo `protobuf:"bytes,1,rep,name=tasks,proto3" json:"tasks,omitempty"`
}

func (x *QueryTasksResponse) Reset() {
	*x = QueryTasksResponse{}
	mi := &file_api_job_v1_job_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryTasksResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryTasksResponse) ProtoMessage() {}

func (x *QueryTasksResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1_job_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryTasksResponse.ProtoReflect.Descriptor instead.
func (*QueryTasksResponse) Descriptor() ([]byte, []int) {
	return file_api_job_v1_job_proto_rawDescGZIP(), []int{3}
}

func (x *QueryTasksResponse) GetTasks() []*QueryTasksResponse_TaskInfo {
	if x != nil {
		return x.Tasks
	}
	return nil
}

// 任务超时积分退还
type PayLoadTaskTimeout struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskId int32 `protobuf:"varint,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
}

func (x *PayLoadTaskTimeout) Reset() {
	*x = PayLoadTaskTimeout{}
	mi := &file_api_job_v1_job_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PayLoadTaskTimeout) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PayLoadTaskTimeout) ProtoMessage() {}

func (x *PayLoadTaskTimeout) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1_job_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PayLoadTaskTimeout.ProtoReflect.Descriptor instead.
func (*PayLoadTaskTimeout) Descriptor() ([]byte, []int) {
	return file_api_job_v1_job_proto_rawDescGZIP(), []int{4}
}

func (x *PayLoadTaskTimeout) GetTaskId() int32 {
	if x != nil {
		return x.TaskId
	}
	return 0
}

type PayLoadTest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *PayLoadTest) Reset() {
	*x = PayLoadTest{}
	mi := &file_api_job_v1_job_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PayLoadTest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PayLoadTest) ProtoMessage() {}

func (x *PayLoadTest) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1_job_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PayLoadTest.ProtoReflect.Descriptor instead.
func (*PayLoadTest) Descriptor() ([]byte, []int) {
	return file_api_job_v1_job_proto_rawDescGZIP(), []int{5}
}

func (x *PayLoadTest) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type PayLoadOrderPaymentCompleted struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId int64 `protobuf:"varint,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
}

func (x *PayLoadOrderPaymentCompleted) Reset() {
	*x = PayLoadOrderPaymentCompleted{}
	mi := &file_api_job_v1_job_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PayLoadOrderPaymentCompleted) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PayLoadOrderPaymentCompleted) ProtoMessage() {}

func (x *PayLoadOrderPaymentCompleted) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1_job_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PayLoadOrderPaymentCompleted.ProtoReflect.Descriptor instead.
func (*PayLoadOrderPaymentCompleted) Descriptor() ([]byte, []int) {
	return file_api_job_v1_job_proto_rawDescGZIP(), []int{6}
}

func (x *PayLoadOrderPaymentCompleted) GetOrderId() int64 {
	if x != nil {
		return x.OrderId
	}
	return 0
}

type PayLoadOrderTimeout struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderNum string `protobuf:"bytes,1,opt,name=order_num,json=orderNum,proto3" json:"order_num,omitempty"`
}

func (x *PayLoadOrderTimeout) Reset() {
	*x = PayLoadOrderTimeout{}
	mi := &file_api_job_v1_job_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PayLoadOrderTimeout) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PayLoadOrderTimeout) ProtoMessage() {}

func (x *PayLoadOrderTimeout) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1_job_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PayLoadOrderTimeout.ProtoReflect.Descriptor instead.
func (*PayLoadOrderTimeout) Descriptor() ([]byte, []int) {
	return file_api_job_v1_job_proto_rawDescGZIP(), []int{7}
}

func (x *PayLoadOrderTimeout) GetOrderNum() string {
	if x != nil {
		return x.OrderNum
	}
	return ""
}

type QueuingTimeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskId int64 `protobuf:"varint,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
}

func (x *QueuingTimeRequest) Reset() {
	*x = QueuingTimeRequest{}
	mi := &file_api_job_v1_job_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueuingTimeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueuingTimeRequest) ProtoMessage() {}

func (x *QueuingTimeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1_job_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueuingTimeRequest.ProtoReflect.Descriptor instead.
func (*QueuingTimeRequest) Descriptor() ([]byte, []int) {
	return file_api_job_v1_job_proto_rawDescGZIP(), []int{8}
}

func (x *QueuingTimeRequest) GetTaskId() int64 {
	if x != nil {
		return x.TaskId
	}
	return 0
}

type QueuingTimeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskId  int64 `protobuf:"varint,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	Seconds int64 `protobuf:"varint,2,opt,name=seconds,proto3" json:"seconds,omitempty"`
}

func (x *QueuingTimeResponse) Reset() {
	*x = QueuingTimeResponse{}
	mi := &file_api_job_v1_job_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueuingTimeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueuingTimeResponse) ProtoMessage() {}

func (x *QueuingTimeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1_job_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueuingTimeResponse.ProtoReflect.Descriptor instead.
func (*QueuingTimeResponse) Descriptor() ([]byte, []int) {
	return file_api_job_v1_job_proto_rawDescGZIP(), []int{9}
}

func (x *QueuingTimeResponse) GetTaskId() int64 {
	if x != nil {
		return x.TaskId
	}
	return 0
}

func (x *QueuingTimeResponse) GetSeconds() int64 {
	if x != nil {
		return x.Seconds
	}
	return 0
}

type QueryTasksResponse_TaskInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskId string `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	Result []byte `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
	Status string `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *QueryTasksResponse_TaskInfo) Reset() {
	*x = QueryTasksResponse_TaskInfo{}
	mi := &file_api_job_v1_job_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryTasksResponse_TaskInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryTasksResponse_TaskInfo) ProtoMessage() {}

func (x *QueryTasksResponse_TaskInfo) ProtoReflect() protoreflect.Message {
	mi := &file_api_job_v1_job_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryTasksResponse_TaskInfo.ProtoReflect.Descriptor instead.
func (*QueryTasksResponse_TaskInfo) Descriptor() ([]byte, []int) {
	return file_api_job_v1_job_proto_rawDescGZIP(), []int{3, 0}
}

func (x *QueryTasksResponse_TaskInfo) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

func (x *QueryTasksResponse_TaskInfo) GetResult() []byte {
	if x != nil {
		return x.Result
	}
	return nil
}

func (x *QueryTasksResponse_TaskInfo) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_api_job_v1_job_proto protoreflect.FileDescriptor

var file_api_job_v1_job_proto_rawDesc = []byte{
	0x0a, 0x14, 0x61, 0x70, 0x69, 0x2f, 0x6a, 0x6f, 0x62, 0x2f, 0x76, 0x31, 0x2f, 0x6a, 0x6f, 0x62,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6a, 0x6f, 0x62, 0x2e, 0x76, 0x31, 0x1a, 0x1c,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xdb, 0x02, 0x0a, 0x0e, 0x45,
	0x6e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a,
	0x08, 0x6a, 0x6f, 0x62, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x0f, 0x2e, 0x6a, 0x6f, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x4a, 0x6f, 0x62, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x07, 0x6a, 0x6f, 0x62, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c,
	0x6f, 0x61, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65, 0x74, 0x72, 0x79, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x72, 0x65, 0x74, 0x72, 0x79, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x61, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x41, 0x74, 0x12, 0x38, 0x0a,
	0x0a, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x69, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x70, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x6e, 0x12, 0x36, 0x0a, 0x08, 0x64, 0x65, 0x61, 0x64, 0x6c,
	0x69, 0x6e, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08, 0x64, 0x65, 0x61, 0x64, 0x6c, 0x69, 0x6e, 0x65, 0x12,
	0x37, 0x0a, 0x09, 0x72, 0x65, 0x74, 0x65, 0x6e, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x72,
	0x65, 0x74, 0x65, 0x6e, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x2a, 0x0a, 0x0f, 0x45, 0x6e, 0x71, 0x75,
	0x65, 0x75, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x74,
	0x61, 0x73, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61,
	0x73, 0x6b, 0x49, 0x64, 0x22, 0x2e, 0x0a, 0x11, 0x51, 0x75, 0x65, 0x72, 0x79, 0x54, 0x61, 0x73,
	0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x61, 0x73,
	0x6b, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x74, 0x61, 0x73,
	0x6b, 0x49, 0x64, 0x73, 0x22, 0xa4, 0x01, 0x0a, 0x12, 0x51, 0x75, 0x65, 0x72, 0x79, 0x54, 0x61,
	0x73, 0x6b, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x05, 0x74,
	0x61, 0x73, 0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x6a, 0x6f, 0x62,
	0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x1a, 0x53, 0x0a, 0x08, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x72,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x2d, 0x0a, 0x12, 0x50,
	0x61, 0x79, 0x4c, 0x6f, 0x61, 0x64, 0x54, 0x61, 0x73, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75,
	0x74, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x22, 0x1f, 0x0a, 0x0b, 0x50, 0x61,
	0x79, 0x4c, 0x6f, 0x61, 0x64, 0x54, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x39, 0x0a, 0x1c, 0x50,
	0x61, 0x79, 0x4c, 0x6f, 0x61, 0x64, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x50, 0x61, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x22, 0x32, 0x0a, 0x13, 0x50, 0x61, 0x79, 0x4c, 0x6f, 0x61,
	0x64, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x1b, 0x0a,
	0x09, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x6e, 0x75, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x4e, 0x75, 0x6d, 0x22, 0x2d, 0x0a, 0x12, 0x51, 0x75,
	0x65, 0x75, 0x69, 0x6e, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x17, 0x0a, 0x07, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x22, 0x48, 0x0a, 0x13, 0x51, 0x75, 0x65,
	0x75, 0x69, 0x6e, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x17, 0x0a, 0x07, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x63,
	0x6f, 0x6e, 0x64, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x73, 0x65, 0x63, 0x6f,
	0x6e, 0x64, 0x73, 0x2a, 0x81, 0x02, 0x0a, 0x07, 0x4a, 0x6f, 0x62, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x0c, 0x0a, 0x08, 0x54, 0x45, 0x53, 0x54, 0x5f, 0x4a, 0x4f, 0x42, 0x10, 0x00, 0x12, 0x11, 0x0a,
	0x0d, 0x54, 0x45, 0x53, 0x54, 0x5f, 0x43, 0x52, 0x4f, 0x4e, 0x5f, 0x4a, 0x4f, 0x42, 0x10, 0x01,
	0x12, 0x1a, 0x0a, 0x16, 0x47, 0x4f, 0x4f, 0x47, 0x4c, 0x45, 0x5f, 0x49, 0x4e, 0x53, 0x54, 0x41,
	0x4e, 0x43, 0x45, 0x5f, 0x41, 0x44, 0x4a, 0x55, 0x53, 0x54, 0x10, 0x02, 0x12, 0x18, 0x0a, 0x14,
	0x53, 0x55, 0x42, 0x53, 0x43, 0x52, 0x49, 0x50, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x52, 0x45, 0x4e,
	0x45, 0x57, 0x41, 0x4c, 0x10, 0x03, 0x12, 0x17, 0x0a, 0x13, 0x54, 0x41, 0x53, 0x4b, 0x5f, 0x54,
	0x49, 0x4d, 0x45, 0x4f, 0x55, 0x54, 0x5f, 0x52, 0x45, 0x46, 0x55, 0x4e, 0x44, 0x10, 0x04, 0x12,
	0x15, 0x0a, 0x11, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x5f, 0x56, 0x49, 0x44, 0x45, 0x4f, 0x5f,
	0x53, 0x4f, 0x52, 0x54, 0x10, 0x05, 0x12, 0x1b, 0x0a, 0x17, 0x43, 0x4c, 0x45, 0x41, 0x4e, 0x5f,
	0x55, 0x53, 0x45, 0x52, 0x5f, 0x55, 0x50, 0x4c, 0x4f, 0x41, 0x44, 0x5f, 0x46, 0x49, 0x4c, 0x45,
	0x53, 0x10, 0x06, 0x12, 0x17, 0x0a, 0x13, 0x43, 0x4c, 0x45, 0x41, 0x4e, 0x5f, 0x4f, 0x4c, 0x44,
	0x5f, 0x4c, 0x4f, 0x47, 0x5f, 0x46, 0x49, 0x4c, 0x45, 0x53, 0x10, 0x07, 0x12, 0x1d, 0x0a, 0x19,
	0x43, 0x4c, 0x45, 0x41, 0x4e, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x45, 0x58, 0x50, 0x49, 0x52,
	0x45, 0x44, 0x5f, 0x50, 0x4f, 0x49, 0x4e, 0x54, 0x53, 0x10, 0x08, 0x12, 0x1a, 0x0a, 0x16, 0x52,
	0x45, 0x53, 0x45, 0x54, 0x5f, 0x44, 0x41, 0x49, 0x4c, 0x59, 0x5f, 0x43, 0x48, 0x41, 0x52, 0x41,
	0x43, 0x54, 0x45, 0x52, 0x53, 0x10, 0x09, 0x2a, 0x40, 0x0a, 0x08, 0x41, 0x64, 0x6d, 0x69, 0x6e,
	0x4a, 0x6f, 0x62, 0x12, 0x19, 0x0a, 0x15, 0x47, 0x45, 0x4e, 0x45, 0x52, 0x41, 0x54, 0x45, 0x5f,
	0x44, 0x41, 0x49, 0x4c, 0x59, 0x5f, 0x52, 0x45, 0x50, 0x4f, 0x52, 0x54, 0x10, 0x00, 0x12, 0x19,
	0x0a, 0x15, 0x47, 0x45, 0x4e, 0x45, 0x52, 0x41, 0x54, 0x45, 0x5f, 0x54, 0x4f, 0x44, 0x41, 0x59,
	0x5f, 0x52, 0x45, 0x50, 0x4f, 0x52, 0x54, 0x10, 0x01, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x79, 0x64, 0x73, 0x73, 0x78, 0x2f, 0x6b, 0x72,
	0x61, 0x74, 0x6f, 0x73, 0x2d, 0x6b, 0x69, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6a, 0x6f, 0x62,
	0x2f, 0x76, 0x31, 0x3b, 0x6a, 0x6f, 0x62, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_api_job_v1_job_proto_rawDescOnce sync.Once
	file_api_job_v1_job_proto_rawDescData = file_api_job_v1_job_proto_rawDesc
)

func file_api_job_v1_job_proto_rawDescGZIP() []byte {
	file_api_job_v1_job_proto_rawDescOnce.Do(func() {
		file_api_job_v1_job_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_job_v1_job_proto_rawDescData)
	})
	return file_api_job_v1_job_proto_rawDescData
}

var file_api_job_v1_job_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_api_job_v1_job_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_api_job_v1_job_proto_goTypes = []any{
	(JobType)(0),                         // 0: job.v1.JobType
	(AdminJob)(0),                        // 1: job.v1.AdminJob
	(*EnqueueRequest)(nil),               // 2: job.v1.EnqueueRequest
	(*EnqueueResponse)(nil),              // 3: job.v1.EnqueueResponse
	(*QueryTasksRequest)(nil),            // 4: job.v1.QueryTasksRequest
	(*QueryTasksResponse)(nil),           // 5: job.v1.QueryTasksResponse
	(*PayLoadTaskTimeout)(nil),           // 6: job.v1.PayLoadTaskTimeout
	(*PayLoadTest)(nil),                  // 7: job.v1.PayLoadTest
	(*PayLoadOrderPaymentCompleted)(nil), // 8: job.v1.PayLoadOrderPaymentCompleted
	(*PayLoadOrderTimeout)(nil),          // 9: job.v1.PayLoadOrderTimeout
	(*QueuingTimeRequest)(nil),           // 10: job.v1.QueuingTimeRequest
	(*QueuingTimeResponse)(nil),          // 11: job.v1.QueuingTimeResponse
	(*QueryTasksResponse_TaskInfo)(nil),  // 12: job.v1.QueryTasksResponse.TaskInfo
	(*timestamppb.Timestamp)(nil),        // 13: google.protobuf.Timestamp
	(*durationpb.Duration)(nil),          // 14: google.protobuf.Duration
}
var file_api_job_v1_job_proto_depIdxs = []int32{
	0,  // 0: job.v1.EnqueueRequest.job_type:type_name -> job.v1.JobType
	13, // 1: job.v1.EnqueueRequest.process_at:type_name -> google.protobuf.Timestamp
	14, // 2: job.v1.EnqueueRequest.process_in:type_name -> google.protobuf.Duration
	13, // 3: job.v1.EnqueueRequest.deadline:type_name -> google.protobuf.Timestamp
	14, // 4: job.v1.EnqueueRequest.retention:type_name -> google.protobuf.Duration
	12, // 5: job.v1.QueryTasksResponse.tasks:type_name -> job.v1.QueryTasksResponse.TaskInfo
	6,  // [6:6] is the sub-list for method output_type
	6,  // [6:6] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_api_job_v1_job_proto_init() }
func file_api_job_v1_job_proto_init() {
	if File_api_job_v1_job_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_job_v1_job_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_job_v1_job_proto_goTypes,
		DependencyIndexes: file_api_job_v1_job_proto_depIdxs,
		EnumInfos:         file_api_job_v1_job_proto_enumTypes,
		MessageInfos:      file_api_job_v1_job_proto_msgTypes,
	}.Build()
	File_api_job_v1_job_proto = out.File
	file_api_job_v1_job_proto_rawDesc = nil
	file_api_job_v1_job_proto_goTypes = nil
	file_api_job_v1_job_proto_depIdxs = nil
}
