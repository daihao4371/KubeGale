package utils

import (
	"KubeGale/global"
	"os"
	"path/filepath"
)

// GetRootDir 获取当前项目根目录
func GetRootDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		global.KUBEGALE_LOG.Error("获取当前工作目录失败：%v")
	}

	// 向上查找直到找到包含conf.yaml的文件
	for {
		configPath := filepath.Join(currentDir, "config.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return currentDir
		}

		// 获取父级目录
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			global.KUBEGALE_LOG.Error("无法找到项目根目录，请确定config.yaml文件存在")
		}
		currentDir = parentDir
	}
}
