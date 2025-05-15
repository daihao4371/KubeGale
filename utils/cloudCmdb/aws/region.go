package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Region AWS 区域结构体
type Region struct{}

// NewRegion 创建新的区域实例
func NewRegion() *Region {
	return &Region{}
}

// List 获取 AWS 区域列表
// AccessKeyID: 访问密钥ID
// AccessKeySecret: 访问密钥密码
// 返回: 区域列表和错误信息
func (r *Region) List(AccessKeyID, AccessKeySecret string) (list []types.Region, err error) {
	// 创建 AWS 配置
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"), // 使用默认区域
		config.WithCredentialsProvider(aws.NewCredentialsCache(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     AccessKeyID,
				SecretAccessKey: AccessKeySecret,
			}, nil
		}))),
	)
	if err != nil {
		return nil, err
	}

	// 创建 EC2 客户端
	client := ec2.NewFromConfig(cfg)

	// 创建请求
	input := &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(true), // 获取所有区域，包括未启用的
	}

	// 发送请求获取区域列表
	response, err := client.DescribeRegions(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return response.Regions, nil
}
