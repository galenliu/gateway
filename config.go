package gateway



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



