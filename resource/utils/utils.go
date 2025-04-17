package utils

import (
	"log"
	"os"
	"path/filepath"
)

// GetRootDir 获取项目根目录路径
func GetRootDir() string {
	// 当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前工作目录失败: %v", err)
	}

	// 向上查找直到找到包含config.yaml的目录
	for {
		// 检查当前目录是否包含config.yaml
		configPath := filepath.Join(currentDir, "config.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return currentDir
		}

		// 获取父目录
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// 已经到达根目录，无法继续向上
			log.Fatalf("无法找到项目根目录，请确保config.yaml文件存在")
		}
		currentDir = parentDir
	}
}
