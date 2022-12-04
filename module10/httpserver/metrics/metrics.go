package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// 定一个常量，描述指标所属应用
	MetricsNamespace = "httpserver"
)

var (
	// 定义一个采集器
	functionLatency = CreateExecutionTimeMetric(MetricsNamespace,
		"Time spent.")
)

// 创建一个 prometheus 的直方图采集器，用于记录执行时间
func CreateExecutionTimeMetric(namespace string, help string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "execution_latency_seconds",
			Help:      help,
			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
		}, []string{"step"},
	)
}

// 将定义的采集器注册到 prometheus 客户端中
func Register() {
	err := prometheus.Register(functionLatency)
	if err != nil {
		fmt.Println(err)
	}
}

// 定义一个使用prometheus直方图记录执行时间的结构体
// usual usage pattern is: timer := NewExecutionTimer(...) ; compute ; timer.ObserveStep() ; ... ; timer.ObserveTotal()
type ExecutionTimer struct {
	histo *prometheus.HistogramVec
	start time.Time
	last  time.Time
}

// 计算执行时间，并将结果记入直方图中
func (t *ExecutionTimer) ObserveTotal() {
	(*t.histo).WithLabelValues("total").Observe(time.Now().Sub(t.start).Seconds())
}

// 使用传入的直方图创建一个从现在开始计时的新计时器 call ObserveXXX() on it to measure
func NewExecutionTimer(histo *prometheus.HistogramVec) *ExecutionTimer {
	now := time.Now()
	return &ExecutionTimer{
		histo: histo,
		start: now,
		last:  now,
	}
}

// 创建一个新计时器，使用上面定义的functionLatency直方图采集器记录指标
func NewTimer() *ExecutionTimer {
	return NewExecutionTimer(functionLatency)
}
