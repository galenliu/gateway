package plugin

import (
	"errors"
	"fmt"
	"github.com/fasthttp/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 8172
)

// ws 的所有连接
// 用于广播
var wsConnAll map[int64]*wsConnection
var maxConnId int64

// 客户端读写消息
type wsMessage struct {
	// websocket.TextMessage 消息类型
	messageType int
	data        []byte
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 客户端连接
type wsConnection struct {
	registered bool
	wsSocket   *websocket.Conn // 底层websocket
	inChan     chan *wsMessage // 读队列
	outChan    chan *wsMessage // 写队列
	mutex      sync.Mutex      // 避免重复关闭管道,加锁处理
	isClosed   bool
	closeChan  chan byte // 关闭通知
	id         int64
}

var clientChan chan *wsConnection

func NewIpcServer(addr string) chan *wsConnection {

	clientChan = make(chan *wsConnection, 64)
	http.HandleFunc("/", serveWs)
	srv := http.Server{
		Addr:    addr,
		Handler: http.DefaultServeMux,
	}

	//
	//go func() {
	//	<-ctx.Done()
	//	wg.Add(1)
	//	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	//	defer cancel()
	//	err := srv.Shutdown(ctx)
	//	if err != nil {
	//		log.Printf("error shutting down:%s", err.Error())
	//	}
	//	wg.Done()
	//	close(errChan)
	//	close(clientChan)
	//}()

	go func() {
		log.Println("listening at " + addr)
		err := srv.ListenAndServe()
		fmt.Println("waiting for the remaining connections to finish...")
		if err != nil && err != http.ErrServerClosed {
			close(clientChan)
		}
		log.Println("gracefully shutdown the http server...")
	}()
	return clientChan
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	maxConnId++
	c := &wsConnection{
		registered: false,
		wsSocket:   ws,
		inChan:     make(chan *wsMessage, 1000),
		outChan:    make(chan *wsMessage, 1000),
		closeChan:  make(chan byte),
		isClosed:   false,
		id:         maxConnId,
	}
	// 读协程
	go c.wsReadLoop()
	// 写协程
	go c.wsWriteLoop()
	select {
	case clientChan <- c:
	}
}

// 处理消息队列中的消息
func (wsConn *wsConnection) wsReadLoop() {
	// 设置消息的最大长度
	wsConn.wsSocket.SetReadLimit(maxMessageSize)
	wsConn.wsSocket.SetReadDeadline(time.Now().Add(pongWait))
	wsConn.wsSocket.SetPongHandler(func(string) error {
		wsConn.wsSocket.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		// 读一个message
		msgType, data, err := wsConn.wsSocket.ReadMessage()
		if err != nil {
			log.Println(websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure))
			log.Println("消息读取出现错误", err.Error())
			wsConn.close()
			return
		}
		req := &wsMessage{
			msgType,
			data,
		}
		// 放入请求队列,消息入栈
		select {
		case wsConn.inChan <- req:
		case <-wsConn.closeChan:
			return
		}
	}
}

// 发送消息给客户端
func (wsConn *wsConnection) wsWriteLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		// 取一个应答
		case msg := <-wsConn.outChan:
			// 写给websocket
			if err := wsConn.wsSocket.WriteMessage(msg.messageType, msg.data); err != nil {
				log.Println("发送消息给客户端发生错误", err.Error())
				// 切断服务
				wsConn.close()
				return
			}
		case <-wsConn.closeChan:
			// 获取到关闭通知
			return
		case <-ticker.C:
			// 出现超时情况
			if err := wsConn.wsSocket.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
			}
			//if err := wsConn.wsSocket.WriteMessage(websocket.PingMessage, nil); err != nil {
			//	return
			//}
		}
	}
}

// 写入消息到队列中
func (wsConn *wsConnection) wsWrite(messageType int, data []byte) error {
	select {
	case wsConn.outChan <- &wsMessage{messageType, data}:
	case <-wsConn.closeChan:
		return errors.New("连接已经关闭")
	}
	return nil
}

// 读取消息队列中的消息
func (wsConn *wsConnection) wsRead() (*wsMessage, error) {
	select {
	case msg := <-wsConn.inChan:
		// 获取到消息队列中的消息
		return msg, nil
	case <-wsConn.closeChan:

	}
	return nil, errors.New("连接已经关闭")
}

// 关闭连接
func (wsConn *wsConnection) close() {
	log.Println("关闭连接被调用了")
	wsConn.wsSocket.Close()
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if wsConn.isClosed == false {
		wsConn.isClosed = true
		// 删除这个连接的变量
		delete(wsConnAll, wsConn.id)
		close(wsConn.closeChan)
	}
}
