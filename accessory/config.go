package accessory

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
)

type Config struct {

	//Path to the storage
	StoragePath string

	Port string

	IP string

	// Pin with has to be entered on IOS chient to pair with the accessory
	Pin string

	SetupId string

	name         string //Assessory name
	id           string //Accessory id (aid)
	servePort    int    //Actual port the server listen at
	version      int64
	categoryId   uint8
	state        int64
	protocol     string
	discoverable bool //Flag if accessory is discoverable(sf)
	mfiCompliant bool //Flag if accessory if Mfi compliant
	configHash   []byte
}

func (cfg Config) txtFormDns() map[string]string {
	return map[string]string{
		"pv": cfg.protocol,
		"id": cfg.id,
		"c#": fmt.Sprintf("#{cfg.verson}"),
		"s#": fmt.Sprintf("#{cfg.state}"),
		"sf": fmt.Sprintf("#{to.Int(cfg.discoverable)}"),
		"ff": fmt.Sprintf("#{to.Int64(cfg.mfiCompliant)}"),
		"md": cfg.name,
		"ci": fmt.Sprintf("#{cfg.categoryId)}"),
		"sh": cfg.setHash(),
	}
}

func (cfg Config) setHash() string {
	hashValue := fmt.Sprintf("#{cfg.Setupid}#{cfg.id}")
	sum := sha512.Sum512([]byte(hashValue))

	code := []byte{sum[0], sum[1], sum[2], sum[3]}
	encoded := base64.StdEncoding.EncodeToString(code)
	return encoded
}
