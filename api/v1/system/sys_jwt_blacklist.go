package system

import (
	"KubeGale/global"
	"KubeGale/model/common/response"
	"KubeGale/model/system"
	"KubeGale/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type JwtApi struct{}

// JsonInBlacklist  jwt加入黑名单
func (j *JwtApi) JsonInBlacklist(c *gin.Context) {
	token := utils.GetToken(c)
	jwt := system.JwtBlacklist{Jwt: token}
	err := jwtService.JsonInBlacklist(jwt)
	if err != nil {
		global.KUBEGALE_LOG.Error("jwt作废失败!", zap.Error(err))
		response.FailWithMessage("jwt作废失败", c)
		return
	}
	utils.ClearToken(c)
	response.OkWithMessage("jwt作废成功", c)
}
