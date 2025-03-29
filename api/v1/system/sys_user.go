package system

import (
	"KubeGale/global"
	"KubeGale/model/common/response"
	"KubeGale/model/system"
	"KubeGale/utils"
	ijwt "KubeGale/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type BaseApi struct{}

// SignUp 用户注册
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
	err = userService.SignUp(&user)
	if err != nil {
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

// Login 用户登录
func (b *BaseApi) Login(c *gin.Context) {
	var user system.User

	if err := c.ShouldBindJSON(&user); err != nil {
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
	ur, err := userService.Login(&user)
	if err != nil {
		// 处理特定错误
		var sysErr *global.SysError
		if errors.As(err, &sysErr) {
			switch sysErr.Code {
			case global.ERROR_USER_NOT_EXIST:
				response.FailWithMessage("用户不存在", c)
				return
			case global.ERROR_PASSWORD_WRONG:
				response.FailWithMessage("密码不正确", c)
				return
			case global.ERROR_USER_DISABLED:
				response.FailWithMessage("用户已被禁用", c)
				return
			}
		}

		// 记录未预期的错误
		global.KUBEGALE_LOG.Error("登录失败", zap.Error(err))
		response.FailWithMessage("登录失败: "+err.Error(), c)
		return
	}

	// 生成JWT令牌
	jwtHandler := utils.NewJWTHandler(global.KUBEGALE_REDIS)
	accessToken, refreshToken, err := jwtHandler.SetLoginToken(c, ur.ID)
	if err != nil {
		global.KUBEGALE_LOG.Error("生成令牌失败", zap.Error(err))
		response.FailWithMessage("登录失败: 生成令牌错误", c)
		return
	}

	// 登录成功，返回用户信息和令牌
	response.OkWithDetailed(gin.H{
		"id":           ur.ID,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"roles":        ur.Roles,
		"desc":         ur.Desc,
		"realName":     ur.RealName,
		"userId":       ur.ID,
		"username":     ur.Username,
	}, "登录成功", c)
}

// Logout 注销登录逻辑
func (b *BaseApi) Logout(c *gin.Context) {
	// 创建JWT处理器
	jwtHandler := utils.NewJWTHandler(global.KUBEGALE_REDIS)

	// 清除token
	if err := jwtHandler.ClearToken(c); err != nil {
		global.KUBEGALE_LOG.Error("注销登录失败", zap.Error(err))
		response.FailWithMessage("注销失败: "+err.Error(), c)
		return
	}

	// 清除cookie中的token
	utils.ClearToken(c)

	// 返回成功响应
	response.OkWithMessage("注销成功", c)
}

// GetProfile 获取用户信息
func (b *BaseApi) GetProfile(c *gin.Context) {
	// 添加调试日志
	global.KUBEGALE_LOG.Info("开始获取用户信息")

	// 从上下文中获取用户信息
	userClaims, exists := c.Get("user")
	if !exists {
		global.KUBEGALE_LOG.Error("上下文中未找到用户信息")

		// 尝试从请求头中直接解析token
		tokenStr := c.GetHeader("Authorization")
		if tokenStr != "" {
			global.KUBEGALE_LOG.Info("尝试从请求头直接解析token", zap.String("token", tokenStr))

			// 提取实际的token字符串（去掉Bearer前缀）
			if strings.HasPrefix(tokenStr, "Bearer ") {
				tokenStr = tokenStr[7:]
			} else if strings.HasPrefix(tokenStr, "bearer ") {
				tokenStr = tokenStr[7:]
			}

			if tokenStr != "" {
				// 解析token
				var uc ijwt.UserClaims
				key := []byte(viper.GetString("jwt.key1"))
				token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
					return key, nil
				})

				if err == nil && token != nil && token.Valid {
					// 检查会话是否有效
					jwtHandler := utils.NewJWTHandler(global.KUBEGALE_REDIS)
					err = jwtHandler.CheckSession(c, uc.Ssid)
					if err == nil {
						// 设置用户信息到上下文
						c.Set("user", uc)
						global.KUBEGALE_LOG.Info("从请求头解析token成功，已设置用户信息", zap.Int("uid", uc.Uid))
						userClaims = uc
						exists = true
					} else {
						global.KUBEGALE_LOG.Error("会话检查失败", zap.Error(err))
					}
				} else {
					global.KUBEGALE_LOG.Error("token解析失败", zap.Error(err))
				}
			}
		}

		// 如果仍然无法获取用户信息，返回错误
		if !exists {
			response.FailWithMessage("用户未登录", c)
			return
		}
	}

	global.KUBEGALE_LOG.Info("获取到用户上下文", zap.Any("userClaims", userClaims))

	// 尝试多种类型断言方式
	var uid int
	var found bool

	// 尝试直接断言为UserClaims
	if claims, ok := userClaims.(ijwt.UserClaims); ok {
		uid = claims.Uid
		found = true
		global.KUBEGALE_LOG.Info("从UserClaims中获取到用户ID", zap.Int("uid", uid))
	}

	// 尝试断言为指针类型
	if !found {
		if claims, ok := userClaims.(*ijwt.UserClaims); ok {
			uid = claims.Uid
			found = true
			global.KUBEGALE_LOG.Info("从*UserClaims中获取到用户ID", zap.Int("uid", uid))
		}
	}

	// 尝试断言为map
	if !found {
		if claimsMap, mapOk := userClaims.(map[string]interface{}); mapOk {
			// 尝试从map中获取uid
			if uidVal, exists := claimsMap["Uid"]; exists {
				// 尝试不同类型的转换
				switch v := uidVal.(type) {
				case int:
					uid = v
					found = true
				case float64:
					uid = int(v)
					found = true
				case string:
					if id, err := strconv.Atoi(v); err == nil {
						uid = id
						found = true
					}
				}
				if found {
					global.KUBEGALE_LOG.Info("从map中获取到用户ID", zap.Int("uid", uid))
				}
			}
		}
	}

	// 如果仍未找到，记录详细信息并返回错误
	if !found {
		global.KUBEGALE_LOG.Error("用户信息类型错误",
			zap.String("actualType", fmt.Sprintf("%T", userClaims)),
			zap.Any("value", userClaims))
		response.FailWithMessage("用户信息类型错误", c)
		return
	}

	// 调用服务层获取用户详细信息
	global.KUBEGALE_LOG.Info("开始调用服务层获取用户信息", zap.Int("uid", uid))
	user, err := userService.GetProfile(uid)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取用户信息失败", zap.Error(err))
		response.FailWithMessage("获取用户信息失败: "+err.Error(), c)
		return
	}

	// 返回用户信息
	response.OkWithDetailed(gin.H{
		"id":           user.ID,
		"roles":        user.Roles,
		"realName":     user.RealName,
		"userId":       user.ID,
		"username":     user.Username,
		"desc":         user.Desc,
		"homePath":     user.HomePath,
		"mobile":       user.Mobile,
		"feiShuUserId": user.FeiShuUserId,
		"menus":        user.Menus,
		"apis":         user.Apis,
		"Roles":        user.Roles,
	}, "获取用户信息成功", c)
}

// RefreshToken 刷新token
func (b *BaseApi) RefreshToken(c *gin.Context) {
	var req system.TokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数错误: "+err.Error(), c)
		return
	}

	if req.RefreshToken == "" {
		response.FailWithMessage("刷新令牌不能为空", c)
		return
	}

	// 获取密钥
	key := viper.GetString("jwt.key2")

	// 解析 token 并获取刷新 claims
	rc := ijwt.RefreshClaims{}
	token, err := jwt.ParseWithClaims(req.RefreshToken, &rc, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		global.KUBEGALE_LOG.Error("token刷新失败", zap.Error(err))
		response.FailWithMessage("token解析失败", c)
		return
	}

	// 检查 token 是否有效
	if token == nil || !token.Valid {
		global.KUBEGALE_LOG.Warn("无效的token")
		response.FailWithMessage("token无效", c)
		return
	}

	// 创建JWT处理器
	jwtHandler := utils.NewJWTHandler(global.KUBEGALE_REDIS)

	// 检查会话状态是否异常
	if err = jwtHandler.CheckSession(c, rc.Ssid); err != nil {
		global.KUBEGALE_LOG.Error("会话检查失败", zap.Error(err))
		response.FailWithMessage("会话检查失败: "+err.Error(), c)
		return
	}

	// 生成新的访问令牌和刷新令牌
	accessToken, refreshToken, err := jwtHandler.SetLoginToken(c, rc.Uid)
	if err != nil {
		global.KUBEGALE_LOG.Error("生成新令牌失败", zap.Error(err))
		response.FailWithMessage("生成新token失败: "+err.Error(), c)
		return
	}

	// 获取用户信息
	user, err := userService.GetProfile(rc.Uid)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取用户信息失败", zap.Error(err))
		response.FailWithMessage("获取用户信息失败: "+err.Error(), c)
		return
	}

	// 返回新的令牌和用户信息
	response.OkWithDetailed(gin.H{
		"id":           user.ID,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"roles":        user.Roles,
		"desc":         user.Desc,
		"realName":     user.RealName,
		"userId":       user.ID,
		"username":     user.Username,
	}, "刷新令牌成功", c)
}

// GetPermCode 获取用户权限码
func (b *BaseApi) GetPermCode(c *gin.Context) {
	// 从上下文中获取用户信息
	userClaims, exists := c.Get("user")
	if !exists {
		global.KUBEGALE_LOG.Error("上下文中未找到用户信息")
		response.FailWithMessage("用户未登录", c)
		return
	}

	global.KUBEGALE_LOG.Info("获取到用户上下文", zap.Any("userClaims", userClaims))

	// 尝试多种类型断言方式
	var uid int
	var found bool

	// 尝试直接断言为UserClaims
	if claims, ok := userClaims.(ijwt.UserClaims); ok {
		uid = claims.Uid
		found = true
		global.KUBEGALE_LOG.Info("从UserClaims中获取到用户ID", zap.Int("uid", uid))
	}

	// 尝试断言为指针类型
	if !found {
		if claims, ok := userClaims.(*ijwt.UserClaims); ok {
			uid = claims.Uid
			found = true
			global.KUBEGALE_LOG.Info("从*UserClaims中获取到用户ID", zap.Int("uid", uid))
		}
	}

	// 尝试断言为map
	if !found {
		if claimsMap, mapOk := userClaims.(map[string]interface{}); mapOk {
			// 尝试从map中获取uid
			if uidVal, exists := claimsMap["Uid"]; exists {
				// 尝试不同类型的转换
				switch v := uidVal.(type) {
				case int:
					uid = v
					found = true
				case float64:
					uid = int(v)
					found = true
				case string:
					if id, err := strconv.Atoi(v); err == nil {
						uid = id
						found = true
					}
				}
				if found {
					global.KUBEGALE_LOG.Info("从map中获取到用户ID", zap.Int("uid", uid))
				}
			}
		}
	}

	// 如果仍未找到，记录详细信息并返回错误
	if !found {
		global.KUBEGALE_LOG.Error("用户信息类型错误",
			zap.String("actualType", fmt.Sprintf("%T", userClaims)),
			zap.Any("value", userClaims))
		response.FailWithMessage("用户信息类型错误", c)
		return
	}

	// 调用服务层获取用户权限码
	codes, err := userService.GetPermCode(uid)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取用户权限码失败", zap.Error(err))
		response.FailWithMessage("获取用户权限码失败: "+err.Error(), c)
		return
	}

	// 返回权限码
	response.OkWithData(codes, c)
}

// GetUserList 获取用户列表
func (b *BaseApi) GetUserList(c *gin.Context) {
	// 获取分页参数
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	keyword := c.DefaultQuery("keyword", "")

	// 转换为整数
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		response.FailWithMessage("页码参数错误", c)
		return
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil {
		response.FailWithMessage("每页数量参数错误", c)
		return
	}

	// 调用服务层获取用户列表
	users, total, err := userService.GetUserList(pageNum, pageSizeNum, keyword)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取用户列表失败", zap.Error(err))
		response.FailWithMessage("获取用户列表失败: "+err.Error(), c)
		return
	}

	// 返回用户列表和总数
	response.OkWithDetailed(gin.H{
		"list":     users,
		"total":    total,
		"page":     pageNum,
		"pageSize": pageSizeNum,
	}, "获取用户列表成功", c)
}

// ChangePassword 修改用户密码
func (b *BaseApi) ChangePassword(c *gin.Context) {
	var req system.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数错误: "+err.Error(), c)
		return
	}

	// 验证新密码和确认密码是否一致
	if req.NewPassword != req.ConfirmPassword {
		response.FailWithMessage("新密码和确认密码不一致", c)
		return
	}

	// 从上下文中获取用户信息
	userClaims, exists := c.Get("user")
	if !exists {
		response.FailWithMessage("用户未登录", c)
		return
	}

	// 类型断言获取用户ID
	claims, ok := userClaims.(utils.UserClaims)
	if !ok {
		global.KUBEGALE_LOG.Error("用户信息类型错误")
		response.FailWithMessage("用户信息类型错误", c)
		return
	}

	// 调用服务层修改密码
	if err := userService.ChangePassword(claims.Uid, req.Password, req.NewPassword); err != nil {
		// 处理特定错误
		var sysErr *global.SysError
		if errors.As(err, &sysErr) {
			switch sysErr.Code {
			case global.ERROR_OLD_PASSWORD_WRONG:
				response.FailWithMessage("原密码不正确", c)
				return
			case global.ERROR_USER_NOT_EXIST:
				response.FailWithMessage("用户不存在", c)
				return
			case global.ERROR_SAME_PASSWORD:
				response.FailWithMessage("新密码不能与旧密码相同", c)
				return
			}
		}

		// 记录未预期的错误
		global.KUBEGALE_LOG.Error("修改密码失败", zap.Error(err))
		response.FailWithMessage("修改密码失败: "+err.Error(), c)
		return
	}

	// 返回成功响应
	response.OkWithMessage("密码修改成功", c)
}

// WriteOff 注销用户
func (b *BaseApi) WriteOff(c *gin.Context) {
	var req system.WriteOffRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数错误: "+err.Error(), c)
		return
	}

	// 验证用户名不能为空
	if req.Username == "" {
		response.FailWithMessage("用户名不能为空", c)
		return
	}

	// 验证密码不能为空
	if req.Password == "" {
		response.FailWithMessage("密码不能为空", c)
		return
	}

	// 调用服务层注销用户
	if err := userService.WriteOff(req.Username, req.Password); err != nil {
		// 处理特定错误
		var sysErr *global.SysError
		if errors.As(err, &sysErr) {
			switch sysErr.Code {
			case global.ERROR_USER_NOT_EXIST:
				response.FailWithMessage("用户不存在", c)
				return
			case global.ERROR_PASSWORD_WRONG:
				response.FailWithMessage("密码不正确", c)
				return
			}
		}

		// 记录未预期的错误
		global.KUBEGALE_LOG.Error("注销账号失败", zap.Error(err))
		response.FailWithMessage("注销账号失败: "+err.Error(), c)
		return
	}

	// 返回成功响应
	response.OkWithMessage("账号注销成功", c)
}

// UpdateProfile 修改用户信息
func (b *BaseApi) UpdateProfile(c *gin.Context) {
	var req system.UpdateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数错误: "+err.Error(), c)
		return
	}

	// 验证用户ID
	if req.UserId <= 0 {
		response.FailWithMessage("用户ID不能为空", c)
		return
	}

	// 构建用户信息
	user := &system.User{
		RealName:     req.RealName,
		Desc:         req.Desc,
		Mobile:       req.Mobile,
		FeiShuUserId: req.FeiShuUserId,
		AccountType:  req.AccountType,
		HomePath:     req.HomePath,
		Enable:       req.Enable,
	}

	// 调用服务层更新用户信息
	if err := userService.UpdateProfile(req.UserId, user); err != nil {
		// 处理特定错误
		var sysErr *global.SysError
		if errors.As(err, &sysErr) {
			switch sysErr.Code {
			case global.ERROR_USER_NOT_EXIST:
				response.FailWithMessage("用户不存在", c)
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
		global.KUBEGALE_LOG.Error("更新用户信息失败", zap.Error(err))
		response.FailWithMessage("更新用户信息失败: "+err.Error(), c)
		return
	}

	// 返回成功响应
	response.OkWithMessage("用户信息更新成功", c)
}

// DeleteUser 删除用户
func (b *BaseApi) DeleteUser(c *gin.Context) {
	// 获取路径参数中的用户ID
	id := c.Param("id")

	// 转换为整数
	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.FailWithMessage("用户ID参数错误", c)
		return
	}

	// 验证用户ID
	if idInt <= 0 {
		response.FailWithMessage("无效的用户ID", c)
		return
	}

	// 调用服务层删除用户
	if err := userService.DeleteUser(idInt); err != nil {
		// 处理特定错误
		var sysErr *global.SysError
		if errors.As(err, &sysErr) {
			switch sysErr.Code {
			case global.ERROR_USER_NOT_EXIST:
				response.FailWithMessage("用户不存在", c)
				return
			case global.ERROR_USER_ID_INVALID:
				response.FailWithMessage("用户ID无效", c)
				return
			}
		}

		// 记录未预期的错误
		global.KUBEGALE_LOG.Error("删除用户失败", zap.Error(err))
		response.FailWithMessage("删除用户失败: "+err.Error(), c)
		return
	}

	// 返回成功响应
	response.OkWithMessage("用户删除成功", c)
}
