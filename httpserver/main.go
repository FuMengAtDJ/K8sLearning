package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/golang/glog"
)

func main() {
	// 设置环境变量的值
	os.Setenv("VERSION", "1.15.6")
	// 设置flag for glog
	flag.Set("v", "4")
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "./log")
	flag.Parse()
	defer glog.Flush()

	glog.V(2).Info("Starting http server...")
	http.HandleFunc("/", rootHandler)
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", rootHandler)
	serveMux.HandleFunc("/healthz", healthz)
	err := http.ListenAndServe(":80", serveMux)
	if err != nil {
		log.Fatal(err)
	}

}

func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "200")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entering root handler")
	io.WriteString(w, "===================Details of the http request header:============\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}

	io.WriteString(w, "===================Details of the OS ENV:VERSION ============\n")
	version := os.Getenv("VERSION")
	if version != "" {
		io.WriteString(w, fmt.Sprintf("VERSION=%s\n", version))
	}

	w.WriteHeader(http.StatusOK)
	glog.V(2).Info(fmt.Sprintf("Client Ip=%s\n", RemoteIp(r)))
	glog.V(2).Info(fmt.Sprintf("Response http statuscode=%d\n", http.StatusOK))
}

func RemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("XRealIP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("XForwardedFor"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}
