package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"unsafe"

	"tailscale.com/tsnet"
)

var server *tsnet.Server

//export startProxy
func startProxy(dataDir *C.char) {
	dir := C.GoString(dataDir)

	log.Println("Starting proxy, data dir:", dir)

	server = &tsnet.Server{
		Hostname: "android-sidecar",
		Dir:      dir,
		AuthKey:  "tskey-xxxx", // ⚠️ 换成你的
	}

	go run(server)
}

func run(s *tsnet.Server) {
	ctx := context.Background()

	// 等待网络 ready
	lc, err := s.LocalClient()
	if err != nil {
		log.Println("LocalClient error:", err)
		return
	}

	st, err := lc.Status(ctx)
	if err != nil {
		log.Println("Status error:", err)
		return
	}

	log.Println("Tailscale IPs:", st.Self.TailscaleIPs)

	// ========= 代理目标 =========
	target, _ := url.Parse("http://your-service.tailnet:8080")

	proxy := httputil.NewSingleHostReverseProxy(target)

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Proxy request:", r.URL.Path)
		proxy.ServeHTTP(w, r)
	})

	log.Println("Listening on 127.0.0.1:8081")

	err = http.ListenAndServe("127.0.0.1:8081", nil)
	if err != nil {
		log.Println("HTTP error:", err)
		os.Exit(1)
	}
}

func main() {}