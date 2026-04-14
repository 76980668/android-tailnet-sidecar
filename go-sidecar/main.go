package main

import (
	"log"
	"time"

	"tailscale.com/tsnet"
)

var server *tsnet.Server

//export Start
func Start(authKey string, hostname string) int {
	server = &tsnet.Server{
		Hostname: hostname,
		AuthKey:  authKey,
		Logf:     log.Printf,
	}

	go func() {
		for {
			err := server.Start()
			if err != nil {
				log.Println("start error:", err)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	select {}
}

//export Stop
func Stop() {
	if server != nil {
		server.Close()
	}
}

func main() {}
