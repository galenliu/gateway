package yeelight

import (
	"bufio"
	"encoding/json"
	"fmt"
	c "github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/pkg/color"
	"github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/pkg/utils"
	"image/color"
	"io"
	"math/rand"
	"net"
	"sync"
	"time"
)

const (
	discoverMSG = "M-SEARCH * HTTP/1.1\r\n HOST:239.255.255.250:1982\r\n MAN:\"ssdp:discover\"\r\n ST:wifi_bulb\r\n"

	// timeout value for TCP and UDP commands
	timeout = time.Second * 2
	//SSDP discover address
	ssdpAddr = "239.255.255.250:1982"
	//CR-LF delimiter
	crlf = "\r\n"
)

type Mode int

type Yeelight struct {
	addr       string
	supports   []string
	lastStatus sync.Map
	rnd        *rand.Rand
	sendChan   chan *Command
}

func New(addr string, supports []string) *Yeelight {
	y := &Yeelight{
		addr:     addr,
		supports: supports,
		rnd:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return y
}

func (y *Yeelight) executeCommand(name string, params ...interface{}) (*CommandResult, error) {
	return y.execute(y.newCommand(name, params))
}

func (y *Yeelight) execute(cmd *Command) (*CommandResult, error) {

	if y.sendChan != nil {
		select {
		case y.sendChan <- cmd:
		}
		return nil, nil
	}
	conn, err := net.Dial("tcp", y.addr)
	if nil != err {
		conn = nil
		return nil, fmt.Errorf("cannot open connection to %s. %s", y.addr, err)
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(timeout))

	//write request/command
	b, _ := json.Marshal(cmd)
	_, err = io.WriteString(conn, string(b)+crlf)
	if err != nil {
		fmt.Println(err.Error() + "\t\n")
	}
	//wait and read for response
	res, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("cannot read command result %s", err)
	}
	var rs CommandResult
	err = json.Unmarshal([]byte(res), &rs)
	if nil != err {
		return nil, fmt.Errorf("cannot parse command result %s", err)
	}
	if nil != rs.Error {
		return nil, fmt.Errorf("command execution error. Code: %d, Message: %s", rs.Error.Code, rs.Error.Message)
	}
	return &rs, nil
}

func (y *Yeelight) newCommand(name string, params []interface{}) *Command {
	return &Command{
		Method: name,
		ID:     y.randID(),
		Params: params,
	}
}

func (y *Yeelight) randID() int {
	i := y.rnd.Intn(100)
	return i
}

func (y *Yeelight) TurnOn() (*CommandResult, error) {
	return y.executeCommand("set_power", "on")
}

func (y *Yeelight) TurnOnWithParams(mode Mode, duration int) (*CommandResult, error) {
	return y.executeCommand("set_power", "on", duration, mode)
}

func (y *Yeelight) TurnOff() (*CommandResult, error) {
	return y.executeCommand("set_power", "off")
}

func (y *Yeelight) IsOn() (bool, error) {
	on, _ := y.lastStatus.Load("power")
	return on.(string) == "on", nil
}

func (y *Yeelight) SetPower(b bool) {
	if b {
		_, _ = y.TurnOn()
	} else {
		_, _ = y.TurnOff()
	}
}

func (y *Yeelight) SetBrightness(brightness int) (*CommandResult, error) {
	on, err := y.IsOn()
	if err == nil {
		if !on {
			_, err := y.TurnOn()
			if err != nil {
				return nil, err
			}
		}
	}
	return y.executeCommand("set_bright", utils.GetBrightnessValue(brightness))
}

func (y *Yeelight) SetRGB(rgba color.RGBA) (*CommandResult, error) {
	on, err := y.IsOn()
	if err == nil {
		if !on {
			_, err := y.TurnOn()
			if err != nil {
				return nil, err
			}
		}
	}
	value := c.RGBToYeelight(rgba)
	return y.executeCommand("set_rgb", value)
}

func (y *Yeelight) SetHSV(hue int, saturation int) (*CommandResult, error) {
	on, err := y.IsOn()
	if err == nil {
		if !on {
			_, err := y.TurnOn()
			if err != nil {
				return nil, err
			}
		}
	}
	return y.executeCommand("set_rgb", hue, saturation)
}

func (y *Yeelight) SetBrightnessWithDuration(brightness int, duration int) (*CommandResult, error) {
	on, err := y.IsOn()
	if err == nil {
		if !on {
			_, err := y.TurnOn()
			if err != nil {
				return nil, err
			}
		}
	}
	return y.executeCommand("set_bright", brightness, duration)
}

func (y *Yeelight) StartFlow(flow *Flow) (*CommandResult, error) {

	params := flow.AsStartParams()
	return y.executeCommand("start_cf", params)
}

func (y *Yeelight) StopFlow() (*CommandResult, error) {
	return y.executeCommand("stop_cf", "")
}

// GetProp method is used to retrieve current property of smart LED.
func (y *Yeelight) GetProp(values ...interface{}) ([]interface{}, error) {
	r, err := y.executeCommand("get_prop", values...)
	if nil != err {
		return nil, err
	}
	return r.Result, nil
}

func (y *Yeelight) GetAddr() string {
	return y.addr
}

func (y *Yeelight) GetSupports() []string {
	return y.supports
}

func (y *Yeelight) SetName(name string) (*CommandResult, error) {
	return y.executeCommand("set_name", name)
}

func (y *Yeelight) Listen() (<-chan *Notification, chan<- struct{}, error) {

	y.sendChan = make(chan *Command, 50)
	notifCh := make(chan *Notification, 50)
	done := make(chan struct{}, 1)

	revTask := func(conn net.Conn) {
		for {
			msg, err := bufio.NewReader(conn).ReadString('\n')
			if err == nil {
				var rs Notification
				err = json.Unmarshal([]byte(msg), &rs)
				if err != nil {
					fmt.Println(err.Error())
					break
				}
				select {
				case <-done:
					return
				case notifCh <- &rs:
				default:
					fmt.Println("Channel is full")
				}
			} else {
				fmt.Printf("rev data err: %v\n", err.Error())
				break
			}
			time.Sleep(1 * time.Second)
		}
	}

	sendTask := func(conn net.Conn, done <-chan struct{}) {
		for {
			select {
			case <-done:
				return
			case cmd := <-y.sendChan:
				b, _ := json.Marshal(cmd)
				_, err := io.WriteString(conn, string(b)+crlf)
				if err != nil {
					fmt.Println(err.Error() + "\t\n")
				}
				//wait and read for response
				res, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					fmt.Println(err.Error())
				}
				var rs CommandResult
				err = json.Unmarshal([]byte(res), &rs)
				if nil != err {
					fmt.Println(err.Error())
				}
				if nil != rs.Error {
					fmt.Printf("command execution error. Code: %d, Message: %s \t\n", rs.Error.Code, rs.Error.Message)
				}
			}

		}
	}

	connection := func() {
		for {
			conn, err := net.Dial("tcp", y.addr)
			done := make(chan struct{})
			if err != nil {
				fmt.Println(err.Error())
			} else {
				defer conn.Close()
				go sendTask(conn, done)
				revTask(conn)
			}
			done <- struct{}{}
			time.Sleep(3 * time.Second)
		}
	}

	go connection()

	return notifCh, done, nil
}

func (y *Yeelight) GetPropertyValue(name string) any {
	v, ok := y.lastStatus.Load(name)
	if ok {
		return v
	}
	return nil
}
