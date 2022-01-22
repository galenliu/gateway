package yeelight

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type (
	//Command represents COMMAND request to Yeelight device
	Command struct {
		ID     int           `json:"id"`
		Method string        `json:"method"`
		Params []interface{} `json:"params"`
	}

	// CommandResult represents response from Yeelight device
	CommandResult struct {
		ID     int           `json:"id"`
		Result []interface{} `json:"result,omitempty"`
		Error  *Error        `json:"error,omitempty"`
	}

	// Notification represents notification response
	Notification struct {
		Method string         `json:"method"`
		Params map[string]any `json:"params"`
	}

	//Error struct represents error part of response
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

func Discover() (*Yeelight, error) {
	var err error

	ssdp, _ := net.ResolveUDPAddr("udp4", ssdpAddr)
	c, _ := net.ListenPacket("udp4", ":0")
	socket := c.(*net.UDPConn)
	socket.WriteToUDP([]byte(discoverMSG), ssdp)
	socket.SetReadDeadline(time.Now().Add(timeout))

	rsBuf := make([]byte, 1024)
	size, _, err := socket.ReadFromUDP(rsBuf)
	if err != nil {
		return nil, errors.New("no devices found")
	}
	rs := rsBuf[0:size]
	addr := parseAddr(string(rs))
	fmt.Printf("Device with ip %s found\n", addr)
	supports := strings.Split(parseSupports(string(rs)), " ")
	y := New(addr, supports)
	y.parseStatus(string(rs))
	go func() {
		_, _, err := y.Listen()
		if err != nil {
			fmt.Print(err.Error())
		}
	}()
	return y, nil
}

//parseAddr parses address from ssdp response
func parseAddr(msg string) string {
	if strings.HasSuffix(msg, crlf) {
		msg = msg + crlf
	}
	resp, err := http.ReadResponse(bufio.NewReader(strings.NewReader(msg)), nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()
	return strings.TrimPrefix(resp.Header.Get("LOCATION"), "yeelight://")
}

func parseSupports(msg string) string {
	if strings.HasSuffix(msg, crlf) {
		msg = msg + crlf
	}
	resp, err := http.ReadResponse(bufio.NewReader(strings.NewReader(msg)), nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()
	return strings.TrimPrefix(resp.Header.Get("support"), " ")
}

func (y *Yeelight) parseStatus(msg string) {
	if strings.HasSuffix(msg, crlf) {
		msg = msg + crlf
	}
	resp, err := http.ReadResponse(bufio.NewReader(strings.NewReader(msg)), nil)
	if err != nil {
		fmt.Println(err)
	}
	power := strings.TrimPrefix(resp.Header.Get("power"), " ")
	if power != "" {
		y.lastStatus.Store("power", power)
	}

	bright := strings.TrimPrefix(resp.Header.Get("bright"), " ")
	if bright != "" {
		y.lastStatus.Store("bright", bright)
	}

	colorMode := strings.TrimPrefix(resp.Header.Get("color_mode"), " ")
	if colorMode != "" {
		y.lastStatus.Store("color_mode", colorMode)
	}

	ct := strings.TrimPrefix(resp.Header.Get("ct"), " ")
	if ct != "" {
		y.lastStatus.Store("ct", ct)
	}

	rgb := strings.TrimPrefix(resp.Header.Get("rgb"), " ")
	if rgb != "" {
		y.lastStatus.Store("rgb", rgb)
	}

	hue := strings.TrimPrefix(resp.Header.Get("hue"), " ")
	if hue != "" {
		y.lastStatus.Store("hue", hue)
	}

	sat := strings.TrimPrefix(resp.Header.Get("sat"), " ")
	if sat != "" {
		y.lastStatus.Store("sat", sat)
	}

	name := strings.TrimPrefix(resp.Header.Get("name"), " ")
	if name != "" {
		y.lastStatus.Store("name", name)
	}
}

func toBool(s string) bool {
	if s == "on" {
		return true
	}
	return false
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return i
	}
	return 0
}

//closeConnection closes network connection
func closeConnection(c net.Conn) {
	if nil != c {
		c.Close()
	}
}
