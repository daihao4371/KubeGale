package aws

import (
	"KubeGale/global"
	model "KubeGale/model/cloudCmdb"
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"go.uber.org/zap"
)

type LoadBalancer struct {
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{}
}

// 获取 ALB/NLB 负载均衡器列表
func (l *LoadBalancer) getV2LoadBalancers(client *elasticloadbalancingv2.Client, pageSize int32) ([]types.LoadBalancer, error) {
	input := &elasticloadbalancingv2.DescribeLoadBalancersInput{
		PageSize: aws.Int32(pageSize),
	}

	response, err := client.DescribeLoadBalancers(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return response.LoadBalancers, nil
}

// 获取 Classic 负载均衡器列表
func (l *LoadBalancer) getClassicLoadBalancers(client *elasticloadbalancing.Client, pageSize int32) ([]elbtypes.LoadBalancerDescription, error) {
	input := &elasticloadbalancing.DescribeLoadBalancersInput{
		PageSize: aws.Int32(pageSize),
	}

	response, err := client.DescribeLoadBalancers(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return response.LoadBalancerDescriptions, nil
}

func (l *LoadBalancer) status(status string) string {
	if statusStr, ok := LoadBalancerStatus[status]; ok {
		return statusStr
	}
	return ""
}

func (l *LoadBalancer) getInstanceIP(instance interface{}, tp string) string {
	switch lb := instance.(type) {
	case types.LoadBalancer:
		// ALB/NLB
		if tp == "private" {
			for _, az := range lb.AvailabilityZones {
				if strings.Contains(*az.SubnetId, "private") {
					return *lb.DNSName
				}
			}
		} else if tp == "public" {
			for _, az := range lb.AvailabilityZones {
				if strings.Contains(*az.SubnetId, "public") {
					return *lb.DNSName
				}
			}
		}
	case elbtypes.LoadBalancerDescription:
		// Classic LB
		if tp == "private" && *lb.Scheme == "internal" {
			return *lb.DNSName
		} else if tp == "public" && *lb.Scheme == "internet-facing" {
			return *lb.DNSName
		}
	}
	return ""
}

func (l *LoadBalancer) List(cloudId uint, region model.CloudRegions, AccessKeyID, AccessKeySecret string) (list []model.LoadBalancer, err error) {
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

	// 创建 ALB/NLB 客户端
	v2Client := elasticloadbalancingv2.NewFromConfig(cfg)
	pageSize := int32(30)

	// 获取 ALB/NLB 负载均衡器
	v2LoadBalancers, err := l.getV2LoadBalancers(v2Client, pageSize)
	if err != nil {
		global.KUBEGALE_LOG.Error("aws get ALB/NLB instances fail!", zap.Error(err))
		return list, err
	}

	// 处理 ALB/NLB 负载均衡器
	for _, instance := range v2LoadBalancers {
		// 获取带宽信息（仅适用于 NLB）
		bandwidth := ""
		if instance.Type == types.LoadBalancerTypeEnumNetwork {
			// NLB 的带宽信息需要从其他 API 获取
			// TODO: 实现获取 NLB 带宽信息的逻辑
		}

		list = append(list, model.LoadBalancer{
			Name:            *instance.LoadBalancerName,
			InstanceId:      *instance.LoadBalancerArn,
			PrivateAddr:     l.getInstanceIP(instance, "private"),
			PublicAddr:      l.getInstanceIP(instance, "public"),
			Bandwidth:       bandwidth,
			Region:          strings.ReplaceAll(region.RegionId, "aws-", ""),
			RegionName:      region.RegionName,
			Status:          l.status(string(instance.State.Code)),
			CreationTime:    instance.CreatedTime.String(),
			CloudPlatformId: cloudId,
		})
	}

	// 创建 Classic LB 客户端
	classicClient := elasticloadbalancing.NewFromConfig(cfg)

	// 获取 Classic 负载均衡器
	classicLoadBalancers, err := l.getClassicLoadBalancers(classicClient, pageSize)
	if err != nil {
		global.KUBEGALE_LOG.Error("aws get Classic LB instances fail!", zap.Error(err))
		return list, err
	}

	// 处理 Classic 负载均衡器
	for _, instance := range classicLoadBalancers {
		// 获取带宽信息
		bandwidth := ""
		if instance.HealthCheck != nil {
			bandwidth = fmt.Sprintf("%d", *instance.HealthCheck.HealthyThreshold)
		}

		list = append(list, model.LoadBalancer{
			Name:            *instance.LoadBalancerName,
			InstanceId:      *instance.LoadBalancerName, // Classic LB 使用名称作为 ID
			PrivateAddr:     l.getInstanceIP(instance, "private"),
			PublicAddr:      l.getInstanceIP(instance, "public"),
			Bandwidth:       bandwidth,
			Region:          strings.ReplaceAll(region.RegionId, "aws-", ""),
			RegionName:      region.RegionName,
			Status:          l.status("active"), // Classic LB 总是返回 active 状态
			CreationTime:    instance.CreatedTime.String(),
			CloudPlatformId: cloudId,
		})
	}

	return list, nil
}
