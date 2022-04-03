package yeelight

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"time"
)

func (c Client) Listen(ctx context.Context) (<-chan string, error) {
	messageChan := make(chan string, 10)

	go func() {
		for {
			fmt.Printf("开始连接:%s \t\n", c.host)
			conn, err := net.Dial("tcp", c.host)
			if err != nil {
				fmt.Printf("连接失败:%s \t\n", c.host)
				break
			}
			for {
				err = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
				bytes, err := bufio.NewReader(conn).ReadBytes('\n')
				e, ok := err.(net.Error)
				if ok && e.Timeout() {
					fmt.Printf("读取超时:%s \t\n", c.host)
					continue
				}
				if err != nil {
					fmt.Printf("读取错误 Ip：%s err: %s \t\n", c.host, err.Error())
					//关闭这个连接
					err := conn.Close()
					if err != nil {
						fmt.Printf("关闭错误： %s \t\n", err.Error())
					}
					break
				}
				select {
				case <-ctx.Done():
					fmt.Printf("关闭连接退出监听 \t\n")
					conn.Close()
					return
				default:
					select {
					case messageChan <- string(bytes):
					default:
						fmt.Printf("消息通道已满\t\n")
					}
				}
			}
			//2秒后重新连接
			fmt.Printf("2秒后重新连接 \t\n")
			time.Sleep(time.Second * 2)
		}
	}()
	return messageChan, nil
}
