package main

import (
	"KubeGale/core"
	"KubeGale/global"
	"KubeGale/initialize"
	"KubeGale/resource/system"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

func main() {
	fmt.Println("=== KubeGale API 数据初始化工具 ===")

	// 添加命令行参数，用于控制是否自动退出
	autoExit := flag.Bool("auto-exit", true, "初始化完成后自动退出")
	flag.Parse()

	// 获取项目根目录的配置文件路径
	rootDir := getRootDir()
	configPath := filepath.Join(rootDir, "config.yaml")
	
	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("配置文件不存在: %s", configPath)
	}
	
	// 初始化配置，显式指定配置文件路径
	global.KUBEGALE_VP = core.Viper(configPath)
	global.KUBEGALE_LOG = core.Zap()
	zap.ReplaceGlobals(global.KUBEGALE_LOG)

	// 初始化数据库连接
	global.KUBEGALE_DB = initialize.Gorm()
	if global.KUBEGALE_DB == nil {
		log.Println("数据库连接失败，无法初始化API数据")
		os.Exit(1)
	}

	// 确保表结构存在
	initialize.RegisterTables()

	// 调用API初始化函数
	system.InitApiData()

	// 根据命令行参数决定是否自动退出
	if *autoExit {
		fmt.Println("初始化完成，程序自动退出")
		os.Exit(0)
	} else {
		fmt.Println("初始化完成，按任意键退出...")
		fmt.Scanln() // 等待用户按键
		os.Exit(0)
	}
}

// getRootDir 获取项目根目录路径
func getRootDir() string {
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
