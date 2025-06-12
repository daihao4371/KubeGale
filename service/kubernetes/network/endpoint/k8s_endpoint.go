package endpoint

import (
	"KubeGale/global"
	"KubeGale/model/kubernetes/endpoint"
	"KubeGale/utils/kubernetes"
	"KubeGale/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"strings"

	"github.com/gofrs/uuid/v5"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type K8sEndPointService struct {
	kubernetes.BaseService
}

func (k *K8sEndPointService) GetEndPointList(req endpoint.GetEndPointListReq, uuid uuid.UUID) (*[]discoveryv1.EndpointSlice, int, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	options := metav1.ListOptions{
		LabelSelector: req.LabelSelector,
	}

	data, err := client.DiscoveryV1().EndpointSlices(req.Namespace).List(context.TODO(), options)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var filterEndPoints []discoveryv1.EndpointSlice
	if req.Keyword != "" {
		for _, item := range data.Items {
			if strings.Contains(item.Name, req.Keyword) {
				filterEndPoints = append(filterEndPoints, item)
			}
		}
	} else {
		filterEndPoints = data.Items
	}

	result, err := paginate.Paginate(filterEndPoints, req.Page, req.PageSize)

	return result.(*[]discoveryv1.EndpointSlice), len(filterEndPoints), nil
}

func (k *K8sEndPointService) DescribeEndPoint(req endpoint.DescribeEndPointReq, uuid uuid.UUID) (*discoveryv1.EndpointSlice, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	EndPointIns, err := client.DiscoveryV1().EndpointSlices(req.Namespace).Get(context.TODO(), req.EndPointName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return EndPointIns, nil
}

func (k *K8sEndPointService) UpdateEndPoint(req endpoint.UpdateEndPointReq, uuid uuid.UUID) (*discoveryv1.EndpointSlice, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.KUBEGALE_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.KUBEGALE_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	data, err := json.Marshal(req.Content)
	if err != nil {
		global.KUBEGALE_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	tmp := string(data)
	var EndPointIns *discoveryv1.EndpointSlice
	err = json.Unmarshal([]byte(tmp), &EndPointIns)
	if err != nil {
		global.KUBEGALE_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.DiscoveryV1().EndpointSlices(req.Namespace).Update(context.TODO(), EndPointIns, metav1.UpdateOptions{})
	if err != nil {
		global.KUBEGALE_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}

func (k *K8sEndPointService) DeleteEndPoint(req endpoint.DeleteEndPointReq, uuid uuid.UUID) error {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.KUBEGALE_LOG.Error("删除失败:" + err.Error())
		return err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.KUBEGALE_LOG.Error("删除失败:" + err.Error())
		return err
	}

	err = client.DiscoveryV1().EndpointSlices(req.Namespace).Delete(context.TODO(), req.EndPointName, metav1.DeleteOptions{})
	if err != nil {
		global.KUBEGALE_LOG.Error("删除失败:" + err.Error())
		return err
	}
	return nil
}

func (k *K8sEndPointService) CreateEndPoint(req endpoint.CreateEndPointReq, uuid uuid.UUID) (*discoveryv1.EndpointSlice, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.KUBEGALE_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	data, err := json.Marshal(req.Content)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	var EndPoint *discoveryv1.EndpointSlice
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &EndPoint)
	ins, err := client.DiscoveryV1().EndpointSlices(req.Namespace).Create(context.TODO(), EndPoint, metav1.CreateOptions{})
	if err != nil {
		global.KUBEGALE_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
