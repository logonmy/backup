package main

import (
	"golang.org/x/net/websocket"
	"io/ioutil"
	"log"
	"sync"
)

// 只有一个连接，不处理 err，不关闭
func Read(ws *websocket.Conn) (data []byte, ok bool) {
	r, err := ws.NewFrameReader()
	if err != nil {
		return
	}
	fr, err := ws.HandleFrame(r)
	if err != nil {
		return
	}
	if fr == nil {
		return
	}
	data, err = ioutil.ReadAll(fr)
	if err != nil {
		return
	}
	return data, true
}

var (
	origin   = "https://ligaofeng.top"
	url      = "wss://ligaofeng.top:8849/chat"
	m        sync.Map
	msgNums  = 51
	readChan = make(chan string, 2000)
)

func Worker() {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Println("worker dial fail:", err)
		return
	}

	exit := make(chan struct{})

	go func() {
		for {
			select {
			case s := <-readChan:
				value, ok := m.Load(s)
				if ok {
					m.Store(s, value.(int)+1)
				} else {
					m.Store(s, 1)
				}
			case <-exit:
				return
			}
		}
	}()

	for {
		data, ok := Read(ws)
		if ok {
			s := string(data)
			if s == "end" {
				close(exit)
				return
			}
			readChan <- s
		}
	}
}

func main() {
	defer func() {
		num := 0
		all := 0
		m.Range(func(key, value interface{}) bool {
			log.Printf("%20s : %3d", key, value)
			num++
			if value == msgNums {
				all++
			}
			return true
		})
		// 接收到的worker数量
		log.Println("total worker: ", num)
		// 全部接收数量
		log.Println("all recv num: ", all)
	}()
	Worker()
}
