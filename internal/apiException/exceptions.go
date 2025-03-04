package apiException

import (
	"net/http"

	"bookrecycle-server/pkg/log"
)

// Error 自定义错误类型
type Error struct {
	Code  int
	Msg   string
	Level log.Level
}

// 自定义错误
var (
	ServerError             = NewError(200500, log.LevelError, "系统异常，请稍后重试")
	ParamsError             = NewError(200501, log.LevelInfo, "参数错误")
	WrongPasswordOrUsername = NewError(200502, log.LevelInfo, "用户名或密码错误")
	UserAlreadyExist        = NewError(200503, log.LevelInfo, "用户已存在")
	InvalidUsername         = NewError(200504, log.LevelInfo, "用户名不合法")
	InvalidPassword         = NewError(200505, log.LevelInfo, "密码不合法")
	NoAccessPermission      = NewError(200506, log.LevelInfo, "无访问权限")
	WebSocketError          = NewError(200507, log.LevelWarn, "WebSocket连接错误")
	FileSizeExceedError     = NewError(200508, log.LevelInfo, "文件大小超限")
	UploadFileError         = NewError(200509, log.LevelInfo, "上传文件失败")
	FileNotImageError       = NewError(200510, log.LevelInfo, "上传的文件不是图片")
	ResourceNotFound        = NewError(200511, log.LevelInfo, "资源不存在")
	UserNotActive           = NewError(200512, log.LevelInfo, "用户未激活")
	BalanceNotEnough        = NewError(200513, log.LevelInfo, "余额不足")
	NotFound                = NewError(200404, log.LevelWarn, http.StatusText(http.StatusNotFound))
)

// Error 实现 error 接口，返回错误的消息内容
func (e *Error) Error() string {
	return e.Msg
}

// NewError 创建并返回一个新的自定义错误实例
func NewError(code int, level log.Level, msg string) *Error {
	return &Error{
		Code:  code,
		Msg:   msg,
		Level: level,
	}
}
