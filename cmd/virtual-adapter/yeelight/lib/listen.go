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
			defer conn.Close()
			if err != nil {
				break
			}
			for {
				bytes, err := bufio.NewReader(conn).ReadBytes('\n')
				if err != nil {
					fmt.Printf("Error reading %s", err.Error())
					break
				}
				select {
				case <-ctx.Done():
					fmt.Printf("listen exit")
					return
				default:
					select {
					case messageChan <- string(bytes):
					default:
						fmt.Printf("read chan if full")
					}
				}
			}
			time.Sleep(time.Second * 10)
		}
	}()
	return messageChan, nil
}
