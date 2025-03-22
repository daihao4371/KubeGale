package system

import (
	"KubeGale/global"
	"KubeGale/model/common/response"
	"KubeGale/model/system"
	"errors"

	"KubeGale/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BaseApi struct{}

func (b *BaseApi) SignUp(c *gin.Context) {
	var user system.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		response.FailWithMessage("请求参数错误: "+err.Error(), c)
		return
	}

	// 参数基本验证
	if user.Username == "" {
		response.FailWithMessage("用户名不能为空", c)
		return
	}
	if user.Password == "" {
		response.FailWithMessage("密码不能为空", c)
		return
	}

	// 调用服务层方法
	userService := service.ServiceGroupApp.SystemServiceGroup.UserService
	if err := userService.SignUp(&user); err != nil {
		// 处理特定错误
		// 使用自定义错误类型进行比较，而不是直接比较错误码
		var sysErr *global.SysError
		if errors.As(err, &sysErr) {
			switch sysErr.Code {
			case global.ERROR_USER_ALREADY_EXIST:
				response.FailWithMessage("用户名已存在", c)
				return
			case global.ERROR_MOBILE_INVALID:
				response.FailWithMessage("手机号格式不正确", c)
				return
			case global.ERROR_MOBILE_USED:
				response.FailWithMessage("手机号已被使用", c)
				return
			}
		}

		// 记录未预期的错误
		global.KUBEGALE_LOG.Error("用户注册失败", zap.Error(err))
		response.FailWithMessage("注册失败: "+err.Error(), c)
		return
	}

	// 注册成功
	response.OkWithMessage("注册成功", c)
}
