package system

import (
	"KubeGale/common"
	"context"
	"fmt"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderCasbin = initOrderApiIgnore + 1

type initCasbin struct{}

func NewInitCasbin() *initCasbin {
	return &initCasbin{}
}

func (i *initCasbin) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, common.ErrMissingDBContext
	}

	// 先创建基本表结构
	err := db.AutoMigrate(&adapter.CasbinRule{})
	if err != nil {
		return ctx, err
	}

	// 检查索引是否存在
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'casbin_rule' AND index_name = 'idx_casbin_rule'").Count(&count)

	// 如果索引不存在，则创建
	if count == 0 {
		err = db.Exec("ALTER TABLE casbin_rule ADD UNIQUE INDEX idx_casbin_rule (ptype, v0, v1, v2, v3, v4, v5)").Error
		if err != nil {
			return ctx, err
		}
	}

	// 修改表字符集为utf8mb4
	err = db.Exec("ALTER TABLE casbin_rule CONVERT TO CHARACTER SET utf8mb4").Error
	if err != nil {
		return ctx, err
	}

	// 设置AUTO_INCREMENT起始值
	err = db.Exec("ALTER TABLE casbin_rule AUTO_INCREMENT = 471").Error
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (i *initCasbin) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&adapter.CasbinRule{})
}

func (i initCasbin) InitializerName() string {
	var entity adapter.CasbinRule
	return entity.TableName()
}

// 修改InitializeData方法中的刷新Casbin权限API路径
func (i *initCasbin) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, common.ErrMissingDBContext
	}

	// 查询所有 API
	var apis []struct {
		Path     string
		Method   string
		ApiGroup string
	}
	if err := db.Table("sys_apis").Select("path, method, api_group").Find(&apis).Error; err != nil {
		return ctx, errors.Wrap(err, "查询所有API失败")
	}

	// 获取现有的权限规则
	var existingRules []adapter.CasbinRule
	if err := db.Find(&existingRules).Error; err != nil {
		return ctx, errors.Wrap(err, "获取现有权限规则失败")
	}

	// 创建权限规则映射，用于快速查找
	ruleMap := make(map[string]adapter.CasbinRule)
	for _, rule := range existingRules {
		key := fmt.Sprintf("%s:%s:%s", rule.Ptype, rule.V0, rule.V1)
		ruleMap[key] = rule
	}

	// 为不同角色生成API权限
	var newRules []adapter.CasbinRule

	// admin角色(888)拥有所有API权限
	for _, api := range apis {
		key := fmt.Sprintf("p:888:%s", api.Path)
		if _, exists := ruleMap[key]; !exists {
			newRules = append(newRules, adapter.CasbinRule{
				Ptype: "p",
				V0:    "888", // admin角色
				V1:    api.Path,
				V2:    api.Method,
			})
		}
	}

	// 开发负责人角色(9528)拥有除系统管理外的所有API权限
	for _, api := range apis {
		if api.ApiGroup != "系统用户" && api.ApiGroup != "角色" && api.ApiGroup != "菜单" {
			key := fmt.Sprintf("p:9528:%s", api.Path)
			if _, exists := ruleMap[key]; !exists {
				newRules = append(newRules, adapter.CasbinRule{
					Ptype: "p",
					V0:    "9528", // 开发负责人角色
					V1:    api.Path,
					V2:    api.Method,
				})
			}
		}
	}

	// 运维角色(8881)拥有基础API权限和即时通讯API权限
	basicApis := []string{
		"/user/getUserInfo",
		"/user/setSelfInfo",
		"/user/changePassword",
		"/menu/getMenu",
		"/menu/getBaseMenuTree",
		"/jwt/jsonInBlacklist",
		"/notification/getNotificationList",
		"/notification/getNotificationById",
		"/notification/getCardContent",
		"/notification/testNotification",
		"/cloud_region/tree",
		"/cloud_platform/list",
		"/virtualMachine/list",
		"/loadBalancer/list",
		"/rds/list",
	}
	for _, api := range apis {
		for _, basicApi := range basicApis {
			if api.Path == basicApi {
				key := fmt.Sprintf("p:8881:%s", api.Path)
				if _, exists := ruleMap[key]; !exists {
					newRules = append(newRules, adapter.CasbinRule{
						Ptype: "p",
						V0:    "8881", // 运维角色
						V1:    api.Path,
						V2:    api.Method,
					})
				}
				break
			}
		}
	}

	// 如果有新的权限规则需要添加
	if len(newRules) > 0 {
		if err := db.Create(&newRules).Error; err != nil {
			return ctx, errors.Wrap(err, "Casbin 表 ("+i.InitializerName()+") 数据初始化失败!")
		}
	}

	// 合并现有规则和新规则
	allRules := append(existingRules, newRules...)
	next := context.WithValue(ctx, i.InitializerName(), allRules)
	return next, nil
}

func (i *initCasbin) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	// 修改为检查 admin 角色(888)的权限
	if errors.Is(db.Where(adapter.CasbinRule{Ptype: "p", V0: "888", V1: "/user/getUserInfo", V2: "GET"}).
		First(&adapter.CasbinRule{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
