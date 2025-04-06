package system

import (
	"KubeGale/global"
	"KubeGale/model/common/response"
	systemRes "KubeGale/model/system/response"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"time"
)

// 当开启多服务器部署时，替换下面的配置，使用redis共享存储验证码
//var store = captcha.NewDefaultRedisStore()

var store = base64Captcha.DefaultMemStore

// 生成验证码
func (b *BaseApi) Captcha(c *gin.Context) {
	// 判断验证码是否开启
	openCaptcha := global.KUBEGALE_CONFIG.Captcha.OpenCaptcha               // 是否开启防爆次数
	openCaptchaTimeOut := global.KUBEGALE_CONFIG.Captcha.OpenCaptchaTimeOut // 缓存超时时间
	key := c.ClientIP()
	v, ok := global.BlackCache.Get(key)
	if !ok {
		global.BlackCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeOut))
	}

	var oc bool
	if openCaptcha == 0 || openCaptcha < interfaceToInt(v) {
		oc = true
	}
	// 字符,公式,验证码配置
	// 生成默认数字的driver
	driver := base64Captcha.NewDriverDigit(global.KUBEGALE_CONFIG.Captcha.ImgHeight, global.KUBEGALE_CONFIG.Captcha.ImgWidth, global.KUBEGALE_CONFIG.Captcha.KeyLong, 0.7, 80)
	// cp := base64Captcha.NewCaptcha(driver, store.UseWithCtx(c))   // v8下使用redis
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		global.KUBEGALE_LOG.Error("验证码获取失败!", zap.Error(err))
		response.FailWithMessage("验证码获取失败", c)
		return
	}
	response.OkWithDetailed(systemRes.SysCaptchaResponse{
		CaptchaId:     id,
		PicPath:       b64s,
		CaptchaLength: global.KUBEGALE_CONFIG.Captcha.KeyLong,
		OpenCaptcha:   oc,
	}, "验证码获取成功", c)
}

// 类型转换
func interfaceToInt(v interface{}) (i int) {
	switch v := v.(type) {
	case int:
		i = v
	default:
		i = 0
	}
	return
}
