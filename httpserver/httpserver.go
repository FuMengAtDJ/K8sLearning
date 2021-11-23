package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/golang/glog"
)

func main() {
	// 设置日志等级
	flag.Set("v", "4")
	flag.Set("alsologtostderr", "true")
	// 默认log出力到/var/log/container下面。便于对log进行统一收集。不再单独出力到下面的目录
	// flag.Set("log_dir", "./log")
	// 解析传入的参数
	flag.Parse()
	// 退出前执行，清空缓存区，将日志写入文件
	defer glog.Flush()

	// 设置日志文件大小，超出时分割出新的文件,默认时1.8GB。
	glog.MaxSize = 10 * 1024 * 1024 // 10MB
	glog.V(2).Info("Starting http server...")
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", rootHandler)
	serveMux.HandleFunc("/healthz", healthz)
	err := http.ListenAndServe(":80", serveMux)
	if err != nil {
		glog.Fatal(err)
	}

}

func healthz(w http.ResponseWriter, r *http.Request) {
	// 下面的日志会因为探活的周期时间设置每5秒出力一次，不出力到log文件。
	glog.V(5).Info("entering health handler, active= OK")
	w.Write([]byte("200"))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	glog.V(2).Info("entering root handler")
	// add request header to response header.
	for k, v := range r.Header {
		w.Header()[k] = v
	}

	// add env to esponse header
	version := os.Getenv("ENV_VERSION")
	if version != "" {
		w.Header().Set("VERSION", version)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, Welcome to kubernetes.¥n"))

	reqTime := time.Now().Format("2006-01-02 15:04:05")
	glog.V(2).Infof("[time: %s]-host: %s-method: %s-code: %d", reqTime, r.RemoteAddr, r.Method, http.StatusOK)
}
