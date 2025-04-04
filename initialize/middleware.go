package initialize

import (
	"KubeGale/middleware"
	ijwt "KubeGale/utils"
	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
	"time"
)

func InitMiddlewares(ih ijwt.Handler, l *zap.Logger, enforcer *casbin.Enforcer) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
			AllowCredentials: true,
			AllowHeaders:     []string{"Content-Type", "Authorization", "X-Refresh-Token"},
			ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token"},
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					return true
				}
				return strings.Contains(origin, "")
			},
			MaxAge: 12 * time.Hour,
		}),
		middleware.NewJWTMiddleware(ih).CheckLogin(),
		middleware.NewCasbinMiddleware(enforcer).CheckCasbin(),
		middleware.NewLogMiddleware(l).Log(),
	}
}
