package system

import (
	"KubeGale/global"
	"KubeGale/model/common/response"
	"KubeGale/model/system"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type SystemApiApi struct{}

func (s *SystemApiApi) CreateAPI(c *gin.Context) {
	var req system.CreateApiRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage("请求参数错误: "+err.Error(), c)
		return
	}

	// 参数基本验证
	if req.Name == "" {
		response.FailWithMessage("API名称不能为空", c)
		return
	}
	if req.Path == "" {
		response.FailWithMessage("API路径不能为空", c)
		return
	}
	if req.Method <= 0 || req.Method > 7 {
		response.FailWithMessage("无效的HTTP方法", c)
		return
	}

	// 构建API对象
	api := &system.Api{
		Name:        req.Name,
		Path:        req.Path,
		Method:      req.Method,
		Description: req.Description,
		Version:     req.Version,
		Category:    req.Category,
		IsPublic:    req.IsPublic,
	}

	// 调用service层创建API
	err = apiService.CreateApi(api)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建API失败", zap.Error(err))
		response.FailWithMessage("创建API失败: "+err.Error(), c)
		return
	}
	// 返回成功响应
	response.OkWithDetailed(gin.H{
		"id": api.ID,
	}, "创建API成功", c)
}

func (a *SystemApiApi) ListApis(c *gin.Context) {
	var req system.ListApisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数错误: "+err.Error(), c)
		return
	}

	// 调用service层获取API列表
	apis, total, err := apiService.ListApis(req.PageNumber, req.PageSize)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取API列表失败", zap.Error(err))
		response.FailWithMessage("获取API列表失败: "+err.Error(), c)
		return
	}
	// 返回API列表和总数
	response.OkWithDetailed(gin.H{
		"list":     apis,
		"total":    total,
		"page":     req.PageNumber,
		"pageSize": req.PageSize,
	}, "获取API列表成功", c)
}

// UpdateAPI 更新API信息
func (a *SystemApiApi) UpdateAPI(c *gin.Context) {
	var req system.UpdateApiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 构建更新的API对象
	api := &system.Api{
		ID:          req.ID,
		Name:        req.Name,
		Path:        req.Path,
		Method:      req.Method,
		Description: req.Description,
		Version:     req.Version,
		Category:    req.Category,
		IsPublic:    req.IsPublic,
	}
	err := apiService.UpdateApi(api)
	if err != nil {
		response.FailWithMessage("更新API失败", c)
		global.KUBEGALE_LOG.Error("更新API失败", zap.Error(err))
		return
	}
	response.OkWithMessage("更新API成功", c)
}

// DeleteAPI 删除API
func (a *SystemApiApi) DeleteAPI(c *gin.Context) {
	// 从URL参数中获取API ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage("无效的API ID参数", c)
		return
	}

	// 验证ID
	if id <= 0 {
		response.FailWithMessage("API ID必须大于0", c)
		return
	}

	// 调用service层删除API
	if err := apiService.DeleteApi(id); err != nil {
		global.KUBEGALE_LOG.Error("删除API失败", zap.Error(err), zap.Int("id", id))
		response.FailWithMessage("删除API失败: "+err.Error(), c)
		return
	}

	// 返回成功响应
	response.OkWithMessage("删除API成功", c)
}
