package global

// 系统错误码常量
const (
	// 通用错误
	SUCCESS               = 0
	ERROR                 = 1
	ERROR_PARAM_INVALID   = 1001
	ERROR_INTERNAL_SERVER = 1002
	ERROR_UNAUTHORIZED    = 1003
	ERROR_FORBIDDEN       = 1004
	ERROR_NOT_FOUND       = 1005

	// 用户相关错误
	ERROR_USER_NOT_EXIST     = 2001
	ERROR_USER_ALREADY_EXIST = 2002
	ERROR_PASSWORD_WRONG     = 2003
	ERROR_USER_DISABLED      = 2004
	ERROR_MOBILE_INVALID     = 2005
	ERROR_MOBILE_USED        = 2006
	ERROR_OLD_PASSWORD_WRONG = 2007
	ERROR_SAME_PASSWORD      = 2008
	ERROR_USER_ID_INVALID    = 2009
)

// 系统错误信息常量
var ErrorMsg = map[int]string{
	SUCCESS:               "操作成功",
	ERROR:                 "操作失败",
	ERROR_PARAM_INVALID:   "参数无效",
	ERROR_INTERNAL_SERVER: "服务器内部错误",
	ERROR_UNAUTHORIZED:    "未授权",
	ERROR_FORBIDDEN:       "禁止访问",
	ERROR_NOT_FOUND:       "资源不存在",

	ERROR_USER_NOT_EXIST:     "用户不存在",
	ERROR_USER_ALREADY_EXIST: "用户名已存在",
	ERROR_PASSWORD_WRONG:     "密码不正确",
	ERROR_USER_DISABLED:      "用户已被禁用",
	ERROR_MOBILE_INVALID:     "手机号格式不正确",
	ERROR_MOBILE_USED:        "手机号已被使用",
	ERROR_OLD_PASSWORD_WRONG: "旧密码不正确",
	ERROR_SAME_PASSWORD:      "新密码不能与旧密码相同",
	ERROR_USER_ID_INVALID:    "无效的用户ID",
}

// 自定义错误类型
type CustomError struct {
	Code int
	Msg  string
}

func (e *CustomError) Error() string {
	return e.Msg
}

// 创建自定义错误
func NewError(code int) *CustomError {
	return &CustomError{
		Code: code,
		Msg:  ErrorMsg[code],
	}
}

// 创建自定义错误并指定消息
func NewErrorWithMsg(code int, msg string) *CustomError {
	return &CustomError{
		Code: code,
		Msg:  msg,
	}
}