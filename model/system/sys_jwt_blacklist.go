package system

import (
	"KubeGale/global"
)

type JwtBlacklist struct {
	global.KUBEGALE_MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
