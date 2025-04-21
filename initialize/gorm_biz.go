package initialize

import (
	"KubeGale/global"
	"KubeGale/model/cmdb"
)

func bizModel() error {
	db := global.KUBEGALE_DB
	err := db.AutoMigrate(
		cmdb.CmdbProjects{},
		cmdb.CmdbHosts{},
	)
	if err != nil {
		return err
	}
	return err
}
