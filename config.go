package gateway

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"github.com/xiam/to"
	"os"
	"path"
)

const (
	YamlConfigFile = "configuration.yaml"
	ConfDirName    = ".smartassistant"
)

type Config struct {
	Name        string
	Service     string
	Domain      string
	Host        string
	TextRecords map[string]string
	Port        int
	Pin         string

	SetupId string

	name     string //Accessory name
	id       string //Accessory id
	protocol string

	discoverable bool // Flag if accessory is discoverable (sf)
	mfiCompliant bool // Flag if accessory if Mfi compliant (ff)

	categoryId int //category type
	servePort  int
	version    int
	state      int

	configHash []byte
}

func GetDefaultConfigDir() string {
	dir, _ := os.UserHomeDir()
	dirPath := path.Join(dir, ConfDirName)
	return dirPath
}

func (cfg *Config) setupHash() string {
	hashvalue := fmt.Sprintf("%s%s", cfg.SetupId, cfg.id)
	sum := sha512.Sum512([]byte(hashvalue))
	// use only first 4 bytes
	code := []byte{sum[0], sum[1], sum[2], sum[3]}
	encoded := base64.StdEncoding.EncodeToString(code)
	return encoded
}

func (cfg *Config) txtRecords() map[string]string {
	return map[string]string{
		"pv": cfg.protocol,
		"id": cfg.id,
		"c#": fmt.Sprintf("%d", cfg.version),
		"s#": fmt.Sprintf("%d", cfg.state),
		"sf": fmt.Sprintf("%d", to.Int64(cfg.discoverable)),
		"ff": fmt.Sprintf("%d", to.Int64(cfg.mfiCompliant)),
		"md": cfg.name,
		"ci": fmt.Sprintf("%d", cfg.categoryId),
		"sh": cfg.setupHash(),
	}
}
