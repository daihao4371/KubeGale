package cloudCmdb

import (
	"fmt"
	model "KubeGale/model/cloudCmdb"
	"KubeGale/model/common/request"
	"KubeGale/global"


)


type CloudPlatformService struct{}

// List  厂商列表
func (p *CloudPlatformService) List(cloud model.CloudPlatform, info request.PageInfo, order string, desc bool) (list []model.CloudPlatform, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.KUBEGALE_DB.Model(model.CloudPlatform{})
	var cloudList []model.CloudPlatform

	if cloud.Name != "" {
		db = db.Where("name LIKE ?", "%"+cloud.Name+"%")
	}

	err = db.Count(&total).Error

	if err != nil {
		return cloudList, total, err
	} else {
		db = db.Limit(limit).Offset(offset)
		if order != "" {
			var OrderStr string
			// 设置有效排序key 防止sql注入
			orderMap := make(map[string]bool, 5)
			orderMap["id"] = true
			orderMap["name"] = true
			if orderMap[order] {
				if desc {
					OrderStr = order + " desc"
				} else {
					OrderStr = order
				}
			} else { // didn't matched any order key in `orderMap`
				err = fmt.Errorf("非法的排序字段: %v", order)
				return cloudList, total, err
			}

			err = db.Order(OrderStr).Find(&cloudList).Error
		} else {
			err = db.Order("id").Find(&cloudList).Error
		}
	}

	return cloudList, total, err
}
