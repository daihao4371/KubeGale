package ws

// TerminalRequest 终端请求结构体
// 用于处理 Kubernetes Pod 终端连接的请求参数
type TerminalRequest struct {
	Name      string `json:"name"  form:"name"`            // 终端会话名称
	PodName   string `json:"pod_name" form:"pod_name"`     // Pod 名称
	Namespace string `json:"namespace" form:"namespace"`   // 命名空间
	ClusterId int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	XToken    string `json:"x-token" form:"x-token"`       // 认证令牌
	Cols      int    `json:"cols" form:"cols"`             // 终端列数
	Rows      int    `json:"rows" form:"rows"`             // 终端行数
}
