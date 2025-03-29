package middleware

import (
	"KubeGale/global"
	ijwt "KubeGale/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
)

type JWTMiddleware struct {
	ijwt.Handler
}

func NewJWTMiddleware(hdl ijwt.Handler) *JWTMiddleware {
	return &JWTMiddleware{
		Handler: hdl,
	}
}

// CheckLogin 校验JWT
func (m *JWTMiddleware) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		// 添加调试日志
		global.KUBEGALE_LOG.Info("请求路径", zap.String("path", path))
		
		// 记录所有请求头，帮助调试
		headers := make(map[string]string)
		for k, v := range ctx.Request.Header {
			if len(v) > 0 {
				headers[k] = v[0]
			}
		}
		global.KUBEGALE_LOG.Info("请求头", zap.Any("headers", headers))

		// 如果请求的路径是下述路径，则不进行token验证
		if path == "/api/user/login" ||
			path == "/api/user/logout" ||
			strings.Contains(path, "hello") ||
			path == "/api/user/refresh_token" ||
			path == "/api/user/signup" ||
			path == "/api/not_auth/getTreeNodeBindIps" ||
			path == "/api/monitor/prometheus_configs/prometheus" ||
			path == "/api/monitor/prometheus_configs/prometheus_alert" ||
			path == "/api/monitor/prometheus_configs/prometheus_record" ||
			path == "/api/monitor/prometheus_configs/alertManager" ||
			path == "/health" {
			ctx.Next()
			return
		}

		var uc ijwt.UserClaims
		var tokenStr string

		// 如果是/api/tree/ecs/console开头的路径，从查询参数获取token
		if strings.HasPrefix(path, "/api/tree/ecs/console") {
			tokenStr = ctx.Query("token")
			if tokenStr != "" {
				tokenStr = "Bearer " + tokenStr
			}
		} else {
			// 从请求中提取token
			tokenStr = m.ExtractToken(ctx)
		}

		global.KUBEGALE_LOG.Info("提取的token", zap.String("token", tokenStr))

		// 提取实际的token字符串（去掉Bearer前缀）
		if strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = tokenStr[7:]
		} else if strings.HasPrefix(tokenStr, "bearer ") {
			tokenStr = tokenStr[7:]
		}
		
		// 如果token为空，直接返回未授权
		if tokenStr == "" {
			global.KUBEGALE_LOG.Error("token为空")
			ctx.AbortWithStatusJSON(200, gin.H{"code": 7, "msg": "用户未登录"})
			return
		}

		// 解析token
		key := []byte(viper.GetString("jwt.key1"))
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})

		if err != nil {
			// token 错误
			global.KUBEGALE_LOG.Error("token解析错误", zap.Error(err))
			ctx.AbortWithStatusJSON(200, gin.H{"code": 7, "msg": "用户未登录"})
			return
		}

		if token == nil || !token.Valid {
			// token 非法或过期
			global.KUBEGALE_LOG.Error("token非法或过期")
			ctx.AbortWithStatusJSON(200, gin.H{"code": 7, "msg": "用户未登录"})
			return
		}

		// 检查会话是否有效
		err = m.CheckSession(ctx, uc.Ssid)
		if err != nil {
			global.KUBEGALE_LOG.Error("会话检查失败", zap.Error(err))
			ctx.AbortWithStatusJSON(200, gin.H{"code": 7, "msg": "用户未登录"})
			return
		}

		// 设置用户信息到上下文，确保类型正确
		global.KUBEGALE_LOG.Info("设置用户信息到上下文", 
			zap.Int("uid", uc.Uid), 
			zap.String("ssid", uc.Ssid),
			zap.String("userAgent", uc.UserAgent))
		
		// 直接设置为结构体，而不是指针
		ctx.Set("user", uc)
		
		// 添加额外的调试信息，确认上下文中是否有用户信息
		testUser, exists := ctx.Get("user")
		if exists {
			global.KUBEGALE_LOG.Info("验证上下文中的用户信息", zap.Any("user", testUser))
		} else {
			global.KUBEGALE_LOG.Error("设置后仍无法获取用户信息")
		}
		
		ctx.Next()
	}
}
