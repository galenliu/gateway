package addon

import (
	"fmt"
	"github.com/gorilla/websocket"
	json "github.com/json-iterator/go"
	"net/url"
)

const (
	Disconnect = "Disconnect"
	Connected  = "Connected"
	Registered = "Registered"
)

type UserProfile struct {
	BaseDir        string `validate:"required" json:"base_dir"`
	DataDir        string `validate:"required" json:"data_dir"`
	AddonsDir      string `validate:"required" json:"addons_dir"`
	ConfigDir      string `validate:"required" json:"config_dir"`
	UploadDir      string `validate:"required" json:"upload_dir"`
	MediaDir       string `validate:"required" json:"media_dir"`
	LogDir         string `validate:"required" json:"log_dir"`
	GatewayVersion string
}

type Preferences struct {
	Language string `validate:"required" json:"language"`
	Units    Units  `validate:"required" json:"units"`
}

type Units struct {
	Temperature string `validate:"required" json:"temperature"`
}

type OnMessage func(data []byte)

type IpcClient struct {
	ws          *websocket.Conn
	manager     *Manager
	url         string
	preferences Preferences
	userProfile UserProfile

	writeCh   chan []byte
	readCh    chan []byte
	closeChan chan any
	reConnect chan any

	gatewayVersion string

	onMessage OnMessage

	status   string
	pluginId string
	origin   string
	verbose  bool
}

// NewClient 新建一个Client，注册消息Handler
func NewClient(PluginId string, manager *Manager) *IpcClient {
	u := url.URL{Scheme: "ws", Host: "localhost:" + IpcDefaultPort, Path: "/"}
	client := &IpcClient{}
	client.pluginId = PluginId
	client.url = u.String()
	client.status = Disconnect

	client.closeChan = make(chan any)
	client.reConnect = make(chan any)

	client.readCh = make(chan []byte)
	client.writeCh = make(chan []byte)

	client.onMessage = manager.onMessage
	client.Register()
	go client.readLoop()
	return client
}

func (client *IpcClient) onData(data []byte) {

	fmt.Printf("read message : %s \t\n", string(data))

	if json.Get(data, "messageType").ToInt() == PluginRegisterResponse {
		client.preferences.Language = json.Get(data, "data", "preferences", "language").ToString()
		client.preferences.Units.Temperature = json.Get(data, "data", "preferences", "units", "temperature").ToString()
		client.userProfile.AddonsDir = json.Get(data, "data", "user_profile", "addons_dir").ToString()
		client.userProfile.BaseDir = json.Get(data, "data", "user_profile", "base_dir").ToString()
		client.userProfile.ConfigDir = json.Get(data, "data", "user_profile", "config_dir").ToString()
		client.userProfile.DataDir = json.Get(data, "data", "user_profile", "data_dir").ToString()
		client.userProfile.GatewayVersion = json.Get(data, "data", "user_profile", "gateway_version").ToString()
		client.userProfile.LogDir = json.Get(data, "data", "user_profile", "log_dir").ToString()
		client.userProfile.MediaDir = json.Get(data, "data", "user_profile", "media_dir").ToString()
		client.userProfile.UploadDir = json.Get(data, "data", "user_profile", "upload_dir").ToString()
		client.status = Registered
	} else {
		client.onMessage(data)
	}
}

func (client *IpcClient) sendMessage(data []byte) {

	if client.ws != nil && client.status == Registered {
		err := client.ws.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			fmt.Printf("ipc client write err")
			client.status = Disconnect
		}
		fmt.Printf("ipc client send message: %s \t\n", string(data))
	}
}

func (client *IpcClient) readMessage() {
	if client.ws != nil {
		_, data, err := client.ws.ReadMessage()
		if err != nil {
			fmt.Printf("read faild, websocket err", err.Error())
			client.status = Disconnect
		}
		client.onData(data)
	}
}

func (client *IpcClient) readLoop() {
	for {
		if client.status == Registered && client.ws != nil {
			client.readMessage()
		}
	}
}

func (client *IpcClient) dial() error {

	var err error = nil
	client.ws, _, err = websocket.DefaultDialer.Dial(client.url, nil)
	if err != nil {
		fmt.Printf("dial err: %s \r\n", err.Error())
		return err
	}
	return nil
}

func (client *IpcClient) Register() {

	err := client.dial()
	if err != nil {
		return
	}

	message := struct {
		MessageType int `json:"messageType"`
		Data        any `json:"data"`
	}{
		MessageType: PluginRegisterRequest,
		Data: struct {
			PluginID string `json:"pluginId"`
		}{PluginID: client.pluginId},
	}

	d, _ := json.MarshalIndent(message, "", " ")
	_ = client.ws.WriteMessage(websocket.BinaryMessage, d)
	_, data, err := client.ws.ReadMessage()
	if err != nil {
		fmt.Printf("read faild, websocket err", err.Error())
		client.status = Disconnect
	}
	client.onData(data)

}

func (client *IpcClient) close() {
	if client.ws != nil {
		_ = client.ws.Close()
	}
}
