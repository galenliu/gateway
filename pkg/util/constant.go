package util

import "fmt"

const (
	MajorVersion = 0
	MinorVersion = 1
	PatchVersion = 0

	AddonsDir = "addons"
	DataDir   = "data"
	ConfigDir = "config"
	MediaDir  = "media"
	UploadDir = "upload"
	LogDir    = "logger"



	DbPrefLang = "preferences.language"
	PrefLangCn = "zh-CN"

	DbPrefUnitsTemp      = "preferences.units.temperature"
	PrefUnitsTempCelsius = "degree celsius"
)

var (
	ShortVersion = fmt.Sprintf("%v.%v", MajorVersion, MinorVersion)
	Version      = fmt.Sprintf("%v.%v", ShortVersion, PatchVersion)
)

const (
	ConfDirName    = ".profile"
)
