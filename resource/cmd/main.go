package main

import (
	"KubeGale/core"
	"KubeGale/global"
	"KubeGale/initialize"
	resourceInit "KubeGale/resource/initialize"
	resourceUtils "KubeGale/resource/utils"
	"context"
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
	rootDir := resourceUtils.GetRootDir()
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

	// 创建上下文
	ctx := context.WithValue(context.Background(), "db", global.KUBEGALE_DB)

	// 初始化所有系统表和数据
	resourceInit.InitializeAllSystemData(ctx)

	// 根据命令行参数决定是否自动退出
	if *autoExit {
		fmt.Println("初始化完成，程序自动退出")
		os.Exit(0)
	}
}
