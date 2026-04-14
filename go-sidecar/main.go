package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"sync"
	"tailscale.com/tsnet"
)

var (
	server *tsnet.Server
	once   sync.Once
	auth   string
)

// 初始化 auth
//export init_auth
func init_auth(c *C.char) {
	auth = C.GoString(c)
}

// 启动
//export start_sidecar
func start_sidecar() C.int {
	once.Do(func() {
		server = &tsnet.Server{
			AuthKey: auth,
		}
	})
	err := server.Start()
	if err != nil {
		return -1
	}
	return 0
}

// 获取 IP
//export get_ip
func get_ip() *C.char {
	if server == nil {
		return C.CString("")
	}
	ip := server.TailscaleIPs()
	if len(ip) == 0 {
		return C.CString("")
	}
	return C.CString(ip[0].String())
}

func main() {}
