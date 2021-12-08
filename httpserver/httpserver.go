package main

import (
	"context"
	"flag"
	"fmt"
	"httpserver/metrics"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	glog.V(2).Info("Starting httpserver...")
	metrics.Register()

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/bar", rootHandler)
	serveMux.HandleFunc("/healthz", healthz)
	serveMux.Handle("/metrics", promhttp.Handler())

	srv := http.Server{
		Addr:    ":80",
		Handler: serveMux,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			glog.Fatalf("listen: %s\n", err)
		}
	}()
	glog.V(2).Info("htteserver Started")

	<-done
	glog.V(2).Info("httpserver Stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		glog.Fatalf("httpserver Shutdown Failed:%+v", err)
	}
	glog.V(2).Info("httpserver Exited Gracefully")

}

func healthz(w http.ResponseWriter, r *http.Request) {
	// 下面的日志会因为探活的周期时间设置每5秒出力一次，不出力到log文件。
	glog.V(5).Info("entering health handler, active= OK")
	w.Write([]byte("200"))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	glog.V(2).Info("entering root handler")

	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	user := r.URL.Query().Get("user")
	delay := randInt(10, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}

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
	w.Write([]byte("Hello, Welcome to kubernetes."))

	reqTime := time.Now().Format("2006-01-02 15:04:05")
	glog.V(2).Infof("[time: %s]-host: %s-method: %s-code: %d", reqTime, r.RemoteAddr, r.Method, http.StatusOK)
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
