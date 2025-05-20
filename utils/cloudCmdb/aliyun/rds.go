package aliyun

import (
	"KubeGale/global"
	model "KubeGale/model/cloudCmdb"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"go.uber.org/zap"
)

type RDS struct{}

func NewRDS() *RDS {
	return &RDS{}
}

func (r *RDS) status(status string) string {
	if _, ok := RdsStatus[status]; ok {
		return RdsStatus[status]
	}
	return ""
}

func (r *RDS) get(client *rds.Client, pageNumber int, pageSize int) ([]rds.DBInstance, error) {
	request := rds.CreateDescribeDBInstancesRequest()
	request.Scheme = "https"
	request.PageSize = requests.NewInteger(pageSize)
	request.PageNumber = requests.NewInteger(pageNumber)
	request.Engine = "MySQL" // 只获取MySQL实例
	request.SetReadTimeout(30 * time.Second)
	request.SetConnectTimeout(10 * time.Second)

	response, err := client.DescribeDBInstances(request)
	if err != nil {
		return nil, fmt.Errorf("获取RDS实例列表失败: %v", err)
	}

	return response.Items.DBInstance, nil
}

func (r *RDS) List(cloudId uint, region model.CloudRegions, AccessKeyID, AccessKeySecret string) (list []model.RDS, err error) {
	// 验证 AccessKey
	if AccessKeyID == "" || AccessKeySecret == "" {
		return nil, fmt.Errorf("AccessKey 不能为空")
	}

	// 获取区域ID
	regionId := strings.ReplaceAll(region.RegionId, "aliyun-", "")

	// 创建配置
	config := sdk.NewConfig()
	config.EnableAsync = true
	config.MaxRetryTime = 3
	config.Timeout = 30 // 设置超时时间为30秒
	config.AutoRetry = true
	config.MaxRetryTime = 3

	// 创建凭证
	credential := credentials.NewAccessKeyCredential(AccessKeyID, AccessKeySecret)

	// 创建客户端
	client, err := rds.NewClientWithOptions(regionId, config, credential)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建阿里云RDS客户端失败",
			zap.String("regionId", regionId),
			zap.Error(err))
		return nil, fmt.Errorf("创建阿里云RDS客户端失败: %v", err)
	}

	// 验证 AccessKey 是否有效
	request := rds.CreateDescribeRegionsRequest()
	request.Scheme = "https"
	request.SetReadTimeout(30 * time.Second)
	request.SetConnectTimeout(10 * time.Second)

	_, err = client.DescribeRegions(request)
	if err != nil {
		if strings.Contains(err.Error(), "InvalidAccessKeyId.NotFound") {
			return nil, fmt.Errorf("AccessKey 无效或不存在")
		}
		if strings.Contains(err.Error(), "timeout") {
			return nil, fmt.Errorf("请求超时，请检查网络连接")
		}
		return nil, fmt.Errorf("验证 AccessKey 失败: %v", err)
	}

	pageNumber := 1
	pageSize := 30

	for {
		response, err := r.get(client, pageNumber, pageSize)
		if err != nil {
			global.KUBEGALE_LOG.Error("获取RDS实例列表失败",
				zap.String("regionId", regionId),
				zap.Int("pageNumber", pageNumber),
				zap.Error(err))
			return list, err
		}

		for _, instance := range response {
			// 获取实例的连接信息
			connectionRequest := rds.CreateDescribeDBInstanceNetInfoRequest()
			connectionRequest.Scheme = "https"
			connectionRequest.DBInstanceId = instance.DBInstanceId
			connectionRequest.SetReadTimeout(30 * time.Second)
			connectionRequest.SetConnectTimeout(10 * time.Second)

			connectionResponse, err := client.DescribeDBInstanceNetInfo(connectionRequest)
			if err != nil {
				global.KUBEGALE_LOG.Error("获取RDS实例连接信息失败",
					zap.String("instanceId", instance.DBInstanceId),
					zap.Error(err))
				continue
			}

			// 获取连接地址
			var privateAddr, publicAddr string
			for _, netInfo := range connectionResponse.DBInstanceNetInfos.DBInstanceNetInfo {
				if netInfo.IPType == "Private" {
					privateAddr = netInfo.ConnectionString
				} else if netInfo.IPType == "Public" {
					publicAddr = netInfo.ConnectionString
				}
			}

			list = append(list, model.RDS{
				Name:            instance.DBInstanceDescription,
				InstanceId:      instance.DBInstanceId,
				PrivateAddr:     privateAddr,
				PublicAddr:      publicAddr,
				Region:          regionId,
				RegionName:      region.RegionName,
				Status:          r.status(instance.DBInstanceStatus),
				CreationTime:    instance.CreateTime,
				CloudPlatformId: cloudId,
			})
		}

		if len(response) < pageSize {
			break
		}

		pageNumber++
	}

	global.KUBEGALE_LOG.Info("成功获取RDS实例列表",
		zap.String("regionId", regionId),
		zap.Int("count", len(list)))

	return list, nil
}
