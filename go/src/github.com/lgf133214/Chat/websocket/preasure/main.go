package main

import (
	"golang.org/x/net/websocket"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	origin = "https://ligaofeng.top"
	url    = "wss://ligaofeng.top:8849/chat"
	start  = make(chan struct{})
	// 计数
	done int32
	// 每个连接发送的信息数量 = msgNums + 1
	msgNums = 50
	wg      sync.WaitGroup
)

func Worker(id int) {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Println("worker ", id, " dial fail:", err)
		return
	}

	// 等待开始
	<-start

	send := 0

	defer func() {
		ws.Close()
		log.Printf("worker %3d done, send:%3d \n", id, send)
		wg.Done()
	}()

	for {
		msg := []byte("{\"msg\":\"worker " + strconv.Itoa(id) + "\"}")
		_, err := ws.Write(msg)
		if err != nil {
			return
		}
		send++
		// 自定义数量
		if send > msgNums {
			// 结束全部任务
			atomic.AddInt32(&done, 1)
			return
		}

		time.Sleep(time.Second)
	}
}

func main() {
	// 建立连接
	for i := range [70][0]int{} {
		go Worker(i)
		wg.Add(1)
	}
	// 开始发送
	close(start)
	// 等待发送任务完成
	wg.Wait()
	// 打印完成结果
	log.Println("done:", done)

	// 发送结束信息，可能不会收到
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		panic(err)
	}
	for range [5][0]int{} {
		msg := []byte("end")
		_, err := ws.Write(msg)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Millisecond * 500)
	}
}
