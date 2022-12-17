package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"syscall"

	"github.com/cloud-native-homework/httpserver/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// var versionPoint *string

var version string

var code int = 200

func main() {
	// 设置日志输出时间格式
	log.SetLevel(log.DebugLevel)

	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGTERM:
				log.Info("捕获 SIGTERM 不退出", s)
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT:
				log.Info("退出", s)
				ExitFunc()
			default:
				fmt.Println("other", s)
			}
		}
	}()

	// 获取启动参数中的 VERSION
	// versionPoint = flag.String("VERSION","0.1","define version")
	// flag.Parse()

	// 获取环境变量中的 VERSION
	version = os.Getenv("VERSION")
	log.Info("环境变量 VERSION: ", version)

	// 注册一个 prometheus 指标采集器
	metrics.Register()

	// 注册不同路径的处理函数
	http.HandleFunc("/", logging(rootHandler))
	http.HandleFunc("/healthz", logging(healthzHandler))
	// 增加 prometheus endpoint
	http.Handle("/metrics", promhttp.Handler())

	// 启动http server
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func ExitFunc() {
	fmt.Println("开始退出...")
	fmt.Println("执行清理...")
	fmt.Println("结束退出...")
	os.Exit(0)
}

// 根路径处理函数
func rootHandler(resp http.ResponseWriter, req *http.Request) {
	// 创建执行时间计时器
	timer := metrics.NewTimer()
	// 函数执行完成，采集器记录执行时间
	defer timer.ObserveTotal()

	// 增加随机延迟
	delay := randInt(10, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))

	// 输出响应
	// fmt.Fprintln(resp,"hello world")
	resp.Write([]byte("hello world\n"))
	io.WriteString(resp, "===================service2 Details of the http request header:============\n")
	for k, v := range req.Header {
		io.WriteString(resp, fmt.Sprintf("%s=%s\n", k, v))
	}
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

// healthz处理函数
func healthzHandler(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	from := req.FormValue("from")
	log.WithFields(log.Fields{
		"调用来源":  from,
		"当前响应码": code,
	}).Debugf("健康检查调用")

	// 返回指定的响应码
	resp.WriteHeader(code)
	io.WriteString(resp, "ok\n")
}

// 记录响应信息的结构体
type RespRecoder struct {
	http.ResponseWriter
	StatusCode int
}

// 重写 http.ResponseWriter WriteHeader 记录 statusCode
func (r *RespRecoder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// 记录请求日志
func logging(handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		// 获取客户端 IP
		clientIP := req.RemoteAddr
		// 获取请求路径
		path := req.URL.Path

		// 包装 http.ResponseWriter
		respRecoder := &RespRecoder{ResponseWriter: resp, StatusCode: 200}

		// 调用实际处理函数
		handler(respRecoder, req)

		// HTTP 返回码
		statusCode := respRecoder.StatusCode

		// 日志输出
		log.WithFields(log.Fields{
			"客户端地址": clientIP,
			"请求路径":  path,
			"响应码":   statusCode,
		}).Debug("收到请求")
	}
}
