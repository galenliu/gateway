module gateway

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/galeuliu/gateway-schema v1.0.0
	github.com/gin-gonic/gin v1.6.3
	github.com/gorilla/websocket v1.4.2
	github.com/json-iterator/go v1.1.10
	github.com/kataras/iris/v12 v12.2.0-alpha.0.20200930073625-5017e3c986be // indirect
	github.com/mattn/go-sqlite3 v1.14.3
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/xiam/to v0.0.0-20200126224905-d60d31e03561
	go.uber.org/zap v1.16.0
)

replace github.com/galeuliu/gateway-schema v1.0.0 => ../gateway-schema
