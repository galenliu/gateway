package characteristic

const (
	PermRead                    = "pr" // can be read
	PermWrite                   = "pw" // can be written
	PermEvents                  = "ev" // sends events
	PermHidden                  = "hd" // is hidden
	PermAdditionalAuthorization = "aa" //额外授权
)

const (
	FormatString string = "string"
	FormatBool   string = "bool"
	FormatFloat  string = "float"
	FormatInt8   string = "int8"
	FormatInt16  string = "int16"
	FormatInt    string = "int"

	FormatUint8  string = "uint8"
	FormatUint16 string = "uint16"
	FormatUint32 string = "uint32"
	FormatUint64 string = "uint64"
	FormatTLV8   string = "tlv8"
	FormatData   string = "data"
)

// HAP characteristic units
const (
	UnitPercentage = "percentage" //百分比
	UnitArcDegrees = "arcdegrees" //度数
	UnitCelsius    = "celsius"
	UnitLux        = "lux"
	UnitSeconds    = "seconds"
	UnitPPM        = "ppm"
)
