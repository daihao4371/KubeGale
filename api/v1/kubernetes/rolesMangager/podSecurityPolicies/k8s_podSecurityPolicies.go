package podSecurityPolicies

import (
	"KubeGale/global"
	"KubeGale/model/common/request"
	"KubeGale/model/common/response"
	"KubeGale/model/kubernetes/podSecurityPolicies"
	"KubeGale/service"
	"KubeGale/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type K8sPodSecurityPoliciesApi struct{}

var k8sPodSecurityPoliciesService = service.ServiceGroupApp.RoleServiceGroup.PodSecurityPoliciesServiceGroup.K8sPodSecurityPoliciesService

// GetPodSecurityPoliciesList 获取 Pod 安全策略列表
// @Tags kubernetes
// @Summary 获取 Pod 安全策略列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query podSecurityPolicies.GetPodSecurityPoliciesListReq true "获取 Pod 安全策略列表"
// @Success 200 {object} response.Response{data=podSecurityPolicies.PodSecurityPoliciesListResponse,msg=string} "获取成功"
// @Router /kubernetes/podSecurityPolicies [get]
func (k *K8sPodSecurityPoliciesApi) GetPodSecurityPoliciesList(c *gin.Context) {
	req := podSecurityPolicies.GetPodSecurityPoliciesListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, total, err := k8sPodSecurityPoliciesService.GetPodSecurityPoliciesList(req, utils.GetUserUuid(c)); err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(podSecurityPolicies.PodSecurityPoliciesListResponse{
			Items: list,
			Total: total,
			PageInfo: request.PageInfo{
				Page:     req.Page,
				PageSize: req.PageSize,
				Keyword:  req.Keyword,
			},
		}, "获取成功", c)
	}
}

// DescribePodSecurityPoliciesInfo 获取 Pod 安全策略详情
// @Tags kubernetes
// @Summary 获取 Pod 安全策略详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query podSecurityPolicies.DescribePodSecurityPoliciesReq true "获取 Pod 安全策略详情"
// @Success 200 {object} response.Response{data=podSecurityPolicies.DescribePodSecurityPoliciesResponse,msg=string} "获取成功"
// @Router /kubernetes/podSecurityPolicies/details [get]
func (k *K8sPodSecurityPoliciesApi) DescribePodSecurityPoliciesInfo(c *gin.Context) {
	req := podSecurityPolicies.DescribePodSecurityPoliciesReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sPodSecurityPoliciesService.DescribePodSecurityPolicies(req, utils.GetUserUuid(c)); err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(podSecurityPolicies.DescribePodSecurityPoliciesResponse{Items: list}, "获取成功", c)
	}
}

// UpdatePodSecurityPolicies 更新 Pod 安全策略
// @Tags kubernetes
// @Summary 更新 Pod 安全策略
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body podSecurityPolicies.UpdatePodSecurityPoliciesReq true "更新 Pod 安全策略"
// @Success 200 {object} response.Response{data=v1.PodSecurityContext,msg=string} "更新成功"
// @Router /kubernetes/podSecurityPolicies [put]
func (k *K8sPodSecurityPoliciesApi) UpdatePodSecurityPolicies(c *gin.Context) {
	req := podSecurityPolicies.UpdatePodSecurityPoliciesReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, err := k8sPodSecurityPoliciesService.UpdatePodSecurityPolicies(req, utils.GetUserUuid(c)); err != nil {
		global.KUBEGALE_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
		return
	} else {
		response.OkWithDetailed(list, "更新成功", c)
	}
}

// DeletePodSecurityPolicies 删除 Pod 安全策略
// @Tags kubernetes
// @Summary 删除 Pod 安全策略
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body podSecurityPolicies.DeletePodSecurityPoliciesReq true "删除 Pod 安全策略"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /kubernetes/podSecurityPolicies [delete]
func (k *K8sPodSecurityPoliciesApi) DeletePodSecurityPolicies(c *gin.Context) {
	req := podSecurityPolicies.DeletePodSecurityPoliciesReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := k8sPodSecurityPoliciesService.DeletePodSecurityPolicies(req, utils.GetUserUuid(c)); err != nil {
		global.KUBEGALE_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
		return
	} else {
		time.Sleep(1 * time.Second)
		response.OkWithMessage("删除成功", c)
	}
}

// CreatePodSecurityPolicies 创建 Pod 安全策略
// @Tags kubernetes
// @Summary 创建 Pod 安全策略
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body podSecurityPolicies.CreatePodSecurityPoliciesReq true "创建 Pod 安全策略"
// @Success 200 {object} response.Response{data=v1.PodSecurityContext,msg=string} "创建成功"
// @Router /kubernetes/podSecurityPolicies [post]
func (k *K8sPodSecurityPoliciesApi) CreatePodSecurityPolicies(c *gin.Context) {
	req := podSecurityPolicies.CreatePodSecurityPoliciesReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if securityContext, err := k8sPodSecurityPoliciesService.CreatePodSecurityPolicies(req, utils.GetUserUuid(c)); err != nil {
		global.KUBEGALE_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败："+err.Error(), c)
		return
	} else {
		response.OkWithData(securityContext, c)
	}
}
