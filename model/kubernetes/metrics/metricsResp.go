package metrics

// MetricsResponse 指标响应结构体
// 用于返回指标查询的结果
type MetricsResponse struct {
	Metrics interface{} `json:"metrics"` // 指标数据
}
