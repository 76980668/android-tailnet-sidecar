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

	"tailscale.com/tsnet"
)

var server *tsnet.Server

//export start
func start(user *C.char, authKey *C.char) {

	u := C.GoString(user)
	k := C.GoString(authKey)

	log.Println("start tailnet:", u)

	server = &tsnet.Server{
		Hostname: "android-sidecar",
		AuthKey:  k,
	}

	go run(server)
}

func run(s *tsnet.Server) {

	lc, err := s.LocalClient()
	if err != nil {
		log.Println("LocalClient error:", err)
		return
	}

	st, err := lc.Status(context.Background())
	if err != nil {
		log.Println("status error:", err)
		return
	}

	log.Println("node:", st.Self.HostName)

	// proxy target
	target, _ := url.Parse("http://100.64.0.1:8080")

	proxy := httputil.NewSingleHostReverseProxy(target)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	log.Println("proxy running :8081")
	http.ListenAndServe("127.0.0.1:8081", nil)
}

func main() {}
