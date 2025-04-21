package middleware

import (
	"KubeGale/global"
	"KubeGale/utils"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"

	"go.uber.org/zap"
	
	"KubeGale/model/common/response"
	"KubeGale/service"

	"github.com/gin-gonic/gin"
)

var jwtService = service.ServiceGroupApp.SystemServiceGroup.JwtService

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := utils.GetToken(c)
		if token == "" {
			global.KUBEGALE_LOG.Error("JWT认证失败: token为空")
			response.NoAuth("未登录或非法访问", c)
			c.Abort()
			return
		}
		if jwtService.IsBlacklist(token) {
			global.KUBEGALE_LOG.Error("JWT认证失败: token在黑名单中")
			response.NoAuth("您的帐户异地登陆或令牌失效", c)
			utils.ClearToken(c)
			c.Abort()
			return
		}
		j := utils.NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			global.KUBEGALE_LOG.Error("JWT认证失败: 解析token失败", zap.Error(err))
			if errors.Is(err, utils.TokenExpired) {
				response.NoAuth("授权已过期", c)
				utils.ClearToken(c)
				c.Abort()
				return
			}
			response.NoAuth(err.Error(), c)
			utils.ClearToken(c)
			c.Abort()
			return
		}

		// 记录当前用户信息，便于调试
		global.KUBEGALE_LOG.Info("JWT认证成功", 
			zap.String("username", claims.Username), 
			zap.Uint("authorityId", claims.AuthorityId))
		
		// 已登录用户被管理员禁用 需要使该用户的jwt失效 此处比较消耗性能 如果需要 请自行打开
		// 用户被删除的逻辑 需要优化 此处比较消耗性能 如果需要 请自行打开

		//if user, err := userService.FindUserByUuid(claims.UUID.String()); err != nil || user.Enable == 2 {
		//	_ = jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: token})
		//	response.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
		//	c.Abort()
		//}
		c.Set("claims", claims)
		if claims.ExpiresAt.Unix()-time.Now().Unix() < claims.BufferTime {
			dr, _ := utils.ParseDuration(global.KUBEGALE_CONFIG.JWT.ExpiresTime)
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(dr))
			newToken, _ := j.CreateTokenByOldToken(token, *claims)
			newClaims, _ := j.ParseToken(newToken)
			c.Header("new-token", newToken)
			c.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt.Unix(), 10))
			utils.SetToken(c, newToken, int(dr.Seconds()))
			if global.KUBEGALE_CONFIG.System.UseMultipoint {
				// 记录新的活跃jwt
				_ = jwtService.SetRedisJWT(newToken, newClaims.Username)
			}
		}
		c.Next()

		if newToken, exists := c.Get("new-token"); exists {
			c.Header("new-token", newToken.(string))
		}
		if newExpiresAt, exists := c.Get("new-expires-at"); exists {
			c.Header("new-expires-at", newExpiresAt.(string))
		}
	}
}
