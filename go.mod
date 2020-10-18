module gateway

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Joker/hpp v1.0.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/galeuliu/gateway-schema v1.0.0
	github.com/gin-gonic/gin v1.6.3
	github.com/gorilla/websocket v1.4.2
	github.com/grandcat/zeroconf v1.0.0
	github.com/json-iterator/go v1.1.10
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/kataras/iris/v12 v12.2.0-alpha.0.20200930073625-5017e3c986be // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-sqlite3 v1.14.3
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/onsi/ginkgo v1.14.1 // indirect
	github.com/onsi/gomega v1.10.2 // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/xiam/to v0.0.0-20200126224905-d60d31e03561
	github.com/yudai/pp v2.0.1+incompatible // indirect
	go.uber.org/zap v1.16.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)

replace github.com/galeuliu/gateway-schema v1.0.0 => ../gateway-schema
