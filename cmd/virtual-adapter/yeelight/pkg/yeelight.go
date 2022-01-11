package yeelight

import (
	"bufio"
	"encoding/json"
	"fmt"
	c "github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/pkg/color"
	"github.com/galenliu/gateway/cmd/virtual-adapter/yeelight/pkg/utils"
	"image/color"
	"math/rand"
	"net"
	"sync"
	"time"
)

const (
	discoverMSG = "M-SEARCH * HTTP/1.1\r\n HOST:239.255.255.250:1982\r\n MAN:\"ssdp:discover\"\r\n ST:wifi_bulb\r\n"

	// timeout value for TCP and UDP commands
	timeout = time.Second * 3
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
	select {
	case y.sendChan <- cmd:
		return nil, nil
	}
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
	return on.(bool), nil
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

	var err error
	notifCh := make(chan *Notification)
	done := make(chan struct{}, 1)

	conn, err := net.DialTimeout("tcp", y.addr, time.Second*3)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot connect to %s. %s", y.addr, err)
	}

	fmt.Println("Connection established")
	go func(c net.Conn) {
		//make sure connection is closed when method returns
		defer closeConnection(conn)
		connReader := bufio.NewReader(c)
		for {
			select {
			case c := <-y.sendChan:
				fmt.Printf("send Channel rev: %v \t\n", c)
				b, _ := json.Marshal(c)
				fmt.Fprint(conn, string(b)+crlf)
			case <-done:
				return
			default:
				data, err := connReader.ReadString('\n')
				if nil == err {
					var rs Notification
					fmt.Println(data)
					json.Unmarshal([]byte(data), &rs)
					select {
					case notifCh <- &rs:
					default:
						fmt.Println("Channel is full")
					}
				}
			}

		}

	}(conn)

	return notifCh, done, nil
}
