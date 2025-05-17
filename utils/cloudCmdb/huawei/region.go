package huawei

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	v3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/region"
)

type Region struct {
}

func NewRegion() *Region {
	return &Region{}
}

type Regions struct {
	Id       string
	Name     string
	RegionId string
}

func (r *Region) List(AccessKeyID, AccessKeySecret string) (list []Regions, err error) {
	// 创建认证信息
	credentials, err := global.NewCredentialsBuilder().
		WithAk(AccessKeyID).
		WithSk(AccessKeySecret).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	// 创建 IAM 客户端
	regionValue, err := region.SafeValueOf("cn-north-4")
	if err != nil {
		return nil, err
	}
	hcClient, err := v3.IamClientBuilder().
		WithRegion(regionValue).
		WithCredential(credentials).
		SafeBuild()
	if err != nil {
		return nil, err
	}
	client := v3.NewIamClient(hcClient)

	// 创建请求
	request := &model.KeystoneListRegionsRequest{}
	response, err := client.KeystoneListRegions(request)
	if err != nil {
		return nil, err
	}

	// 处理响应
	var regions []Regions
	for _, re := range *response.Regions {
		// 获取区域名称，优先使用中文名，如果没有则使用ID作为名称
		regionName := re.Id
		if re.Locales != nil && re.Locales.ZhCn != "" {
			regionName = re.Locales.ZhCn
		}

		regions = append(regions, Regions{
			Id:       re.Id,
			Name:     regionName,
			RegionId: re.Id,
		})
	}

	return regions, nil
}
