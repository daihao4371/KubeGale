package velero

import (
	"KubeGale/global"
	veleroReq "KubeGale/model/velero/request"
	"KubeGale/utils/kubernetes"
	"KubeGale/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"strings"

	"github.com/gofrs/uuid/v5"
	v1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type K8sVeleroRestoresService struct {
	kubernetes.BaseService
}

var restoreGVR = schema.GroupVersionResource{
	Group:    "velero.io",
	Version:  "v1",
	Resource: "restores",
}

// GetK8sVeleroRestoreList
//
// @Description: 获取velero 恢复列表
// @receiver k8sVeleroRestoresService
// @param req query veleroReq.K8sVeleroRestoresSearchReq true "velero恢复列表请求参数"
// @param uuid path uuid.UUID true "用户UUID"
// @return list
// @return total
// @return err
func (k8sVeleroRestoresService *K8sVeleroRestoresService) GetK8sVeleroRestoreList(
	req veleroReq.K8sVeleroRestoresSearchReq, uuid uuid.UUID) (list *[]v1.Restore, total int, err error) {

	kubernetes, err := k8sVeleroRestoresService.Generic(&req, uuid)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	config, err := kubernetes.Config()
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	unstructuredList, err := dynamicClient.Resource(restoreGVR).Namespace(req.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	restoreList := &v1.RestoreList{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredList.UnstructuredContent(), restoreList)
	if err != nil {
		global.KUBEGALE_LOG.Error("转换失败:" + err.Error())
		return nil, 0, err
	}

	var filteredList []v1.Restore
	if req.Keyword != "" {
		for _, item := range restoreList.Items {
			if strings.Contains(item.Name, req.Keyword) {
				filteredList = append(filteredList, item)
			}
		}
	} else {
		filteredList = restoreList.Items
	}

	result, err := paginate.Paginate(filteredList, req.Page, req.PageSize)
	datas, _ := result.(*[]v1.Restore)

	return datas, len(filteredList), nil
}

// DescribeVeleroRestore
//
// @Description: 查看velero恢复详情
// @receiver K8sVeleroRestoresService
// @param req query veleroReq.DescribeVeleroRestoreReq true "velero恢复详情请求参数"
// @param uuid path uuid.UUID true "用户UUID"
// @return schedule
// @return err
func (K8sVeleroRestoresService *K8sVeleroRestoresService) DescribeVeleroRestore(req *veleroReq.DescribeVeleroRestoreReq, uuid uuid.UUID) (schedule *v1.Restore, err error) {
	kubernetes, err := K8sVeleroRestoresService.Generic(req, uuid)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	config, err := kubernetes.Config()
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	unstructured, err := dynamicClient.Resource(restoreGVR).Namespace(req.Namespace).Get(context.TODO(), req.VeleroRestoreName, metav1.GetOptions{})
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	restore := &v1.Restore{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured.UnstructuredContent(), restore)
	if err != nil {
		global.KUBEGALE_LOG.Error("转换失败:" + err.Error())
		return nil, err
	}
	return restore, nil
}

// DeleteK8sVeleroRestore
//
// @Description: 删除velero恢复记录
// @receiver K8sVeleroRestoresService
// @param req body veleroReq.DeleteVeleroRestoreReq true "删除velero恢复记录的请求参数"
// @param uuid path uuid.UUID true "用户UUID"
// @return err
func (K8sVeleroRestoresService *K8sVeleroRestoresService) DeleteK8sVeleroRestore(req *veleroReq.DeleteVeleroRestoreReq, uuid uuid.UUID) (err error) {
	kubernetes, err := K8sVeleroRestoresService.Generic(req, uuid)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建失败:" + err.Error())
		return err
	}
	config, err := kubernetes.Config()
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败:" + err.Error())
		return err
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败:" + err.Error())
		return err
	}
	err = dynamicClient.Resource(restoreGVR).Namespace(req.Namespace).Delete(context.TODO(), req.VeleroRestoreName, metav1.DeleteOptions{})
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败:" + err.Error())
		return err
	}
	return nil
}

// CreateK8sVeleroRestore
//
// @Description: 创建velero恢复
// @receiver K8sVeleroRestoresService
// @param req body veleroReq.CreateVeleroRestoreReq true "创建velero恢复的请求参数"
// @param uuid path uuid.UUID true "用户UUID"
// @return *v1.Restore
// @return error
func (K8sVeleroRestoresService *K8sVeleroRestoresService) CreateK8sVeleroRestore(req *veleroReq.CreateVeleroRestoreReq, uuid uuid.UUID) (*v1.Restore, error) {
	kubernetes, err := K8sVeleroRestoresService.Generic(req, uuid)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	config, err := kubernetes.Config()
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	data, err := json.Marshal(req.Content)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}
	var restore *v1.Restore
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &restore)

	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(restore)
	if err != nil {
		global.KUBEGALE_LOG.Error("转换失败:" + err.Error())
		return nil, err
	}

	created, err := dynamicClient.Resource(restoreGVR).Namespace("velero").Create(context.TODO(), &unstructured.Unstructured{Object: unstructuredMap}, metav1.CreateOptions{})
	if err != nil {
		global.KUBEGALE_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}

	result := &v1.Restore{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(created.UnstructuredContent(), result)
	if err != nil {
		global.KUBEGALE_LOG.Error("转换失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
