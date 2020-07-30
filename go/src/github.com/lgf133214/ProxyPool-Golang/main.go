package main

import (
	"github.com/lgf133214/ProxyPool-Golang/scheduler"
	"github.com/lgf133214/ProxyPool-Golang/server"
)

func main() {
	scheduler.Run()
	server.Run()
}
