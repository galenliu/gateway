package gateway

import "fmt"

const (
	MajorVersion = 0
	MinorVersion = 1
	PatchVersion = "0.dev0"

	AddonsDir = "addons"
	DataDir   = "data"
	ConfigDir = "config"
	MediaDir  = "media"
	UploadDir = "upload"
	LogDir    = "log"

	ApiPrefix  = "/api"
	ThingsPath = "/things"
	AddonsPath = "/addons"
)

var (
	ShortVersion = fmt.Sprintf("%v.%v", MajorVersion, MinorVersion)
	Version      = fmt.Sprintf("%v.%v", ShortVersion, PatchVersion)
)
