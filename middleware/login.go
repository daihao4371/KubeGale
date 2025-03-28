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

		// 如果请求的路径是下述路径，则不进行token验证
		if path == "/api/user/login" ||
			//path == "/api/user/signup" ||   // 不允许用户自己注册账号
			path == "/api/user/logout" ||
			strings.Contains(path, "hello") ||
			path == "/api/user/refresh_token" ||
			path == "/api/user/signup" ||
			path == "/api/not_auth/getTreeNodeBindIps" ||
			path == "/api/monitor/prometheus_configs/prometheus" ||
			path == "/api/monitor/prometheus_configs/prometheus_alert" ||
			path == "/api/monitor/prometheus_configs/prometheus_record" ||
			path == "/api/monitor/prometheus_configs/alertManager" {
			return
		}

		var uc ijwt.UserClaims
		var tokenStr string

		// 如果是/api/tree/ecs/console开头的路径，从查询参数获取token
		if strings.HasPrefix(path, "/api/tree/ecs/console") {
			tokenStr = ctx.Query("token")
		} else {
			// 从请求中提取token
			tokenStr = m.ExtractToken(ctx)
			global.KUBEGALE_LOG.Info("提取的token", zap.String("token", tokenStr))
		}

		// 在CheckLogin函数中
		key := []byte(viper.GetString("jwt.key1"))
		// 解析token
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})

		if err != nil {
			// token 错误
			global.KUBEGALE_LOG.Error("token解析错误", zap.Error(err))
			ctx.AbortWithStatus(401)
			return
		}

		if token == nil || !token.Valid {
			// token 非法或过期
			ctx.AbortWithStatus(401)
			return
		}

		// 检查是否携带ua头
		if uc.UserAgent == "" {
			ctx.AbortWithStatus(401)
			return
		}

		// 检查会话是否有效
		err = m.CheckSession(ctx, uc.Ssid)

		if err != nil {
			ctx.AbortWithStatus(401)
			return
		}

		ctx.Set("user", uc)
	}
}
