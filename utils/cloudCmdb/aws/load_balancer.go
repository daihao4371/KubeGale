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
	global.KUBEGALE_LOG.Info("开始同步AWS负载均衡器",
		zap.Uint("cloudId", cloudId),
		zap.String("region", region.RegionId),
		zap.String("regionName", region.RegionName))

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
		global.KUBEGALE_LOG.Error("AWS配置加载失败",
			zap.Error(err),
			zap.String("region", region.RegionId))
		return nil, fmt.Errorf("AWS配置加载失败: %v", err)
	}

	// 创建 ALB/NLB 客户端
	v2Client := elasticloadbalancingv2.NewFromConfig(cfg)
	pageSize := int32(30)

	// 获取 ALB/NLB 负载均衡器
	v2LoadBalancers, err := l.getV2LoadBalancers(v2Client, pageSize)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取ALB/NLB实例失败",
			zap.Error(err),
			zap.String("region", region.RegionId))
		return list, fmt.Errorf("获取ALB/NLB实例失败: %v", err)
	}

	global.KUBEGALE_LOG.Info("成功获取ALB/NLB实例",
		zap.Int("count", len(v2LoadBalancers)),
		zap.String("region", region.RegionId))

	// 处理 ALB/NLB 负载均衡器
	for _, instance := range v2LoadBalancers {
		global.KUBEGALE_LOG.Debug("处理ALB/NLB实例",
			zap.String("name", aws.ToString(instance.LoadBalancerName)),
			zap.String("type", string(instance.Type)))

		// 获取带宽信息（仅适用于 NLB）
		bandwidth := ""
		if instance.Type == types.LoadBalancerTypeEnumNetwork {
			// 获取 NLB 的属性信息
			attributes, err := v2Client.DescribeLoadBalancerAttributes(context.TODO(), &elasticloadbalancingv2.DescribeLoadBalancerAttributesInput{
				LoadBalancerArn: instance.LoadBalancerArn,
			})
			if err != nil {
				global.KUBEGALE_LOG.Error("aws get NLB attributes fail!", zap.Error(err))
			} else {
				// 遍历属性查找带宽信息
				for _, attr := range attributes.Attributes {
					if *attr.Key == "load_balancing.cross_zone.enabled" {
						// 如果是跨可用区，带宽会更高
						if *attr.Value == "true" {
							bandwidth = "1000" // 跨可用区 NLB 默认带宽为 1000Mbps
						} else {
							bandwidth = "500" // 单可用区 NLB 默认带宽为 500Mbps
						}
						break
					}
				}
			}
		}

		privateAddr := ""
		publicAddr := ""
		if instance.Scheme == types.LoadBalancerSchemeEnumInternetFacing {
			publicAddr = aws.ToString(instance.DNSName)
		} else if instance.Scheme == types.LoadBalancerSchemeEnumInternal {
			privateAddr = aws.ToString(instance.DNSName)
		}

		lb := model.LoadBalancer{
			Name:            aws.ToString(instance.LoadBalancerName),
			InstanceId:      aws.ToString(instance.LoadBalancerArn),
			PrivateAddr:     privateAddr,
			PublicAddr:      publicAddr,
			Bandwidth:       bandwidth,
			Region:          strings.ReplaceAll(region.RegionId, "aws-", ""),
			RegionName:      region.RegionName,
			Status:          l.status(string(instance.State.Code)),
			CreationTime:    instance.CreatedTime.String(),
			CloudPlatformId: cloudId,
		}

		global.KUBEGALE_LOG.Debug("添加ALB/NLB到列表",
			zap.String("name", lb.Name),
			zap.String("instanceId", lb.InstanceId))

		list = append(list, lb)
	}

	// 创建 Classic LB 客户端
	classicClient := elasticloadbalancing.NewFromConfig(cfg)

	// 获取 Classic 负载均衡器
	classicLoadBalancers, err := l.getClassicLoadBalancers(classicClient, pageSize)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取Classic LB实例失败",
			zap.Error(err),
			zap.String("region", region.RegionId))
		return list, fmt.Errorf("获取Classic LB实例失败: %v", err)
	}

	global.KUBEGALE_LOG.Info("成功获取Classic LB实例",
		zap.Int("count", len(classicLoadBalancers)),
		zap.String("region", region.RegionId))

	// 处理 Classic 负载均衡器
	for _, instance := range classicLoadBalancers {
		global.KUBEGALE_LOG.Debug("处理Classic LB实例",
			zap.String("name", aws.ToString(instance.LoadBalancerName)))

		// 获取带宽信息
		bandwidth := ""
		if instance.HealthCheck != nil {
			bandwidth = fmt.Sprintf("%d", *instance.HealthCheck.HealthyThreshold)
		}

		privateAddr := ""
		publicAddr := ""
		if instance.Scheme != nil && *instance.Scheme == "internet-facing" {
			publicAddr = aws.ToString(instance.DNSName)
		} else if instance.Scheme != nil && *instance.Scheme == "internal" {
			privateAddr = aws.ToString(instance.DNSName)
		}

		lb := model.LoadBalancer{
			Name:            aws.ToString(instance.LoadBalancerName),
			InstanceId:      aws.ToString(instance.LoadBalancerName),
			PrivateAddr:     privateAddr,
			PublicAddr:      publicAddr,
			Bandwidth:       bandwidth,
			Region:          strings.ReplaceAll(region.RegionId, "aws-", ""),
			RegionName:      region.RegionName,
			Status:          l.status("active"),
			CreationTime:    instance.CreatedTime.String(),
			CloudPlatformId: cloudId,
		}

		global.KUBEGALE_LOG.Debug("添加Classic LB到列表",
			zap.String("name", lb.Name),
			zap.String("instanceId", lb.InstanceId))

		list = append(list, lb)
	}

	global.KUBEGALE_LOG.Info("负载均衡器同步完成",
		zap.Int("totalCount", len(list)),
		zap.String("region", region.RegionId))

	return list, nil
}
