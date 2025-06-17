package system

import (
	"KubeGale/global"
	"KubeGale/model/system"
	"errors"
)

type AuthorityApiService struct{}

var AuthorityApiServiceApp = new(AuthorityApiService)

// SetApisForAuthority 为角色分配API权限（覆盖式）
func (s *AuthorityApiService) SetApisForAuthority(authorityId uint, apiIds []uint) error {
	// 1. 删除原有权限
	if err := global.KUBEGALE_DB.Where("authority_id = ?", authorityId).Delete(&system.SysAuthorityApi{}).Error; err != nil {
		return err
	}
	// 2. 批量插入新权限
	var records []system.SysAuthorityApi
	for _, apiId := range apiIds {
		records = append(records, system.SysAuthorityApi{
			AuthorityId: authorityId,
			ApiId:       apiId,
		})
	}
	if len(records) > 0 {
		if err := global.KUBEGALE_DB.Create(&records).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetApisByAuthority 查询角色拥有的API权限
func (s *AuthorityApiService) GetApisByAuthority(authorityId uint) ([]system.SysApi, error) {
	var apis []system.SysApi
	err := global.KUBEGALE_DB.Joins("JOIN sys_authority_apis ON sys_authority_apis.api_id = sys_apis.id").
		Where("sys_authority_apis.authority_id = ?", authorityId).
		Find(&apis).Error
	return apis, err
}

// UserHasApi 判断用户是否有某API权限
func (s *AuthorityApiService) UserHasApi(userId uint, path, method string) (bool, error) {
	// 1. 查询用户所有角色
	var user system.SysUser
	if err := global.KUBEGALE_DB.Preload("Authorities").Where("id = ?", userId).First(&user).Error; err != nil {
		return false, err
	}
	if len(user.Authorities) == 0 {
		return false, errors.New("用户未分配任何角色")
	}
	// 2. 查询这些角色拥有的API权限
	var count int64
	err := global.KUBEGALE_DB.Table("sys_authority_apis").
		Joins("JOIN sys_apis ON sys_authority_apis.api_id = sys_apis.id").
		Where("sys_authority_apis.authority_id IN ? AND sys_apis.path = ? AND sys_apis.method = ?", user.Authorities, path, method).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
