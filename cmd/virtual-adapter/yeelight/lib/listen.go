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
			conn, err := net.Dial("tcp", c.host)
			err = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			if err != nil {
				break
			}
			defer conn.Close()
			if err != nil {
				break
			}
			for {
				bytes, err := bufio.NewReader(conn).ReadBytes('\n')
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
					return
				default:
					select {
					case messageChan <- string(bytes):
					default:
						fmt.Printf("消息通道已满\t\n")
					}
				}
			}
			//10秒后重新连接
			time.Sleep(time.Second * 2)
			fmt.Printf("10秒后重新连接 \t\n")
		}
	}()
	return messageChan, nil
}
