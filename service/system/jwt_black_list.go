package system

import (
	"KubeGale/global"
	"KubeGale/model/system"
	"KubeGale/utils"
	"context"
	"go.uber.org/zap"
)

type JwtService struct{}

var JwtServiceApp = new(JwtService)

// 拉黑jwt
func (jwtService *JwtService) JsonInBlacklist(jwtList system.JwtBlacklist) (err error) {
	err = global.KUBEGALE_DB.Create(&jwtList).Error
	if err != nil {
		return
	}
	global.BlackCache.SetDefault(jwtList.Jwt, struct{}{})
	return
}

// 判断JWT是否在黑名单内部
func (jwtService *JwtService) IsBlacklist(jwt string) bool {
	_, ok := global.BlackCache.Get(jwt)
	return ok
	// err := global.DYCLOUD_DB.Where("jwt = ?", jwt).First(&system.JwtBlacklist{}).Error
	// isNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	// return !isNotFound
}

// 从redis取jwt
func (jwtService *JwtService) GetRedisJWT(userName string) (redisJWT string, err error) {
	redisJWT, err = global.KUBEGALE_REDIS.Get(context.Background(), userName).Result()
	return redisJWT, err
}

// jwt存入redis并设置过期时间
func (jwtService *JwtService) SetRedisJWT(jwt string, userName string) (err error) {
	// 此处过期时间等于jwt过期时间
	dr, err := utils.ParseDuration(global.KUBEGALE_CONFIG.JWT.ExpiresTime)
	if err != nil {
		return err
	}
	timer := dr
	err = global.KUBEGALE_REDIS.Set(context.Background(), userName, jwt, timer).Err()
	return err
}

func LoadAll() {
	var data []string
	err := global.KUBEGALE_DB.Model(&system.JwtBlacklist{}).Select("jwt").Find(&data).Error
	if err != nil {
		global.KUBEGALE_LOG.Error("加载数据库jwt黑名单失败!", zap.Error(err))
		return
	}
	for i := 0; i < len(data); i++ {
		global.BlackCache.SetDefault(data[i], struct{}{})
	} // jwt黑名单 加入 BlackCache 中
}
