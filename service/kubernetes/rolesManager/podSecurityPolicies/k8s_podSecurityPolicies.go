package podSecurityPolicies

import (
	"KubeGale/global"
	"KubeGale/model/kubernetes/podSecurityPolicies"
	"KubeGale/utils/kubernetes"
	"KubeGale/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"strings"

	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type K8sPodSecurityPoliciesService struct {
	kubernetes.BaseService
}

func (k *K8sPodSecurityPoliciesService) GetPodSecurityPoliciesList(req podSecurityPolicies.GetPodSecurityPoliciesListReq, uuid uuid.UUID) (*[]v1.PodSecurityContext, int, error) {
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

	// 获取所有 Pod
	pods, err := client.CoreV1().Pods(req.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: req.LabelSelector,
	})
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}

	// 提取所有 Pod 的 SecurityContext
	var securityContexts []v1.PodSecurityContext
	for _, pod := range pods.Items {
		if pod.Spec.SecurityContext != nil {
			securityContexts = append(securityContexts, *pod.Spec.SecurityContext)
		}
	}

	// 关键字过滤
	var filteredContexts []v1.PodSecurityContext
	if req.Keyword != "" {
		for _, ctx := range securityContexts {
			// 这里可以根据需要添加更多的过滤条件
			if strings.Contains(ctx.String(), req.Keyword) {
				filteredContexts = append(filteredContexts, ctx)
			}
		}
	} else {
		filteredContexts = securityContexts
	}

	// 分页
	result, err := paginate.Paginate(filteredContexts, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	return result.(*[]v1.PodSecurityContext), len(filteredContexts), nil
}

func (k *K8sPodSecurityPoliciesService) DescribePodSecurityPolicies(req podSecurityPolicies.DescribePodSecurityPoliciesReq, uuid uuid.UUID) (*v1.PodSecurityContext, error) {
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

	// 获取指定 Pod
	pod, err := client.CoreV1().Pods(req.Namespace).Get(context.TODO(), req.PodSecurityPoliciesName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if pod.Spec.SecurityContext == nil {
		return nil, nil
	}

	return pod.Spec.SecurityContext, nil
}

func (k *K8sPodSecurityPoliciesService) UpdatePodSecurityPolicies(req podSecurityPolicies.UpdatePodSecurityPoliciesReq, uuid uuid.UUID) (*v1.PodSecurityContext, error) {
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

	// 获取当前 Pod
	pod, err := client.CoreV1().Pods(req.Namespace).Get(context.TODO(), req.PodSecurityPoliciesName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// 解析更新内容
	data, err := json.Marshal(req.Content)
	if err != nil {
		global.KUBEGALE_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}

	var securityContext v1.PodSecurityContext
	if err := json.Unmarshal(data, &securityContext); err != nil {
		global.KUBEGALE_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}

	// 更新 SecurityContext
	pod.Spec.SecurityContext = &securityContext

	// 更新 Pod
	updatedPod, err := client.CoreV1().Pods(req.Namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return updatedPod.Spec.SecurityContext, nil
}

func (k *K8sPodSecurityPoliciesService) DeletePodSecurityPolicies(req podSecurityPolicies.DeletePodSecurityPoliciesReq, uuid uuid.UUID) error {
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

	// 获取当前 Pod
	pod, err := client.CoreV1().Pods(req.Namespace).Get(context.TODO(), req.PodSecurityPoliciesName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// 清除 SecurityContext
	pod.Spec.SecurityContext = nil

	// 更新 Pod
	_, err = client.CoreV1().Pods(req.Namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	return err
}

func (k *K8sPodSecurityPoliciesService) CreatePodSecurityPolicies(req podSecurityPolicies.CreatePodSecurityPoliciesReq, uuid uuid.UUID) (*v1.PodSecurityContext, error) {
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

	// 解析创建内容
	data, err := json.Marshal(req.Content)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}

	var pod v1.Pod
	if err := json.Unmarshal(data, &pod); err != nil {
		global.KUBEGALE_LOG.Error("创建失败:" + err.Error())
		return nil, err
	}

	// 创建 Pod
	createdPod, err := client.CoreV1().Pods(req.Namespace).Create(context.TODO(), &pod, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return createdPod.Spec.SecurityContext, nil
}
