package system

import (
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (s *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	// 创建用户路由组
	userGroup := Router.Group("user") // 注意这里不要加前导斜杠

	// 需要记录操作日志的路由
	userRouterWithRecord := userGroup.Use(middleware.OperationRecord())

	// 不需要记录操作日志的路由
	userRouterWithoutRecord := userGroup

	// 获取API控制器

	// 不需要记录操作日志的接口
	{
		userRouterWithoutRecord.POST("/signup", baseApi.SignUp)              // 注册
		userRouterWithoutRecord.POST("/login", baseApi.Login)                // 登录
		userRouterWithoutRecord.POST("/refresh_token", baseApi.RefreshToken) // 刷新token
		userRouterWithoutRecord.POST("/logout", baseApi.Logout)              // 退出登录
		userRouterWithoutRecord.GET("/profile", baseApi.GetProfile)          // 用户信息
		userRouterWithoutRecord.GET("/codes", baseApi.GetPermCode)           // 前端所需状态码
	}

	// 需要记录操作日志的接口
	{
		userRouterWithRecord.GET("/list", baseApi.GetUserList)                // 用户列表
		userRouterWithRecord.POST("/change_password", baseApi.ChangePassword) // 修改密码
		userRouterWithRecord.POST("/write_off", baseApi.WriteOff)             // 注销账号
		userRouterWithRecord.POST("/profile/update", baseApi.UpdateProfile)   // 更新用户信息
		userRouterWithRecord.DELETE("/:id", baseApi.DeleteUser)               // 删除用户
	}
}
