package main

import (
	"fmt"
	"github.com/lgf133214/UnblockNeteaseMusic-GO/proxy"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	exit := make(chan bool, 1)
	go func() {
		sig := <-signalChan
		fmt.Println("\nreceive signal:", sig)
		exit <- true
	}()
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGSEGV)
	proxy.InitProxy()
	<-exit
	fmt.Println("exiting UnblockNeteaseMusic")
}
