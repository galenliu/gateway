module gateway

go 1.15

require (
	github.com/galeuliu/gateway-schema v1.0.0
	github.com/gin-gonic/gin v1.6.3 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/json-iterator/go v1.1.10
	go.uber.org/zap v1.16.0 // indirect
)

replace github.com/galeuliu/gateway-schema v1.0.0 => ../gateway-schema
