package core

import (
	"KubeGale/core/internal"
	"KubeGale/global"
	"KubeGale/utils"

	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func Zap() (logger *zap.Logger) {
	if ok, _ := utils.PathExists(global.KUBEGALE_CONFIG.Zap.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", global.KUBEGALE_CONFIG.Zap.Director)
		_ = os.Mkdir(global.KUBEGALE_CONFIG.Zap.Director, os.ModePerm)
	}
	levels := global.KUBEGALE_CONFIG.Zap.Levels()
	length := len(levels)
	cores := make([]zapcore.Core, 0, length)
	for i := 0; i < length; i++ {
		core := internal.NewZapCore(levels[i])
		cores = append(cores, core)
	}
	logger = zap.New(zapcore.NewTee(cores...))
	if global.KUBEGALE_CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}
