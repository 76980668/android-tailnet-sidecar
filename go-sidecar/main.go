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
	authKey string
)

//export Java_com_example_tailnet_TailnetBridge_initAuth
func Java_com_example_tailnet_TailnetBridge_initAuth(_, key *C.char) {
	authKey = C.GoString(key)
}

//export Java_com_example_tailnet_TailnetBridge_start
func Java_com_example_tailnet_TailnetBridge_start(_, _ *C.char) C.int {

	once.Do(func() {
		server = &tsnet.Server{
			AuthKey: authKey, // ⭐ 自动登录核心
		}
	})

	err := server.Start()
	if err != nil {
		return -1
	}

	return 0
}

//export Java_com_example_tailnet_TailnetBridge_getIP
func Java_com_example_tailnet_TailnetBridge_getIP(_, _ *C.char) *C.char {
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
