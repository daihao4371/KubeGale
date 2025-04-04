package middleware

import (
	ijwt "KubeGale/utils"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type CasbinMiddleware struct {
	enforcer *casbin.Enforcer
}

func NewCasbinMiddleware(enforcer *casbin.Enforcer) *CasbinMiddleware {
	return &CasbinMiddleware{
		enforcer: enforcer,
	}
}

func (cm *CasbinMiddleware) CheckCasbin() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		// 如果请求的路径是下述路径，则不进行权限验证
		if path == "/api/user/login" ||
			path == "/api/user/logout" ||
			strings.Contains(path, "hello") ||
			path == "/api/user/refresh_token" ||
			path == "/api/user/codes" ||
			path == "/api/user/signup" ||
			path == "/api/not_auth/getTreeNodeBindIps" ||
			path == "/api/monitor/prometheus_configs/prometheus" ||
			path == "/api/monitor/prometheus_configs/prometheus_alert" ||
			path == "/api/monitor/prometheus_configs/prometheus_record" ||
			path == "/api/monitor/prometheus_configs/alertManager" {
			c.Next()
			return
		}

		// 获取用户身份
		userClaims, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1,
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		sub, ok := userClaims.(ijwt.UserClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1,
				"message": "Invalid user claims",
			})
			c.Abort()
			return
		}

		if sub.Uid == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1,
				"message": "Invalid user ID",
			})
			c.Abort()
			return
		}

		//// 打印已有的策略
		//res, err := cm.enforcer.GetPolicy()
		//if err != nil {
		//	log.Println("获取策略失败", err)
		//}
		//
		//log.Println("已有的策略", res)

		// 将用户ID转换为字符串
		userIDStr := strconv.Itoa(sub.Uid)
		act := c.Request.Method

		cm.enforcer.LoadPolicy()

		// 使用 Casbin 检查权限
		ok, err := cm.enforcer.Enforce(userIDStr, path, act)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    1,
				"message": "Error occurred when enforcing policy",
			})
			c.Abort()
			return
		}

		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    1,
				"message": "You don't have permission to access this system",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
