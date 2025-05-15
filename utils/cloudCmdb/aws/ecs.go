package aws

import (
	"KubeGale/global"
	model "KubeGale/model/cloudCmdb"
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"go.uber.org/zap"
)

type ECS struct {
}

func NewECS() *ECS {
	return &ECS{}
}

func (e *ECS) get(client *ec2.Client, pageSize int32) (*ec2.DescribeInstancesOutput, error) {
	input := &ec2.DescribeInstancesInput{
		MaxResults: aws.Int32(pageSize),
	}

	response, err := client.DescribeInstances(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (e *ECS) getInstanceIP(ip []string) string {
	if len(ip) == 0 {
		return ""
	}
	return ip[0]
}

func (e *ECS) status(status types.InstanceStateName) string {
	if statusStr, ok := EC2Status[string(status)]; ok {
		return statusStr
	}
	return ""
}

func (e *ECS) List(cloudId uint, region model.CloudRegions, AccessKeyID, AccessKeySecret string) (list []model.VirtualMachine, err error) {
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

	// 创建 EC2 客户端
	client := ec2.NewFromConfig(cfg)
	pageSize := int32(30)

	response, err := e.get(client, pageSize)
	if err != nil {
		global.KUBEGALE_LOG.Error("aws getInstances fail!", zap.Error(err))
		return list, err
	}

	for _, reservation := range response.Reservations {
		for _, instance := range reservation.Instances {
			// 获取实例名称（从标签中获取）
			instanceName := ""
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					instanceName = *tag.Value
					break
				}
			}

			// 获取私有IP和公有IP
			privateIP := ""
			if instance.PrivateIpAddress != nil {
				privateIP = *instance.PrivateIpAddress
			}

			publicIP := ""
			if instance.PublicIpAddress != nil {
				publicIP = *instance.PublicIpAddress
			}

			// 获取平台信息
			platformDetails := ""
			if instance.PlatformDetails != nil {
				platformDetails = *instance.PlatformDetails
			}

			platform := string(instance.Platform)

			vm := model.VirtualMachine{
				Name:            instanceName,
				InstanceId:      *instance.InstanceId,
				UserName:        "ec2-user", // AWS 默认用户名
				Password:        "changeme",
				Port:            "22",
				OS:              platformDetails,
				OSType:          platform,
				PrivateAddr:     privateIP,
				PublicAddr:      publicIP,
				Region:          strings.ReplaceAll(region.RegionId, "aws-", ""),
				RegionName:      region.RegionName,
				Status:          e.status(instance.State.Name),
				CreationTime:    instance.LaunchTime.String(),
				ExpiredTime:     "", // AWS 按需实例没有过期时间
				CloudPlatformId: cloudId,
			}

			list = append(list, vm)
		}
	}

	return list, nil
}
