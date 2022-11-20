package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// var versionPoint *string

var version string

func main() {
	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("退出", s)
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
	log.Println(version)

	// 注册不同路径的处理函数
	http.HandleFunc("/", logging(rootHandler))
	http.HandleFunc("/healthz", logging(healthzHandler))
	// 启动http server
	err := http.ListenAndServe(":8080", nil)
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
	// 获取request header
	reqHeaders := req.Header
	// 写入response header
	resp.Header().Set("VERSION", version)
	for k, v := range reqHeaders {
		value := strings.Join(v, ",")
		resp.Header().Set(k, value)
	}
	// 获取参数
	// req.ParseForm()
	// codeStr := req.FormValue("code")
	// log.Println("codeStr")
	// code, _ := strconv.Atoi(codeStr)
	// 设置响应码，不指定默认返回200
	// resp.WriteHeader(code)
	// 输出响应
	// fmt.Fprintln(resp,"hello world")
	resp.Write([]byte("hello world\n"))
}

// healthz处理函数
func healthzHandler(resp http.ResponseWriter, req *http.Request) {
	// 返回200
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
		log.Printf("客户端地址: %s, 请求路径: %s, 响应码: %d", clientIP, path, statusCode)
	}
}