package replicaSet

import "KubeGale/model/common/request"

// GetReplicaSetListReq 获取副本集列表请求结构体
type GetReplicaSetListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetReplicaSetListReq) GetClusterID() int {
	return r.ClusterId
}
