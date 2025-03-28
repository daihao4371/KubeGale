package system

import (
	"KubeGale/global"
	"KubeGale/model/system"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	AdminUsername    = "admin"
	AdminPassword    = "123456"
	AdminAccountType = 2
)

type UserMock struct {
	db *gorm.DB
	ce *casbin.Enforcer
}

func NewUserMock(db *gorm.DB, ce *casbin.Enforcer) *UserMock {
	return &UserMock{
		db: db,
		ce: ce,
	}
}

func (u *UserMock) CreateUserAdmin() error {
	// 检查是否已经初始化过用户
	var count int64
	u.db.Model(&system.User{}).Count(&count)
	if count > 0 {
		if global.KUBEGALE_LOG != nil {
			global.KUBEGALE_LOG.Info("[用户已经初始化过,跳过Mock]")
		} else {
			log.Println("[用户已经初始化过,跳过Mock]")
		}
		return nil
	}

	if global.KUBEGALE_LOG != nil {
		global.KUBEGALE_LOG.Info("[用户模块Mock开始]")
	} else {
		log.Println("[用户模块Mock开始]")
	}

	// 生成加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		if global.KUBEGALE_LOG != nil {
			global.KUBEGALE_LOG.Error("生成密码失败", zap.Error(err))
			global.KUBEGALE_LOG.Info("[用户模块Mock结束]")
		} else {
			log.Printf("生成密码失败: %v\n", err)
			log.Println("[用户模块Mock结束]")
		}
		return err
	}

	// 创建管理员用户实例
	adminUser := system.User{
		Username:    AdminUsername,
		Password:    string(hashedPassword),
		RealName:    "花海",
		AccountType: AdminAccountType,
	}

	// 使用 FirstOrCreate 方法查找或创建管理员用户
	result := u.db.Where("username = ?", adminUser.Username).FirstOrCreate(&adminUser)

	// 检查操作是否成功
	if result.Error != nil {
		if global.KUBEGALE_LOG != nil {
			global.KUBEGALE_LOG.Error("创建或获取管理员用户失败", zap.Error(result.Error))
			global.KUBEGALE_LOG.Info("[用户模块Mock结束]")
		} else {
			log.Printf("创建或获取管理员用户失败: %v\n", result.Error)
			log.Println("[用户模块Mock结束]")
		}
		return result.Error
	}

	// 根据 RowsAffected 判断用户是否已存在或新创建
	if result.RowsAffected == 1 {
		if global.KUBEGALE_LOG != nil {
			global.KUBEGALE_LOG.Info("管理员用户创建成功")
		} else {
			log.Println("管理员用户创建成功")
		}
	} else {
		if global.KUBEGALE_LOG != nil {
			global.KUBEGALE_LOG.Info("管理员用户已存在，跳过创建")
		} else {
			log.Println("管理员用户已存在，跳过创建")
		}
	}

	// 为管理员用户添加所有权限
	userIDStr := strconv.FormatInt(int64(adminUser.ID), 10)
	paths := []string{"/*"}
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}

	for _, path := range paths {
		for _, method := range methods {
			if ok, err := u.ce.AddPolicy(userIDStr, path, method); err == nil && ok {
				if global.KUBEGALE_LOG != nil {
					global.KUBEGALE_LOG.Info("成功添加权限策略",
						zap.String("用户", userIDStr),
						zap.String("路径", path),
						zap.String("方法", method))
				} else {
					log.Printf("成功添加权限策略: 用户=%s, 路径=%s, 方法=%s", userIDStr, path, method)
				}
			} else if err != nil {
				if global.KUBEGALE_LOG != nil {
					global.KUBEGALE_LOG.Error("添加权限策略失败", zap.Error(err))
				} else {
					log.Printf("添加权限策略失败: %v", err)
				}
			} else {
				if global.KUBEGALE_LOG != nil {
					global.KUBEGALE_LOG.Info("权限策略已存在",
						zap.String("用户", userIDStr),
						zap.String("路径", path),
						zap.String("方法", method))
				} else {
					log.Printf("权限策略已存在: 用户=%s, 路径=%s, 方法=%s", userIDStr, path, method)
				}
			}
		}
	}

	err = u.ce.SavePolicy()
	if err != nil {
		if global.KUBEGALE_LOG != nil {
			global.KUBEGALE_LOG.Error("保存策略失败", zap.Error(err))
		} else {
			log.Printf("保存策略失败: %v\n", err)
		}
		return err
	}

	menuIds := []int{
		1, 2, 3, 4, 5, 6, // 菜单ID
	}

	// 构建批量插入的数据
	userMenus := make([]map[string]interface{}, 0, len(menuIds))
	for _, menuId := range menuIds {
		userMenus = append(userMenus, map[string]interface{}{
			"user_id": adminUser.ID,
			"menu_id": menuId,
		})
	}

	// 先删除已有的关联
	if err := u.db.Table("user_menus").Where("user_id = ?", adminUser.ID).Delete(nil).Error; err != nil {
		if global.KUBEGALE_LOG != nil {
			global.KUBEGALE_LOG.Error("删除已有用户菜单关联失败", zap.Error(err))
		} else {
			log.Printf("删除已有用户菜单关联失败: %v", err)
		}
		return err
	}

	// 批量创建新的关联
	if err := u.db.Table("user_menus").Create(userMenus).Error; err != nil {
		if global.KUBEGALE_LOG != nil {
			global.KUBEGALE_LOG.Error("添加用户菜单关联失败", zap.Error(err))
		} else {
			log.Printf("添加用户菜单关联失败: %v", err)
		}
		return err
	}

	if global.KUBEGALE_LOG != nil {
		global.KUBEGALE_LOG.Info("[用户模块Mock结束]")
	} else {
		log.Println("[用户模块Mock结束]")
	}
	return nil
}

// InitUserData 提供一个可以在main中直接调用的函数
func InitUserData() {
	// 检查数据库连接是否已初始化
	if global.KUBEGALE_DB == nil {
		if global.KUBEGALE_LOG != nil {
			global.KUBEGALE_LOG.Error("错误: 数据库连接未初始化，无法初始化用户数据")
		} else {
			log.Println("错误: 数据库连接未初始化，无法初始化用户数据")
		}
		return
	}

	// 初始化casbin enforcer (将middleware中的逻辑移植到这里)
	e := initCasbinForResource(global.KUBEGALE_DB)
	if e == nil {
		if global.KUBEGALE_LOG != nil {
			global.KUBEGALE_LOG.Error("初始化casbin enforcer失败")
		} else {
			log.Println("初始化casbin enforcer失败")
		}
		return
	}

	userMock := NewUserMock(global.KUBEGALE_DB, e)
	if err := userMock.CreateUserAdmin(); err != nil {
		if global.KUBEGALE_LOG != nil {
			global.KUBEGALE_LOG.Error("初始化用户数据失败", zap.Error(err))
		} else {
			log.Printf("初始化用户数据失败: %v", err)
		}
	}
}

// initCasbinForResource 在resource包内初始化casbin，不依赖middleware
func initCasbinForResource(db *gorm.DB) *casbin.Enforcer {
	// 创建适配器
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		if global.KUBEGALE_LOG != nil {
			global.KUBEGALE_LOG.Error("创建Casbin适配器失败", zap.Error(err))
		} else {
			log.Printf("创建Casbin适配器失败: %v", err)
		}
		return nil
	}
	
	// 获取model.conf文件路径
	modelPath := getModelConfPath()
	if modelPath == "" {
		if global.KUBEGALE_LOG != nil {
			global.KUBEGALE_LOG.Error("无法找到model.conf文件")
		} else {
			log.Println("无法找到model.conf文件")
		}
		return nil
	}
	
	// 使用找到的model.conf文件创建enforcer
	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		if global.KUBEGALE_LOG != nil {
			global.KUBEGALE_LOG.Error("创建Casbin enforcer失败", zap.Error(err))
		} else {
			log.Printf("创建Casbin enforcer失败: %v", err)
		}
		return nil
	}
	
	return enforcer
}

// getModelConfPath 获取model.conf文件的绝对路径
func getModelConfPath() string {
	// 尝试直接在当前目录查找
	if _, err := os.Stat("model.conf"); err == nil {
		absPath, _ := filepath.Abs("model.conf")
		return absPath
	}
	
	// 尝试在项目根目录查找
	currentDir, err := os.Getwd()
	if err != nil {
		return ""
	}
	
	// 向上查找直到找到model.conf文件
	for {
		modelPath := filepath.Join(currentDir, "model.conf")
		if _, err := os.Stat(modelPath); err == nil {
			return modelPath
		}
		
		// 获取父目录
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// 已经到达根目录，无法继续向上
			return ""
		}
		currentDir = parentDir
	}
}
