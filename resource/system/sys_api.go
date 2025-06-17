package system

import (
	"KubeGale/common"
	sysModel "KubeGale/model/system"
	"context"
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderApi = common.InitOrderSystem + 1

type InitApi struct{}

func (i InitApi) InitializerName() string {
	return sysModel.SysApi{}.TableName()
}

func (i *InitApi) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, common.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysApi{})
}

func (i *InitApi) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysApi{})
}

func (i *InitApi) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, common.ErrMissingDBContext
	}

	// 获取现有的API列表
	var existingApis []sysModel.SysApi
	if err := db.Find(&existingApis).Error; err != nil {
		return ctx, errors.Wrap(err, "获取现有API列表失败")
	}

	// 创建API映射，用于快速查找
	apiMap := make(map[string]sysModel.SysApi)
	for _, api := range existingApis {
		key := fmt.Sprintf("%s:%s:%s", api.ApiGroup, api.Method, api.Path)
		apiMap[key] = api
	}

	// 定义需要初始化的API列表
	apis := []sysModel.SysApi{
		{ApiGroup: "jwt", Method: "POST", Path: "/jwt/jsonInBlacklist", Description: "jwt加入黑名单(退出，必选)"},

		{ApiGroup: "系统用户", Method: "DELETE", Path: "/user/deleteUser", Description: "删除用户"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/admin_register", Description: "用户注册"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/getUserList", Description: "获取用户列表"},
		{ApiGroup: "系统用户", Method: "PUT", Path: "/user/setUserInfo", Description: "设置用户信息"},
		{ApiGroup: "系统用户", Method: "PUT", Path: "/user/setSelfInfo", Description: "设置自身信息(必选)"},
		{ApiGroup: "系统用户", Method: "GET", Path: "/user/getUserInfo", Description: "获取自身信息(必选)"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/setUserAuthorities", Description: "设置权限组"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/changePassword", Description: "修改密码（建议选择)"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/setUserAuthority", Description: "修改用户角色(必选)"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/resetPassword", Description: "重置用户密码"},
		{ApiGroup: "系统用户", Method: "PUT", Path: "/user/setSelfSetting", Description: "用户界面配置"},

		{ApiGroup: "api", Method: "POST", Path: "/api/createApi", Description: "创建api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/deleteApi", Description: "删除Api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/updateApi", Description: "更新Api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/getApiList", Description: "获取api列表"},
		{ApiGroup: "api", Method: "POST", Path: "/api/getAllApis", Description: "获取所有api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/getApiById", Description: "获取api详细信息"},
		{ApiGroup: "api", Method: "DELETE", Path: "/api/deleteApisByIds", Description: "批量删除api"},
		{ApiGroup: "api", Method: "GET", Path: "/api/syncApi", Description: "获取待同步API"},
		{ApiGroup: "api", Method: "GET", Path: "/api/getApiGroups", Description: "获取路由组"},
		{ApiGroup: "api", Method: "POST", Path: "/api/enterSyncApi", Description: "确认同步API"},
		{ApiGroup: "api", Method: "POST", Path: "/api/ignoreApi", Description: "忽略API"},

		{ApiGroup: "角色", Method: "POST", Path: "/authority/copyAuthority", Description: "拷贝角色"},
		{ApiGroup: "角色", Method: "POST", Path: "/authority/createAuthority", Description: "创建角色"},
		{ApiGroup: "角色", Method: "POST", Path: "/authority/deleteAuthority", Description: "删除角色"},
		{ApiGroup: "角色", Method: "PUT", Path: "/authority/updateAuthority", Description: "更新角色信息"},
		{ApiGroup: "角色", Method: "POST", Path: "/authority/getAuthorityList", Description: "获取角色列表"},
		{ApiGroup: "角色", Method: "POST", Path: "/authority/setDataAuthority", Description: "设置角色资源权限"},

		{ApiGroup: "casbin", Method: "POST", Path: "/casbin/updateCasbin", Description: "更改角色api权限"},
		{ApiGroup: "casbin", Method: "POST", Path: "/casbin/getPolicyPathByAuthorityId", Description: "获取权限列表"},

		{ApiGroup: "操作记录", Method: "POST", Path: "/sysOperationRecord/createSysOperationRecord", Description: "新增操作记录"},
		{ApiGroup: "操作记录", Method: "GET", Path: "/sysOperationRecord/findSysOperationRecord", Description: "根据ID获取操作记录"},
		{ApiGroup: "操作记录", Method: "GET", Path: "/sysOperationRecord/getSysOperationRecordList", Description: "获取操作记录列表"},
		{ApiGroup: "操作记录", Method: "DELETE", Path: "/sysOperationRecord/deleteSysOperationRecord", Description: "删除操作记录"},
		{ApiGroup: "操作记录", Method: "DELETE", Path: "/sysOperationRecord/deleteSysOperationRecordByIds", Description: "批量删除操作历史"},

		{ApiGroup: "按钮权限", Method: "POST", Path: "/authorityBtn/setAuthorityBtn", Description: "设置按钮权限"},
		{ApiGroup: "按钮权限", Method: "POST", Path: "/authorityBtn/getAuthorityBtn", Description: "获取已有按钮权限"},
		{ApiGroup: "按钮权限", Method: "POST", Path: "/authorityBtn/canRemoveAuthorityBtn", Description: "删除按钮"},

		{ApiGroup: "即时通讯", Method: "POST", Path: "/notification/createFeiShu", Description: "创建飞书通知"},
		{ApiGroup: "即时通讯", Method: "PUT", Path: "/notification/updateFeiShu", Description: "更新飞书通知"},
		{ApiGroup: "即时通讯", Method: "DELETE", Path: "/notification/deleteNotification", Description: "删除通知配置"},
		{ApiGroup: "即时通讯", Method: "POST", Path: "/notification/testNotification", Description: "测试通知发送"},
		{ApiGroup: "即时通讯", Method: "POST", Path: "/notification/createCardContent", Description: "创建卡片内容"},
		{ApiGroup: "即时通讯", Method: "PUT", Path: "/notification/updateCardContent", Description: "更新卡片内容"},

		{ApiGroup: "即时通讯", Method: "POST", Path: "/notification/createDingTalk", Description: "创建钉钉通知"},
		{ApiGroup: "即时通讯", Method: "PUT", Path: "/notification/updateDingTalk", Description: "更新钉钉通知"},

		{ApiGroup: "即时通讯", Method: "POST", Path: "/notification/getNotificationList", Description: "获取通知配置列表"},
		{ApiGroup: "即时通讯", Method: "GET", Path: "/notification/getNotificationById", Description: "根据ID获取通知配置"},
		{ApiGroup: "即时通讯", Method: "GET", Path: "/notification/getCardContent", Description: "根据通知ID获取卡片内容"},

		// CMDB项目管理
		{ApiGroup: "cmdb", Method: "POST", Path: "/cmdb/projects", Description: "新建项目"},
		{ApiGroup: "cmdb", Method: "DELETE", Path: "/cmdb/projects", Description: "删除项目"},
		{ApiGroup: "cmdb", Method: "DELETE", Path: "/cmdb/projectsByIds", Description: "批量删除项目"},
		{ApiGroup: "cmdb", Method: "PUT", Path: "/cmdb/projects", Description: "更新项目"},
		{ApiGroup: "cmdb", Method: "GET", Path: "/cmdb/projectsById", Description: "根据ID获取项目"},
		{ApiGroup: "cmdb", Method: "GET", Path: "/cmdb/projects", Description: "获取项目列表"},

		// CMDB主机管理
		{ApiGroup: "cmdb", Method: "POST", Path: "/cmdb/hosts", Description: "新建主机"},
		{ApiGroup: "cmdb", Method: "DELETE", Path: "/cmdb/hosts", Description: "删除主机"},
		{ApiGroup: "cmdb", Method: "DELETE", Path: "/cmdb/hostsByIds", Description: "批量删除主机"},
		{ApiGroup: "cmdb", Method: "PUT", Path: "/cmdb/hosts", Description: "更新主机"},
		{ApiGroup: "cmdb", Method: "POST", Path: "/cmdb/hosts/authentication", Description: "SSH认证主机"},
		{ApiGroup: "cmdb", Method: "POST", Path: "/cmdb/hosts/import", Description: "导入主机"},
		{ApiGroup: "cmdb", Method: "GET", Path: "/cmdb/hostsById", Description: "根据ID获取主机"},
		{ApiGroup: "cmdb", Method: "POST", Path: "/cmdb/hostsList", Description: "获取主机列表"},
		{ApiGroup: "cmdb", Method: "GET", Path: "/cmdb/hosts/terminal", Description: "SSH终端连接"},

		// 批量操作
		{ApiGroup: "cmdb", Method: "POST", Path: "/cmdb/batchOperations/execute", Description: "执行批量操作"},
		{ApiGroup: "cmdb", Method: "GET", Path: "/cmdb/batchOperations/execLogs", Description: "获取执行记录"},

		// 云平台管理
		{ApiGroup: "cloud_platform", Method: "POST", Path: "/cloud_platform/getById", Description: "获取云平台信息"},
		{ApiGroup: "cloud_platform", Method: "POST", Path: "/cloud_platform/create", Description: "创建云平台"},
		{ApiGroup: "cloud_platform", Method: "PUT", Path: "/cloud_platform/update", Description: "更新云平台"},
		{ApiGroup: "cloud_platform", Method: "DELETE", Path: "/cloud_platform/delete", Description: "删除云平台"},
		{ApiGroup: "cloud_platform", Method: "DELETE", Path: "/cloud_platform/deleteByIds", Description: "批量删除云平台"},
		{ApiGroup: "cloud_platform", Method: "POST", Path: "/cloud_platform/list", Description: "获取云平台列表"},

		// 云区域管理
		{ApiGroup: "cloud_region", Method: "POST", Path: "/cloud_region/syncRegion", Description: "同步区域信息"},
		{ApiGroup: "cloud_region", Method: "GET", Path: "/cloud_region/tree", Description: "获取区域树形结构"},

		// 云主机管理
		{ApiGroup: "virtualMachine", Method: "POST", Path: "/virtualMachine/sync", Description: "同步云主机"},
		{ApiGroup: "virtualMachine", Method: "POST", Path: "/virtualMachine/tree", Description: "获取云主机目录树"},
		{ApiGroup: "virtualMachine", Method: "POST", Path: "/virtualMachine/list", Description: "获取云主机列表"},

		// 负载均衡管理
		{ApiGroup: "loadBalancer", Method: "POST", Path: "/loadBalancer/sync", Description: "同步负载均衡"},
		{ApiGroup: "loadBalancer", Method: "POST", Path: "/loadBalancer/tree", Description: "获取负载均衡目录树"},
		{ApiGroup: "loadBalancer", Method: "POST", Path: "/loadBalancer/list", Description: "获取负载均衡列表"},

		// RDS管理
		{ApiGroup: "rds", Method: "POST", Path: "/rds/sync", Description: "同步RDS"},
		{ApiGroup: "rds", Method: "POST", Path: "/rds/tree", Description: "获取RDS目录树"},
		{ApiGroup: "rds", Method: "POST", Path: "/rds/list", Description: "获取RDS列表"},
		{ApiGroup: "rds", Method: "POST", Path: "/rds/get", Description: "RDS实例信息"},

		// 客户管理
		{ApiGroup: "customer", Method: "POST", Path: "/customer/customer", Description: "创建客户"},
		{ApiGroup: "customer", Method: "PUT", Path: "/customer/customer", Description: "更新客户"},
		{ApiGroup: "customer", Method: "DELETE", Path: "/customer/customer", Description: "删除客户"},
		{ApiGroup: "customer", Method: "GET", Path: "/customer/customer", Description: "获取单一客户信息"},
		{ApiGroup: "customer", Method: "GET", Path: "/customer/customerList", Description: "获取客户列表"},

		// 文件上传下载管理
		{ApiGroup: "fileUploadAndDownload", Method: "POST", Path: "/fileUploadAndDownload/upload", Description: "上传文件"},
		{ApiGroup: "fileUploadAndDownload", Method: "POST", Path: "/fileUploadAndDownload/getFileList", Description: "获取上传文件列表"},
		{ApiGroup: "fileUploadAndDownload", Method: "POST", Path: "/fileUploadAndDownload/deleteFile", Description: "删除指定文件"},
		{ApiGroup: "fileUploadAndDownload", Method: "POST", Path: "/fileUploadAndDownload/editFileName", Description: "编辑文件名或者备注"},
		{ApiGroup: "fileUploadAndDownload", Method: "POST", Path: "/fileUploadAndDownload/breakpointContinue", Description: "断点续传"},
		{ApiGroup: "fileUploadAndDownload", Method: "GET", Path: "/fileUploadAndDownload/findFile", Description: "查询当前文件成功的切片"},
		{ApiGroup: "fileUploadAndDownload", Method: "POST", Path: "/fileUploadAndDownload/breakpointContinueFinish", Description: "切片传输完成"},
		{ApiGroup: "fileUploadAndDownload", Method: "POST", Path: "/fileUploadAndDownload/removeChunk", Description: "删除切片"},

		// Kubernetes CloudTTY
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/cloudtty/get", Description: "CloudTTY终端连接"},
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/nodetty/get", Description: "NodeTTY终端连接"},

		// Kubernetes 集群管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/cluster", Description: "新建Kubernetes集群"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/cluster", Description: "删除Kubernetes集群"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/clusterByIds", Description: "批量删除Kubernetes集群"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/cluster", Description: "更新Kubernetes集群"},
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/credential", Description: "创建集群凭据"},
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/getUserById", Description: "获取集群用户信息"},
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/getClusterRoles", Description: "获取集群角色列表"},
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/getClusterApiGroups", Description: "获取集群资源分组"},
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/createClusterRole", Description: "创建集群角色"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/updateClusterRole", Description: "更新集群角色"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/deleteClusterRole", Description: "删除集群角色"},
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/createClusterUser", Description: "创建集群用户授权"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/updateClusterUser", Description: "更新集群用户授权"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/deleteClusterUser", Description: "删除集群用户授权"},
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/getClusterUserNamespace", Description: "获取用户命名空间"},
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/getClusterListNamespace", Description: "获取集群命名空间列表"},
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/clusterById", Description: "根据ID获取集群信息"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/clusterByName", Description: "根据名称获取集群信息"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/clusterList", Description: "获取集群列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/clusterPublic", Description: "获取公开集群信息"},

		// Kubernetes ConfigMap 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/configMap", Description: "创建ConfigMap"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/configMap", Description: "删除ConfigMap"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/configMap", Description: "更新ConfigMap"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/configMap", Description: "获取ConfigMap列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/configMapDetails", Description: "获取ConfigMap详细信息"},

		// Kubernetes HorizontalPod 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/horizontalPod", Description: "创建HorizontalPod"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/horizontalPod", Description: "删除HorizontalPod"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/horizontalPod", Description: "更新HorizontalPod"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/horizontalPod", Description: "获取HorizontalPod列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/horizontalPodDetail", Description: "获取HorizontalPod详细信息"},

		// Kubernetes LimitRange 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/limitRange", Description: "创建LimitRange"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/limitRange", Description: "删除LimitRange"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/limitRange", Description: "更新LimitRange"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/limitRange", Description: "获取LimitRange列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/limitRangeDetails", Description: "获取LimitRange详细信息"},

		// Kubernetes PodDisruptionBudget 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/Poddisruptionbudget", Description: "创建PodDisruptionBudget"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/Poddisruptionbudget", Description: "删除PodDisruptionBudget"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/Poddisruptionbudget", Description: "更新PodDisruptionBudget"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/Poddisruptionbudget", Description: "获取PodDisruptionBudget列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/PoddisruptionbudgetDetails", Description: "获取PodDisruptionBudget详细信息"},

		// Kubernetes ResourceQuota 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/ResourceQuotas", Description: "创建ResourceQuota"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/ResourceQuotas", Description: "删除ResourceQuota"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/ResourceQuotas", Description: "更新ResourceQuota"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/ResourceQuotas", Description: "获取ResourceQuota列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/ResourceQuotaDetails", Description: "获取ResourceQuota详细信息"},

		// Kubernetes Secret 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/secret", Description: "创建Secret"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/secret", Description: "删除Secret"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/secret", Description: "更新Secret"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/secret", Description: "获取Secret列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/secretDetails", Description: "获取Secret详细信息"},

		// Kubernetes Metrics 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/metrics/get", Description: "获取监控数据"},

		// Kubernetes Namespace 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/namespace", Description: "创建Namespace"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/namespace", Description: "删除Namespace"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/namespace", Description: "更新Namespace"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/namespace", Description: "获取Namespace列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/namespaceDetails", Description: "获取Namespace详细信息"},

		// Kubernetes Endpoint 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/endpoint", Description: "创建Endpoint"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/endpoint", Description: "删除Endpoint"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/endpoint", Description: "更新Endpoint"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/endpoint", Description: "获取Endpoint列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/endpointDetails", Description: "获取Endpoint详细信息"},

		// Kubernetes Ingress 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/ingress", Description: "创建Ingress"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/ingress", Description: "删除Ingress"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/ingress", Description: "更新Ingress"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/ingress", Description: "获取Ingress列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/ingressDetails", Description: "获取Ingress详细信息"},

		// Kubernetes Service 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/service", Description: "创建Service"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/service", Description: "删除Service"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/service", Description: "更新Service"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/service", Description: "获取Service列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/serviceDetails", Description: "获取Service详细信息"},

		// Kubernetes Node 管理
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/nodes", Description: "更新节点信息"},
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/nodes/EvictAllPod", Description: "驱逐节点所有Pod"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/nodes", Description: "获取节点列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/nodes/metrics", Description: "获取节点监控指标"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/nodeDetails", Description: "获取节点详细信息"},

		// Kubernetes Record 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/Record", Description: "创建记录"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/Record", Description: "删除记录"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/Record", Description: "更新记录"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/Record", Description: "获取记录列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/RecordDetails", Description: "获取记录详细信息"},

		// Kubernetes ClusterRole 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/ClusterRole", Description: "创建ClusterRole"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/ClusterRole", Description: "删除ClusterRole"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/ClusterRole", Description: "更新ClusterRole"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/ClusterRole", Description: "获取ClusterRole列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/ClusterRoleDetails", Description: "获取ClusterRole详细信息"},

		// Kubernetes ClusterRoleBinding 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/ClusterRoleBinding", Description: "创建ClusterRoleBinding"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/ClusterRoleBinding", Description: "删除ClusterRoleBinding"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/ClusterRoleBinding", Description: "更新ClusterRoleBinding"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/ClusterRoleBinding", Description: "获取ClusterRoleBinding列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/ClusterRoleBindingDetails", Description: "获取ClusterRoleBinding详细信息"},

		// Kubernetes Role 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/Role", Description: "创建Role"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/Role", Description: "删除Role"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/Role", Description: "更新Role"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/Role", Description: "获取Role列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/RoleDetails", Description: "获取Role详细信息"},

		// Kubernetes RoleBinding 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/RoleBinding", Description: "创建RoleBinding"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/RoleBinding", Description: "删除RoleBinding"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/RoleBinding", Description: "更新RoleBinding"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/RoleBinding", Description: "获取RoleBinding列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/RoleBindingDetails", Description: "获取RoleBinding详细信息"},

		// Kubernetes ServiceAccount 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/serviceAccount", Description: "创建ServiceAccount"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/serviceAccount", Description: "删除ServiceAccount"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/serviceAccount", Description: "更新ServiceAccount"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/serviceAccount", Description: "获取ServiceAccount列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/serviceAccountDetails", Description: "获取ServiceAccount详细信息"},

		// Kubernetes PV 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/pv", Description: "创建PersistentVolume"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/pv", Description: "删除PersistentVolume"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/pv", Description: "更新PersistentVolume"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/pv", Description: "获取PersistentVolume列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/pvDetails", Description: "获取PersistentVolume详细信息"},

		// Kubernetes PVC 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/pvc", Description: "创建PersistentVolumeClaim"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/pvc", Description: "删除PersistentVolumeClaim"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/pvc", Description: "更新PersistentVolumeClaim"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/pvc", Description: "获取PersistentVolumeClaim列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/pvcDetails", Description: "获取PersistentVolumeClaim详细信息"},

		// Kubernetes StorageClass 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/storageClass", Description: "创建StorageClass"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/storageClass", Description: "删除StorageClass"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/storageClass", Description: "更新StorageClass"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/storageClass", Description: "获取StorageClass列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/storageClassDetails", Description: "获取StorageClass详细信息"},

		// Kubernetes CronJob 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/cronJob", Description: "创建CronJob"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/cronJob", Description: "删除CronJob"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/cronJob", Description: "更新CronJob"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/cronJob", Description: "获取CronJob列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/cronJobDetails", Description: "获取CronJob详细信息"},

		// Kubernetes DaemonSet 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/daemonset", Description: "创建DaemonSet"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/daemonset", Description: "删除DaemonSet"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/daemonset", Description: "更新DaemonSet"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/daemonset", Description: "获取DaemonSet列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/daemonsetDetails", Description: "获取DaemonSet详细信息"},

		// Kubernetes Deployment 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/deployment", Description: "创建Deployment"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/deployment", Description: "删除Deployment"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/deployment", Description: "更新Deployment"},
		{ApiGroup: "kubernetes", Method: "PATCH", Path: "/kubernetes/deployment", Description: "回滚Deployment"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/deployment", Description: "获取Deployment列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/deployment/detail", Description: "获取Deployment详细信息"},

		// Kubernetes Job 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/job", Description: "创建Job"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/job", Description: "删除Job"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/job", Description: "更新Job"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/job", Description: "获取Job列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/jobDetails", Description: "获取Job详细信息"},

		// Kubernetes Pod 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/pods", Description: "创建Pod"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/pods", Description: "删除Pod"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/pods", Description: "更新Pod"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/pods", Description: "获取Pod列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/pods/metrics", Description: "获取Pod监控指标"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/podDetails", Description: "获取Pod详细信息"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/events", Description: "获取Pod事件"},
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/pods/listFiles", Description: "列出Pod文件"},
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/pods/uploadFile", Description: "上传文件到Pod"},
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/pods/deleteFiles", Description: "删除Pod文件"},

		// Kubernetes ReplicaSet 管理
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/replicaSet", Description: "获取ReplicaSet列表"},

		// Kubernetes StatefulSet 管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/statefulset", Description: "创建StatefulSet"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/statefulset", Description: "删除StatefulSet"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/statefulset", Description: "更新StatefulSet"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/statefulset", Description: "获取StatefulSet列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/statefulsetDetails", Description: "获取StatefulSet详细信息"},

		// Kubernetes Velero 备份恢复管理
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/velero/tasks", Description: "创建Velero任务"},
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/velero/record", Description: "创建Velero记录"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/velero/tasks", Description: "删除Velero任务"},
		{ApiGroup: "kubernetes", Method: "PUT", Path: "/kubernetes/velero/tasks", Description: "更新Velero任务"},
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/velero", Description: "创建Velero"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/velero/record", Description: "删除Velero记录"},
		{ApiGroup: "kubernetes", Method: "POST", Path: "/kubernetes/velero/record/reduction", Description: "还原Velero记录"},
		{ApiGroup: "kubernetes", Method: "DELETE", Path: "/kubernetes/velero/restore", Description: "删除Velero还原"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/velero/record", Description: "获取Velero记录列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/velero/recordDetail", Description: "获取Velero记录详细信息"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/velero/restore", Description: "获取Velero还原列表"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/velero/restoreDetail", Description: "获取Velero还原详细信息"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/velero/taskDetail", Description: "获取Velero任务详细信息"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/velero/tasks", Description: "获取Velero任务列表"},

		// Kubernetes WebSocket 管理
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/pods/terminal", Description: "Pod终端连接"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/pods/logs", Description: "获取Pod容器日志"},
		{ApiGroup: "kubernetes", Method: "GET", Path: "/kubernetes/pods/downloadFile", Description: "下载Pod文件"},
	}

	// 用于存储需要新增的API
	var newApis []sysModel.SysApi

	// 检查每个API是否需要新增
	for _, api := range apis {
		key := fmt.Sprintf("%s:%s:%s", api.ApiGroup, api.Method, api.Path)
		if _, exists := apiMap[key]; !exists {
			newApis = append(newApis, api)
		}
	}

	// 如果有新的API需要添加
	if len(newApis) > 0 {
		if err := db.Create(&newApis).Error; err != nil {
			return ctx, errors.Wrap(err, sysModel.SysApi{}.TableName()+"表数据初始化失败!")
		}
	}

	// 合并现有API和新API
	allApis := append(existingApis, newApis...)
	next := context.WithValue(ctx, i.InitializerName(), allApis)
	return next, nil
}
