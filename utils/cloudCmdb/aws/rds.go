package aws

import (
	"KubeGale/global"
	model "KubeGale/model/cloudCmdb"
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"go.uber.org/zap"
)

// RDS AWS RDS 结构体
type RDS struct {
}

// NewRDS 创建新的 RDS 实例
func NewRDS() *RDS {
	return &RDS{}
}

// status 转换 RDS 状态
// status: 原始状态
// 返回: 转换后的状态
func (r *RDS) status(status string) string {
	if statusStr, ok := RDSStatus[status]; ok {
		return statusStr
	}
	return ""
}

// get 获取 RDS 实例列表
// client: RDS 客户端
// pageSize: 每页数量
// 返回: RDS 实例列表和错误信息
func (r *RDS) get(client *rds.Client, pageSize int32) ([]types.DBInstance, error) {
	input := &rds.DescribeDBInstancesInput{
		MaxRecords: aws.Int32(pageSize),
	}

	response, err := client.DescribeDBInstances(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return response.DBInstances, nil
}

// List 获取 RDS 实例列表
// cloudId: 云平台ID
// region: 区域信息
// AccessKeyID: 访问密钥ID
// AccessKeySecret: 访问密钥密码
// 返回: RDS 实例列表和错误信息
func (r *RDS) List(cloudId uint, region model.CloudRegions, AccessKeyID, AccessKeySecret string) (list []model.RDS, err error) {
	// 创建 AWS 配置
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(strings.ReplaceAll(region.RegionId, "aws-", "")),
		config.WithCredentialsProvider(aws.NewCredentialsCache(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     AccessKeyID,
				SecretAccessKey: AccessKeySecret,
			}, nil
		}))),
	)
	if err != nil {
		global.KUBEGALE_LOG.Error("aws config load fail!", zap.Error(err))
		return nil, err
	}

	// 创建 RDS 客户端
	client := rds.NewFromConfig(cfg)
	pageSize := int32(30)

	// 获取 RDS 实例列表
	instances, err := r.get(client, pageSize)
	if err != nil {
		global.KUBEGALE_LOG.Error("aws get RDS instances fail!", zap.Error(err))
		return list, err
	}

	// 处理每个 RDS 实例
	for _, instance := range instances {
		// 获取连接地址
		privateAddr := ""
		if instance.Endpoint != nil && instance.Endpoint.Address != nil {
			privateAddr = *instance.Endpoint.Address
		}

		// 构建 RDS 实例信息
		list = append(list, model.RDS{
			Name:            *instance.DBInstanceIdentifier,
			InstanceId:      *instance.DBInstanceIdentifier,
			PrivateAddr:     privateAddr,
			Region:          strings.ReplaceAll(region.RegionId, "aws-", ""),
			RegionName:      region.RegionName,
			Status:          r.status(*instance.DBInstanceStatus),
			CreationTime:    instance.InstanceCreateTime.String(),
			CloudPlatformId: cloudId,
		})
	}

	return list, nil
}
