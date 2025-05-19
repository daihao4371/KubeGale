package aliyun

import (
	"KubeGale/global"
	model "KubeGale/model/cloudCmdb"
	"fmt"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"go.uber.org/zap"
)

type LoadBalancer struct {
	Type string // "slb", "alb"
}

func NewLoadBalancer(lbType string) *LoadBalancer {
	return &LoadBalancer{
		Type: lbType,
	}
}

// SLB相关方法
func (l *LoadBalancer) getSLB(client *slb.Client, pageNumber int, pageSize int) ([]slb.LoadBalancer, error) {
	request := slb.CreateDescribeLoadBalancersRequest()
	request.Scheme = "https"
	request.PageSize = requests.NewInteger(pageSize)
	request.PageNumber = requests.NewInteger(pageNumber)

	response, err := client.DescribeLoadBalancers(request)
	if err != nil {
		return nil, err
	}

	return response.LoadBalancers.LoadBalancer, err
}

// ALB相关方法
func (l *LoadBalancer) getALB(client *alb.Client, pageNumber int, pageSize int) ([]alb.LoadBalancer, error) {
	request := alb.CreateListLoadBalancersRequest()
	request.Scheme = "https"
	request.MaxResults = requests.NewInteger(pageSize)

	// 第一页不需要 NextToken
	if pageNumber > 1 {
		request.NextToken = strconv.Itoa((pageNumber - 1) * pageSize)
	}

	response, err := client.ListLoadBalancers(request)
	if err != nil {
		return nil, err
	}

	return response.LoadBalancers, err
}

func (l *LoadBalancer) status(status string) string {
	if _, ok := LoadBalancerStatus[status]; ok {
		return LoadBalancerStatus[status]
	}
	return ""
}

func (l *LoadBalancer) getInstanceIP(instance interface{}, tp string) string {
	switch v := instance.(type) {
	case slb.LoadBalancer:
		if v.AddressType == tp {
			return v.Address
		}
	case alb.LoadBalancer:
		if v.AddressType == tp {
			return v.DNSName
		}
	}
	return ""
}

func (l *LoadBalancer) List(cloudId uint, region model.CloudRegions, AccessKeyID, AccessKeySecret string) (list []model.LoadBalancer, err error) {
	config := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential(AccessKeyID, AccessKeySecret)
	regionId := strings.ReplaceAll(region.RegionId, "aliyun-", "")

	switch l.Type {
	case "slb":
		return l.listSLB(cloudId, region, regionId, config, credential)
	case "alb":
		return l.listALB(cloudId, region, regionId, config, credential)
	default:
		global.KUBEGALE_LOG.Error("Unsupported load balancer type", zap.String("type", l.Type))
		return nil, fmt.Errorf("unsupported load balancer type: %s", l.Type)
	}
}

func (l *LoadBalancer) listSLB(cloudId uint, region model.CloudRegions, regionId string, config *sdk.Config, credential *credentials.AccessKeyCredential) (list []model.LoadBalancer, err error) {
	client, err := slb.NewClientWithOptions(regionId, config, credential)
	if err != nil {
		global.KUBEGALE_LOG.Error("SLB new Client fail!", zap.Error(err))
		return
	}

	pageNumber := 1
	pageSize := 30

	for {
		response, err := l.getSLB(client, pageNumber, pageSize)
		if err != nil {
			global.KUBEGALE_LOG.Error("SLB getInstances fail!", zap.Error(err))
			return list, err
		}

		for _, instance := range response {
			bandwidth := ""
			if instance.Bandwidth != 0 {
				bandwidth = strconv.Itoa(instance.Bandwidth)
			}

			list = append(list, model.LoadBalancer{
				Name:            instance.LoadBalancerName,
				InstanceId:      instance.LoadBalancerId,
				PrivateAddr:     l.getInstanceIP(instance, "intranet"),
				PublicAddr:      l.getInstanceIP(instance, "internet"),
				Bandwidth:       bandwidth,
				Region:          regionId,
				RegionName:      region.RegionName,
				Status:          l.status(instance.LoadBalancerStatus),
				CreationTime:    instance.CreateTime,
				CloudPlatformId: cloudId,
				Type:            "SLB",
			})
		}

		if len(response) < pageSize {
			break
		}

		pageNumber++
	}

	return list, err
}

func (l *LoadBalancer) listALB(cloudId uint, region model.CloudRegions, regionId string, config *sdk.Config, credential *credentials.AccessKeyCredential) (list []model.LoadBalancer, err error) {
	client, err := alb.NewClientWithOptions(regionId, config, credential)
	if err != nil {
		global.KUBEGALE_LOG.Error("ALB new Client fail!", zap.Error(err))
		return
	}

	pageSize := 50
	var nextToken string

	for {
		request := alb.CreateListLoadBalancersRequest()
		request.Scheme = "https"
		request.MaxResults = requests.NewInteger(pageSize)

		// 只有在有 nextToken 时才设置
		if nextToken != "" {
			request.NextToken = nextToken
		}

		response, err := client.ListLoadBalancers(request)
		if err != nil {
			global.KUBEGALE_LOG.Error("ALB getInstances fail!", zap.Error(err))
			return list, err
		}

		for _, instance := range response.LoadBalancers {
			// 获取负载均衡器的详细信息
			detailRequest := alb.CreateGetLoadBalancerAttributeRequest()
			detailRequest.Scheme = "https"
			detailRequest.LoadBalancerId = instance.LoadBalancerId

			detailResponse, err := client.GetLoadBalancerAttribute(detailRequest)
			if err != nil {
				global.KUBEGALE_LOG.Error("ALB get detail fail!", zap.Error(err))
				continue
			}

			list = append(list, model.LoadBalancer{
				Name:            instance.LoadBalancerName,
				InstanceId:      instance.LoadBalancerId,
				PrivateAddr:     detailResponse.VpcId, // 使用 VPC ID 作为私网地址
				PublicAddr:      detailResponse.DNSName,
				Region:          regionId,
				RegionName:      region.RegionName,
				Status:          l.status(instance.LoadBalancerStatus),
				CreationTime:    instance.CreateTime,
				CloudPlatformId: cloudId,
				Type:            "ALB",
			})
		}

		// 检查是否有下一页
		if response.NextToken == "" {
			break
		}
		nextToken = response.NextToken
	}

	return list, err
}
