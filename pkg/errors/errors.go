package errors

import (
	serrors "errors"
	"fmt"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/pkg/errors"
)

var (
	// ErrNotFound 资源不存在
	ErrNotFound = New("resource not found")
	// ErrAlreadyExists 资源已存在
	ErrAlreadyExists = New("resource already exists")
	// Unauthorized 未授权
	Unauthorized = kerrors.New(401, "unauthorized", "unauthorized")
	// ErrForbidden 禁止访问
	ErrForbidden = New("forbidden")
	// ErrBadRequest 请求错误
	ErrBadRequest = New("bad request")
	// ErrInternal 内部错误
	ErrInternal = New("internal error")
	// ErrTimeout 超时
	ErrTimeout = New("timeout")
	// ErrUnavailable 资源不可用
	ErrUnavailable = New("resource unavailable")
	// ErrLimited 资源已达上限
	ErrLimited = New("resource limited")
	// ErrConflict 冲突
	ErrConflict = New("conflict")
	// ErrCancelled 已取消
	ErrCancelled = New("cancelled")

	// 积分不足
	ErrInsufficientCredits = kerrors.New(410, "insufficient credits", "insufficient credits")
	// 免费次数不足
	ErrInsufficientFreeCount = kerrors.New(411, "insufficient free count", "insufficient free count")

	ErrVideoLengthExceed15s = kerrors.New(412, "video length exceeds", "The video is longer than 15 seconds, please subscribe to take action")
	// 视频长度超过10分钟
	ErrVideoLengthExceed10m = kerrors.New(-1, "video length exceeds", "Video longer than 10 minutes")
	// 文件大小超过限制
	ErrFileSizeExceed = kerrors.New(-1, "file size exceeds", "Video size exceeds 100Mb")
	// 上传失败请重新上传
	ErrUploadFailed = kerrors.New(-1, "upload failed", "Upload failed, please upload again")

	// 创建订阅失败
	ErrCreateSubscription = kerrors.New(-1, "create subscription failed", "Error message.Please re-enter.")
	// 验证码已发送
	ErrVerificationCodeSent = kerrors.New(413, "verification code sent", "Verification code has been sent")

	New       = errors.New
	Join      = serrors.Join
	Unwrap    = errors.Unwrap
	Is        = errors.Is
	As        = errors.As
	Wrap      = errors.Wrap
	Errorf    = errors.Errorf
	Wrapf     = errors.Wrapf
	WithStack = errors.WithStack
)

// UserError 表示用户级别的错误
type UserError struct {
	Ke *kerrors.Error
}

func (e *UserError) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s ", e.Ke.Code, e.Ke.Reason, e.Ke.Message)
}

// NewUserError 创建用户级别的错误
func NewUserError(message string) *UserError {
	return &UserError{
		Ke: kerrors.New(-1, message, message),
	}
}

// WrapAsUserError 将错误包装为用户级别的错误
func WrapAsUserError(err error) *UserError {
	return &UserError{Ke: kerrors.FromError(err)}
}
