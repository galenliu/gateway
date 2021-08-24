package plugin

type Units struct {
	Temperature string `json:"temperature"`
}

type Preferences struct {
	Units    Units  `json:"units"`
	Language string `json:"language"`
}

type UserProfile struct {
	BaseDir        string
	DataDir        string
	AddonsDir      string
	ConfigDir      string
	UploadDir      string
	MediaDir       string
	LogDir         string
	GatewayVersion string
}
