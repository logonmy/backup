package rpc

import (

	"github.com/lgf133214/ProxyPool-Golang/model"
	"github.com/lgf133214/ProxyPool-Golang/storage"
	"net/rpc"
	"time"
)

func init() {
	go func() {
		for {
			select {
			case <-time.After(time.Minute * 5):
				for {
					if SendByRPC() {
						logger.Logger.Debug("send by RPC success one time")
						break
					} else {
						logger.Logger.Debug("send by RPC file one time, waiting...")
						time.Sleep(time.Second * 30)
					}
				}
			}
		}
	}()
}

type Request struct {
}

type Item struct {
	ProxyIp string
	Source  string
}

type Items struct {
	Items []Item
}

type ReceiveResponse struct {
}

func (r *Request) Receive(i Items, resp *ReceiveResponse) error {
	for _, i := range i.Items {
		storage.StoreToBufferDB(model.BufferItem{ProxyIp: i.ProxyIp, Source: i.Source})
	}
	return nil
}

func SendByRPC() bool {
	conn, err := rpc.DialHTTP("tcp", "213.136.64.204:8123")
	if err != nil {
		log.Println(err)
		return false
	}
	defer conn.Close()

	items, ok := storage.GetSuccessItemsForRPC()
	if !ok {
		return false
	}
	err = conn.Call("Request.Receive", items, new(ReceiveResponse))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
