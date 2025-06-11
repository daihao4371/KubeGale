package metrics

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// MetricsCategory 指标分类结构体
// 用于定义不同类别的指标查询参数
type MetricsCategory struct {
	ClusterId uint   `json:"cluster_id"`          // 集群ID
	Category  string `json:"category"`            // 指标类别
	Nodes     string `json:"nodes,omitempty"`     // 节点列表
	PVC       string `json:"pvc,omitempty"`       // 持久卷声明
	Pods      string `json:"pods,omitempty"`      // Pod列表
	Ingress   string `json:"ingress,omitempty"`   // Ingress名称
	Selector  string `json:"selector,omitempty"`  // 选择器
	Namespace string `json:"namespace,omitempty"` // 命名空间
	Status    string `json:"status,omitempty"`    // 状态
	Start     int64  `json:"start,omitempty"`     // 开始时间
	End       int64  `json:"end,omitempty"`       // 结束时间
}

// bytesSent 生成Ingress发送字节数的Prometheus查询语句
func bytesSent(ingress string, namespace string, statuses string) string {
	return fmt.Sprintf(`sum(rate(nginx_ingress_controller_bytes_sent_sum{ingress="%s",namespace="%s",status=~"%s"}[1m])) by (ingress, namespace)`, ingress, namespace, statuses)
}

// GenerateQuery 根据指标类别生成Prometheus查询语句
func (mc *MetricsCategory) GenerateQuery() *PrometheusQuery {
	switch mc.Category {
	case "cluster":
		return &PrometheusQuery{
			MemoryUsage:               strings.Replace("sum(node_memory_MemTotal_bytes - (node_memory_MemFree_bytes + node_memory_Buffers_bytes + node_memory_Cached_bytes)) by (instance)", "_bytes", fmt.Sprintf("_bytes{instance=~\"%s\"}", mc.Nodes), -1),
			MemoryRequests:            fmt.Sprintf(`sum(kube_pod_container_resource_requests{node=~"%s", resource="memory"}) by (component)`, mc.Nodes),
			MemoryLimits:              fmt.Sprintf(`sum(kube_pod_container_resource_limits{node=~"%s", resource="memory"}) by (component)`, mc.Nodes),
			MemoryCapacity:            fmt.Sprintf(`sum(kube_node_status_capacity{node=~"%s", resource="memory"}) by (component)`, mc.Nodes),
			MemoryAllocatableCapacity: fmt.Sprintf(`sum(kube_node_status_allocatable{node=~"%s", resource="memory"}) by (component)`, mc.Nodes),
			CpuUsage:                  fmt.Sprintf(`sum(rate(node_cpu_seconds_total{instance=~"%s", mode=~"user|system"}[1m]))`, mc.Nodes),
			CpuRequests:               fmt.Sprintf(`sum(kube_pod_container_resource_requests{node=~"%s", resource="cpu"}) by (component)`, mc.Nodes),
			CpuLimits:                 fmt.Sprintf(`sum(kube_pod_container_resource_limits{node=~"%s", resource="cpu"}) by (component)`, mc.Nodes),
			CpuCapacity:               fmt.Sprintf(`sum(kube_node_status_capacity{node=~"%s", resource="cpu"}) by (component)`, mc.Nodes),
			CpuAllocatableCapacity:    fmt.Sprintf(`sum(kube_node_status_allocatable{node=~"%s", resource="cpu"}) by (component)`, mc.Nodes),
			PodUsage:                  fmt.Sprintf(`sum({__name__=~"kubelet_running_pod_count|kubelet_running_pods", node=~"%s"})`, mc.Nodes),
			PodCapacity:               fmt.Sprintf(`sum(kube_node_status_capacity{node=~"%s", resource="pods"}) by (component)`, mc.Nodes),
			PodAllocatableCapacity:    fmt.Sprintf(`sum(kube_node_status_allocatable{node=~"%s", resource="pods"}) by (component)`, mc.Nodes),
			FsSize:                    fmt.Sprintf(`sum(node_filesystem_size_bytes{instance=~"%s", mountpoint="/"}) by (kubernetes_node)`, mc.Nodes),
			FsUsage:                   fmt.Sprintf(`sum(node_filesystem_size_bytes{instance=~"%s", mountpoint="/"} - node_filesystem_avail_bytes{instance=~"%s", mountpoint="/"}) by (kubernetes_node)`, mc.Nodes, mc.Nodes),
		}
	case "nodes":
		return &PrometheusQuery{
			MemoryUsage:            `sum(node_memory_MemTotal_bytes - (node_memory_MemFree_bytes + node_memory_Buffers_bytes + node_memory_Cached_bytes)) by (instance)`,
			MemoryCapacity:         `sum(kube_node_status_capacity{resource="memory"}) by (node)`,
			MemoryRequests:         fmt.Sprintf(`sum(kube_pod_container_resource_requests{node=~"%s", resource="memory"}) by (node)`, mc.Nodes),
			CpuUsage:               `sum(rate(node_cpu_seconds_total{mode=~"user|system"}[1m])) by (instance)`,
			CpuRequests:            fmt.Sprintf(`sum(kube_pod_container_resource_requests{node=~"%s", resource="cpu"}) by (node)`, mc.Nodes),
			CpuCapacity:            fmt.Sprintf(`sum(kube_node_status_capacity{resource="cpu", node=~"%s"}) by (node)`, mc.Nodes),
			FsSize:                 fmt.Sprintf(`sum(node_filesystem_size_bytes{mountpoint="/", instance=~"%s"}) by (instance)`, mc.Nodes),
			FsUsage:                `sum(node_filesystem_size_bytes{mountpoint="/"} - node_filesystem_avail_bytes{mountpoint="/"}) by (instance)`,
			PodUsage:               fmt.Sprintf(`sum({__name__=~"kubelet_running_pod_count|kubelet_running_pods", node=~"%s"})`, mc.Nodes),
			PodCapacity:            fmt.Sprintf(`sum(kube_node_status_capacity{node=~"%s", resource="pods"}) by (node)`, mc.Nodes),
			PodAllocatableCapacity: fmt.Sprintf(`sum(kube_node_status_allocatable{node=~"%s", resource="pods"}) by (component)`, mc.Nodes),
		}
	case "pods":
		return &PrometheusQuery{
			CpuUsage:              fmt.Sprintf(`sum(rate(container_cpu_usage_seconds_total{image!="",pod="%s", namespace="%s"}[1m])) by (pod,namespace)`, mc.Pods, mc.Namespace),
			CpuRequests:           fmt.Sprintf(`sum(kube_pod_container_resource_requests{resource="cpu", pod=~"%s", namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			CpuLimits:             fmt.Sprintf(`sum(kube_pod_container_resource_limits{resource="cpu", pod=~"%s", namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			MemoryUsage:           fmt.Sprintf(`sum(container_memory_working_set_bytes{image!="", pod=~"%s", namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			MemoryRequests:        fmt.Sprintf(`sum(kube_pod_container_resource_requests{resource="memory", pod=~"%s", namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			MemoryLimits:          fmt.Sprintf(`sum(kube_pod_container_resource_limits{resource="memory", pod=~"%s", namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			FsUsage:               fmt.Sprintf(`sum(container_fs_usage_bytes{container!="POD",container!="",pod=~"%s", namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			FsWrite:               fmt.Sprintf(`sum(container_fs_writes_bytes_total{container!="", pod=~"%s", namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			FsRead:                fmt.Sprintf(`sum(container_fs_reads_bytes_total{container!="", pod="%s", namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			NetworkReceive:        fmt.Sprintf(`sum(container_network_receive_bytes_total{pod=~"%s",namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			NetworkTransmit:       fmt.Sprintf(`sum(container_network_transmit_bytes_total{pod=~"%s",namespace="%s"}) by (pod, namespace)`, mc.Pods, mc.Namespace),
			PodTcpEstablishedConn: fmt.Sprintf(`sum(inspector_pod_tcpsummarytcpestablishedconn{target_pod=~"%s", target_namespace="%s"} ) by  (target_pod, target_namespace)`, mc.Pods, mc.Namespace),
			PodTcpTimewaitConn:    fmt.Sprintf(`sum(inspector_pod_tcpsummarytcptimewaitconn{target_pod=~"%s", target_namespace="%s"} ) by  (target_pod, target_namespace)`, mc.Pods, mc.Namespace),
		}
	case "ingress":
		return &PrometheusQuery{
			BytesSentSuccess:        bytesSent(mc.Ingress, mc.Namespace, "^2\\\\d*"),
			BytesSent3XX:            bytesSent(mc.Ingress, mc.Namespace, "^3\\\\d*"),
			BytesSent4XX:            bytesSent(mc.Ingress, mc.Namespace, "^4\\\\d*"),
			BytesSentFailure:        bytesSent(mc.Ingress, mc.Namespace, "^5\\\\d*"),
			RequestDurationSeconds:  fmt.Sprintf(`sum(rate(nginx_ingress_controller_request_duration_seconds_sum{ingress="%s",namespace="%s"}[1m])) by (ingress, namespace)`, mc.Ingress, mc.Namespace),
			ResponseDurationSeconds: fmt.Sprintf(`sum(rate(nginx_ingress_controller_response_duration_seconds_sum{ingress="%s",namespace="%s"}[1m])) by (ingress, namespace)`, mc.Ingress, mc.Namespace),
		}
	}
	return nil
}

// MetricsQuery 指标查询结构体
// 用于定义各种类型的指标查询
type MetricsQuery struct {
	MemoryUsage               *MetricsCategory `json:"memoryUsage,omitempty"`               // 内存使用量
	MemoryRequests            *MetricsCategory `json:"memoryRequests,omitempty"`            // 内存请求量
	MemoryLimits              *MetricsCategory `json:"memoryLimits,omitempty"`              // 内存限制量
	MemoryCapacity            *MetricsCategory `json:"memoryCapacity,omitempty"`            // 内存容量
	MemoryAllocatableCapacity *MetricsCategory `json:"memoryAllocatableCapacity,omitempty"` // 可分配内存容量
	CpuUsage                  *MetricsCategory `json:"cpuUsage,omitempty"`                  // CPU使用量
	CpuLimits                 *MetricsCategory `json:"cpuLimits,omitempty"`                 // CPU限制量
	CpuRequests               *MetricsCategory `json:"cpuRequests,omitempty"`               // CPU请求量
	CpuCapacity               *MetricsCategory `json:"cpuCapacity,omitempty"`               // CPU容量
	CpuAllocatableCapacity    *MetricsCategory `json:"cpuAllocatableCapacity,omitempty"`    // 可分配CPU容量
	FsSize                    *MetricsCategory `json:"fsSize,omitempty"`                    // 文件系统大小
	FsUsage                   *MetricsCategory `json:"fsUsage,omitempty"`                   // 文件系统使用量
	FsWrite                   *MetricsCategory `json:"fsWrite,omitempty"`                   // 文件系统写入量
	FsRead                    *MetricsCategory `json:"fsRead,omitempty"`                    // 文件系统读取量
	PodUsage                  *MetricsCategory `json:"podUsage,omitempty"`                  // Pod使用量
	PodCapacity               *MetricsCategory `json:"podCapacity,omitempty"`               // Pod容量
	PodAllocatableCapacity    *MetricsCategory `json:"podAllocatableCapacity,omitempty"`    // 可分配Pod容量
	NetworkReceive            *MetricsCategory `json:"networkReceive,omitempty"`            // 网络接收量
	NetworkTransmit           *MetricsCategory `json:"networkTransmit,omitempty"`           // 网络发送量
	BytesSentSuccess          *MetricsCategory `json:"bytesSentSuccess,omitempty"`          // 成功发送字节数
	BytesSent3XX              *MetricsCategory `json:"bytesSent3XX,omitempty"`              // 3XX状态码发送字节数
	BytesSent4XX              *MetricsCategory `json:"bytesSent4XX,omitempty"`              // 4XX状态码发送字节数
	BytesSentFailure          *MetricsCategory `json:"bytesSentFailure,omitempty"`          // 失败发送字节数
	RequestDurationSeconds    *MetricsCategory `json:"requestDurationSeconds,omitempty"`    // 请求持续时间
	ResponseDurationSeconds   *MetricsCategory `json:"responseDurationSeconds,omitempty"`   // 响应持续时间
	WorkloadMemoryUsage       *MetricsCategory `json:"workloadMemoryUsage,omitempty"`       // 工作负载内存使用量
	PodTcpEstablishedConn     *MetricsCategory `json:"podTcpEstablishedConn,omitempty"`     // Pod TCP已建立连接数
	PodTcpTimewaitConn        *MetricsCategory `json:"podTcpTimewaitConn,omitempty"`        // Pod TCP等待连接数
}

// PrometheusQuery Prometheus查询结构体
// 用于存储Prometheus查询语句
type PrometheusQuery struct {
	CpuUsage                  string // CPU使用量查询
	CpuRequests               string // CPU请求量查询
	CpuLimits                 string // CPU限制量查询
	CpuCapacity               string // CPU容量查询
	WorkloadMemoryUsage       string // 工作负载内存使用量查询
	CpuAllocatableCapacity    string // 可分配CPU容量查询
	MemoryUsage               string // 内存使用量查询
	MemoryCapacity            string // 内存容量查询
	MemoryRequests            string // 内存请求量查询
	MemoryLimits              string // 内存限制量查询
	MemoryAllocatableCapacity string // 可分配内存容量查询
	FsUsage                   string // 文件系统使用量查询
	FsSize                    string // 文件系统大小查询
	FsWrite                   string // 文件系统写入量查询
	FsRead                    string // 文件系统读取量查询
	NetworkReceive            string // 网络接收量查询
	NetworkTransmit           string // 网络发送量查询
	PodUsage                  string // Pod使用量查询
	PodCapacity               string // Pod容量查询
	PodAllocatableCapacity    string // 可分配Pod容量查询
	DiskUsage                 string // 磁盘使用量查询
	DiskCapacity              string // 磁盘容量查询
	BytesSentSuccess          string // 成功发送字节数查询
	BytesSent3XX              string // 3XX状态码发送字节数查询
	BytesSent4XX              string // 4XX状态码发送字节数查询
	BytesSentFailure          string // 失败发送字节数查询
	RequestDurationSeconds    string // 请求持续时间查询
	ResponseDurationSeconds   string // 响应持续时间查询
	PodTcpEstablishedConn     string // Pod TCP已建立连接数查询
	PodTcpTimewaitConn        string // Pod TCP等待连接数查询
}

// GetValueByField 根据字段名获取查询语句
func (pq *PrometheusQuery) GetValueByField(field string) string {
	e := reflect.ValueOf(pq).Elem()
	for i := 0; i < e.NumField(); i++ {
		if e.Type().Field(i).Name == field {
			return e.Field(i).Interface().(string)
		}
	}
	return ""
}

// PrometheusQueryResp Prometheus查询响应结构体
type PrometheusQueryResp struct {
	Status string                   `json:"status"` // 响应状态
	Data   *PrometheusQueryRespData `json:"data"`   // 响应数据
}

// PrometheusQueryRespData Prometheus查询响应数据结构体
type PrometheusQueryRespData struct {
	ResultType string                      `json:"resultType"` // 结果类型
	Result     []PrometheusQueryRespResult `json:"result"`     // 查询结果
}

// PrometheusQueryRespResult Prometheus查询结果结构体
type PrometheusQueryRespResult struct {
	Metric interface{}   `json:"metric"` // 指标信息
	Values []interface{} `json:"values"` // 指标值
}

// PrometheusTracker Prometheus指标追踪器
// 用于存储和管理Prometheus查询结果
type PrometheusTracker struct {
	sync.RWMutex                                 // 读写锁
	Metrics      map[string]*PrometheusQueryResp // 指标查询结果映射
}

// NewPrometheusTracker 创建新的Prometheus追踪器
func NewPrometheusTracker() *PrometheusTracker {
	return &PrometheusTracker{Metrics: map[string]*PrometheusQueryResp{}}
}

// Get 获取指定键的查询结果
func (pt *PrometheusTracker) Get(key string) (*PrometheusQueryResp, bool) {
	pt.RLock()
	defer pt.RUnlock()
	val, ext := pt.Metrics[key]
	return val, ext
}

// Set 设置指定键的查询结果
func (pt *PrometheusTracker) Set(key string, val *PrometheusQueryResp) {
	pt.Lock()
	defer pt.Unlock()
	pt.Metrics[key] = val
}
