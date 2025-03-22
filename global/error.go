package global

import "fmt"

// SysError 系统错误类型
type SysError struct {
	Code int    // 错误码
	Msg  string // 错误信息
}

// Error 实现 error 接口
func (e *SysError) Error() string {
	return fmt.Sprintf("错误码: %d, 错误信息: %s", e.Code, e.Msg)
}

// CreateSysError 创建一个新的系统错误
func CreateSysError(code int) *SysError {
	return &SysError{
		Code: code,
		Msg:  ErrorMsg[code],
	}
}

// CreateSysErrorWithMsg 创建一个带自定义消息的系统错误
func CreateSysErrorWithMsg(code int, msg string) *SysError {
	return &SysError{
		Code: code,
		Msg:  msg,
	}
}