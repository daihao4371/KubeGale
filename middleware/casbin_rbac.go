package middleware

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"log"
)

func InitCasbin(db *gorm.DB) *casbin.Enforcer {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalf("Failed to create adapter: %v", err)
	}
	enforcer, err := casbin.NewEnforcer("./model.conf", adapter)
	if err != nil {
		log.Fatalf("Failed to create enforcer: %v", err)
	}
	return enforcer
}
