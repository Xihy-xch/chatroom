package main

import (
	"log"
	"net/http"

	"github.com/Xihy-xch/tcp-chatroom/global"
	"github.com/Xihy-xch/tcp-chatroom/internal/server"
)

var (
	addr = ":2022"
	banner = "Xihy-chatroom: start on: %s"
)

func init() {
	global.Init()
}

func main() {
	log.Printf(banner + "\n", addr)

	server.RegisterHandle()

	log.Fatal(http.ListenAndServe(addr, nil))
}
