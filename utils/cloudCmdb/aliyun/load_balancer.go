package aliyun

import (
	"KubeGale/global"
	model "KubeGale/model/cloudCmdb"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"go.uber.org/zap"
)

// LoadBalancer 阿里云负载均衡器结构体
type LoadBalancer struct {
}

// NewLoadBalancer 创建新的负载均衡器实例
func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{}
}

// get 获取负载均衡器列表
// client: SLB客户端
// pageNumber: 页码
// pageSize: 每页数量
// 返回: 负载均衡器列表和错误信息
func (l *LoadBalancer) get(client *slb.Client, pageNumber int, pageSize int) ([]slb.LoadBalancer, error) {
	// 创建请求对象
	request := slb.CreateDescribeLoadBalancersRequest()
	request.Scheme = "https"
	request.PageSize = requests.NewInteger(pageSize)
	request.PageNumber = requests.NewInteger(pageNumber)

	// 发送请求获取负载均衡器列表
	response, err := client.DescribeLoadBalancers(request)
	if err != nil {
		return nil, err
	}

	return response.LoadBalancers.LoadBalancer, err
}

// status 转换负载均衡器状态
// status: 原始状态
// 返回: 转换后的状态
func (l *LoadBalancer) status(status string) string {
	if _, ok := LoadBalancerStatus[status]; ok {
		return LoadBalancerStatus[status]
	}
	return ""
}

// getInstanceIP 获取负载均衡器IP地址
// instance: 负载均衡器实例
// tp: IP类型（intranet/internet）
// 返回: IP地址
func (l *LoadBalancer) getInstanceIP(instance slb.LoadBalancer, tp string) string {
	// 检查是否为内网IP
	if instance.AddressType == "intranet" && instance.AddressType == tp {
		return instance.Address
	}

	// 检查是否为公网IP
	if instance.AddressType == "internet" && instance.AddressType == tp {
		return instance.Address
	}

	return ""
}

// List 获取负载均衡器列表
// cloudId: 云平台ID
// region: 区域信息
// AccessKeyID: 访问密钥ID
// AccessKeySecret: 访问密钥密码
// 返回: 负载均衡器列表和错误信息
func (l LoadBalancer) List(cloudId uint, region model.CloudRegions, AccessKeyID, AccessKeySecret string) (list []model.LoadBalancer, err error) {
	// 创建SDK配置
	config := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential(AccessKeyID, AccessKeySecret)

	// 创建SLB客户端
	client, err := slb.NewClientWithOptions(strings.ReplaceAll(region.RegionId, "aliyun-", ""), config, credential)
	if err != nil {
		global.KUBEGALE_LOG.Error("LoadBalancer new Client fail!", zap.Error(err))
		return
	}

	// 设置分页参数
	pageNumber := 1
	pageSize := 30

	// 循环获取所有负载均衡器
	for {
		// 获取当前页的负载均衡器列表
		response, err := l.get(client, pageNumber, pageSize)
		if err != nil {
			global.KUBEGALE_LOG.Error("LoadBalancer getInstances fail!", zap.Error(err))
			return list, err
		}

		// 处理每个负载均衡器实例
		for _, instance := range response {
			// 处理带宽信息
			bandwidth := ""
			if instance.Bandwidth != 0 {
				bandwidth = strconv.Itoa(instance.Bandwidth)
			}

			// 构建负载均衡器信息
			list = append(list, model.LoadBalancer{
				Name:            instance.LoadBalancerName,
				InstanceId:      instance.LoadBalancerId,
				PrivateAddr:     l.getInstanceIP(instance, "intranet"),
				PublicAddr:      l.getInstanceIP(instance, "internet"),
				Bandwidth:       bandwidth,
				Region:          strings.ReplaceAll(region.RegionId, "aliyun-", ""),
				RegionName:      region.RegionName,
				Status:          l.status(instance.LoadBalancerStatus),
				CreationTime:    instance.CreateTime,
				CloudPlatformId: cloudId,
			})
		}

		// 如果返回的实例数小于页大小，说明已经获取完所有数据
		if len(response) < pageSize {
			break
		}

		// 继续获取下一页
		pageNumber++
	}

	return list, err
}
