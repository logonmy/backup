package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/unrolled/secure"
	"log"
	"net/http"
	"sync"
	"time"
)

// 自定义客户端
type Client struct {
	// 引入锁
	sync.Mutex
	// websocket连接
	conn *websocket.Conn
	// 读到的信息在这里
	readChan chan []byte
	// 需要发送的信息在这里
	writeChan chan []byte
	// 协程之间通信，连接是否已关闭
	closeChan chan struct{}
	// close closeChan 使用，避免panic
	closed bool
}

// 读取一条信息
// false 代表连接关闭失效/等待超时
func (c *Client) ReadMessage() ([]byte, bool) {
	select {
	// 超时，关闭连接
	// todo 心跳信息没有处理，在外部处理(自定义)
	case <-time.After(time.Minute*5 + time.Second*10):
		c.Close()
		return nil, false
	case data := <-c.readChan:
		return data, true
	case <-c.closeChan:
		return nil, false
	}
}

// 写入一条信息
// 不能判断成功/失败
func (c *Client) WriteMessage(data []byte) {
	select {
	case c.writeChan <- data:
	case <-c.closeChan:
	case <-time.After(time.Second * 3):
		// 为了防止阻塞，如果处理过慢，这里可以增加服务的稳定，不会导致一直阻塞影响其他接收者接收
		c.Close()
	}
}

// 关闭连接
func (c *Client) Close() {
	c.Lock()
	if !c.closed {
		close(c.closeChan)
		c.closed = true
		// ===
		// conn.Close 可以不加锁,多次使用,线程安全
		c.conn.Close()
		// ===
	}
	c.Unlock()
}

func (c *Client) readLoop() {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			c.Close()
			return
		}
		select {
		case c.readChan <- msg:
		case <-c.closeChan:
			return
		}
	}
}

func (c *Client) writeLoop() {
	for {
		select {
		case msg := <-c.writeChan:
			// 因为缓冲区的存在所以，不会立即发送，加上异常关闭的连接也被视为存在
			// 可能会出现 broken pipe
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				c.Close()
				return
			}
		case <-c.closeChan:
			return
		}
	}
}

// 升级成websocket，开启读（生产者）、写（消费者）的协程，返回自定义*Client
func NewClient(c *gin.Context) (*Client, error) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("NewClient err:", err)
		return nil, err
	}

	client := &Client{
		conn: ws,
		// 没必要设置太大，1也没问题，一般情况消费的足够快，除非是程序在发消息
		readChan:  make(chan []byte, 5),
		writeChan: make(chan []byte, 5),
		closeChan: make(chan struct{}),
	}

	go client.readLoop()
	go client.writeLoop()

	return client, nil
}

// 消息广播时使用
// 要广播的信息，源客户端信息和要发送的信息
type Message struct {
	client *Client
	msg    []byte
}

var (
	// key:client, value:struct{}{}
	clients sync.Map
	// 广播消息队列
	messages = make(chan Message, 100)

	// 广播消息发送 worker
	sendWorkers = make(chan struct{}, 100)

	// 测试时使用，客户端连接数
	//clientNum int32

	// 升级为 websocket
	upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func Chat(c *gin.Context) {
	// 获取封装好的自定义client
	conn, err := NewClient(c)
	if err != nil {
		return
	}

	// 存储至连接队列
	clients.Store(conn, struct{}{})

	//atomic.AddInt32(&clientNum, 1)
	//log.Println("now clients:", clientNum)

	go func() {
		// 收尾的函数
		defer func() {
			clients.Delete(conn)

			//atomic.AddInt32(&clientNum, -1)
			//log.Println("now clients:", clientNum)
		}()
		for {
			i, ok := conn.ReadMessage()
			if !ok {
				return
			}
			// 处理自定义心跳信息
			if string(i) != "{\"msg\":\"beat\"}" {
				tmp := Message{client: conn, msg: i}
				messages <- tmp
			}
		}
	}()
}

func sendWorker(c *Client, msg []byte) {
	sendWorkers <- struct{}{}
	go func() {
		c.WriteMessage(msg)
		<-sendWorkers
	}()
}

// 广播函数
func BroadCast() {
	// 取消息
	for msg := range messages {
		c := msg.client
		// 遍历客户端队列
		clients.Range(func(key, value interface{}) bool {
			// 不是本身就发送
			if key != c {
				// 因为 key 的类型确定，所以可以这么用
				sendWorker(key.(*Client), msg.msg)
			}
			// 继续迭代
			return true
		})
	}
}

// https 处理, 即变为 wss://*
func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     ":8849",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		if err != nil {
			return
		}
	}
}

func main() {
	// goroutine 广播消息
	go BroadCast()

	// 发布模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// 路由设定
	r.GET("/chat", Chat)
	// TLS 中间件设置
	r.Use(TlsHandler())
	// 启动服务
	// TLS 需要对应文件
	// todo 修改这里！！
	err := r.RunTLS(":8849", "*.pem", "*.key")
	if err != nil {
		panic(err)
	}
}

/*
没有条件升级https或不会升级的
注释掉
r.Use(TlsHandler())
改
err := r.RunTLS(":8849", "*.pem", "*.key")
为
err := r.Run(":8849")

使用时前端改wss:// 为 ws://
此时只能在极少部分的网站可以使用，现在大部分是https站点
*/
