package request

// NodeTTY 表示Kubernetes集群节点的TTY终端请求结构
type NodeTTY struct {
	ClusterId int    `json:"cluster_id"` // ClusterId 集群的唯一标识符
	NodeName  string `json:"node_name"`  // NodeName 目标节点的名称
}
