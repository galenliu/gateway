package lib

import (
	"bufio"
	"errors"
	"fmt"
	json "github.com/json-iterator/go"
	"math/rand"
	"net"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	ssdpSearch = "M-SEARCH * HTTP/1.1\r\n MAN:\"ssdp:discover\"\r\n ST:wifi_bulb\r\n"

	// timeout value for TCP and UDP commands
	timeout = time.Second * 3

	//SSDP discover address
	ssdpAddr = "239.255.255.250:1982"
)

// Light properties
type Light struct {
	Location   string   `json:"location,omitempty"`
	ID         string   `json:"id,omitempty"`
	Model      string   `json:"model,omitempty"`
	FWVersion  int      `json:"fw_ver,omitempty"`
	Support    []string `json:"support,omitempty"`
	Power      string   `json:"power,omitempty"`
	Bright     int      `json:"bright,omitempty"`
	ColorMode  int      `json:"color_mode,omitempty"`
	ColorTemp  int      `json:"ct,omitempty"`
	RGB        int      `json:"rgb,omitempty"`
	Hue        int      `json:"hue,omitempty"`
	Saturation int      `json:"sat,omitempty"`
	Name       string   `json:"name,omitempty"`
}

// Command to send to the light
type Command struct {
	ID     int32       `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

// Result reply of the sent command
type Result struct {
	ID     int32       `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

func Discover() []Light {

	lAddr, _ := net.ResolveUDPAddr("udp4", ":0")

	mAddr, _ := net.ResolveUDPAddr("udp4", ssdpAddr)

	conn, _ := net.ListenUDP("udp4", lAddr)

	defer conn.Close()

	go func() {
		conn.WriteToUDP([]byte(ssdpSearch), mAddr)
	}()
	_ = conn.SetDeadline(time.Now().Add(timeout))

	answers := make(map[string]string)
	for {
		answer := make([]byte, 1024)
		n, src, err := conn.ReadFromUDP(answer)
		if err != nil {
			break
		}
		answers[src.String()] = string(answer[:n])
	}

	var lights []Light
	for _, answer := range answers {
		tp := textproto.NewReader(bufio.NewReader(strings.NewReader(answer)))
		tp.ReadLine()
		header, _ := tp.ReadMIMEHeader()

		var light Light
		location, _ := url.Parse(header.Get("location"))
		light.Location = location.Host
		light.ID = header.Get("id")
		light.Model = header.Get("model")
		light.FWVersion, _ = strconv.Atoi(header.Get("fw_ver"))
		light.Support = strings.Split(header.Get("support"), " ")
		light.Power = header.Get("power")
		light.Bright, _ = strconv.Atoi(header.Get("bright"))
		light.ColorMode, _ = strconv.Atoi(header.Get("color_mode"))
		light.ColorTemp, _ = strconv.Atoi(header.Get("ct"))
		light.RGB, _ = strconv.Atoi(header.Get("rgb"))
		light.Hue, _ = strconv.Atoi(header.Get("hue"))
		light.Saturation, _ = strconv.Atoi(header.Get("sat"))
		light.Name = header.Get("name")
		lights = append(lights, light)
	}
	return lights
}

func (m *Light) PowerOn(duration int) (string, error) {
	cmd := Command{
		Method: "set_power",
		Params: []interface{}{"on", "smooth", duration},
	}
	return m.SendCommand(cmd)
}

func (m *Light) PowerOff(duration int) (string, error) {
	cmd := Command{
		Method: "set_power",
		Params: []interface{}{"off", "smooth", duration},
	}
	return m.SendCommand(cmd)
}

// Toggle toggles light state on or off
func (m *Light) Toggle() (string, error) {
	cmd := Command{
		Method: "toggle",
		Params: []interface{}{},
	}
	return m.SendCommand(cmd)
}

// SetBrightness takes the brightness (0-100), transition duration and set brightness of the light
func (m *Light) SetBrightness(level, duration int) (string, error) {
	cmd := Command{
		Method: "set_bright",
		Params: []interface{}{level, "smooth", duration},
	}
	return m.SendCommand(cmd)
}

// SetRGB takes r, g, b values (0-255), transition duration and set color of the light
func (m *Light) SetRGB(red, green, blue, duration int) (string, error) {
	rgb := (red << 16) + (green << 8) + blue

	cmd := Command{
		Method: "set_rgb",
		Params: []interface{}{rgb, "smooth", duration},
	}
	return m.SendCommand(cmd)
}

// SetTemp takes temperature values (1700-6500K), transition duration and set color temperature of the light
func (m *Light) SetTemp(temp, duration int) (string, error) {
	cmd := Command{
		Method: "set_ct_abx",
		Params: []interface{}{temp, "smooth", duration},
	}
	return m.SendCommand(cmd)
}

// Update the current status of the light
func (m *Light) Update() error {
	cmd := Command{
		Method: "get_prop",
		Params: []interface{}{"power", "bright", "ct", "rgb", "hue", "sat", "color_mode"},
	}

	res, err := m.SendCommand(cmd)
	if err != nil {
		return err
	}
	var ures Result
	if err := json.Unmarshal([]byte(res), &ures); err != nil {
		return err
	}
	m.Power = ures.Result.([]interface{})[0].(string)
	m.Bright, _ = strconv.Atoi(ures.Result.([]interface{})[1].(string))
	m.ColorTemp, _ = strconv.Atoi(ures.Result.([]interface{})[2].(string))
	m.RGB, _ = strconv.Atoi(ures.Result.([]interface{})[3].(string))
	m.Hue, _ = strconv.Atoi(ures.Result.([]interface{})[4].(string))
	m.Saturation, _ = strconv.Atoi(ures.Result.([]interface{})[5].(string))
	m.ColorMode, _ = strconv.Atoi(ures.Result.([]interface{})[6].(string))

	return nil
}

func (m *Light) SendCommand(cmd Command) (string, error) {
	if cmd.ID == 0 {
		r := rand.NewSource(time.Now().UnixNano())
		cmd.ID = rand.New(r).Int31()
	}

	conn, err := net.Dial("tcp", m.Location)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	cmdJSON, err := json.MarshalIndent(cmd, "", "")
	if err != nil {
		return "", err
	}
	if _, err := fmt.Fprintf(conn, "%s\r\n", cmdJSON); err != nil {
		return "", err
	}

	return bufio.NewReader(conn).ReadString('\n')

}

type RGB struct {
	R, G, B int
}

func HTMLToRGB(in string) (RGB, error) {
	if in[0] == '#' {
		in = in[1:]
	}

	if len(in) != 6 {
		return RGB{}, errors.New("Invalid string length")
	}

	var r, g, b byte
	if n, err := fmt.Sscanf(in, "%2x%2x%2x", &r, &g, &b); err != nil || n != 3 {
		return RGB{}, err
	}

	return RGB{int(r) / 255, int(g) / 255, int(b) / 255}, nil
}
