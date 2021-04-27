module github.com/galenliu/gateway

go 1.16

require (
	github.com/asaskevich/EventBus v0.0.0-20200907212545-49d423059eef
	github.com/brutella/hc v1.2.4
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fasthttp/websocket v1.4.3 // indirect
	github.com/galenliu/gateway-addon v1.0.0
	github.com/gin-gonic/gin v1.7.1
	github.com/go-oauth2/oauth2/v4 v4.3.0
	github.com/gofiber/fiber/v2 v2.7.1
	github.com/gofiber/websocket/v2 v2.0.3
	github.com/gorilla/websocket v1.4.2
	github.com/json-iterator/go v1.1.10
	github.com/klauspost/compress v1.12.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.5
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/savsgio/gotils v0.0.0-20210316171653-c54912823645 // indirect
	github.com/tidwall/gjson v1.7.4
	github.com/valyala/fasthttp v1.23.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/xiam/to v0.0.0-20200126224905-d60d31e03561
	go.uber.org/zap v1.16.0
	golang.org/x/sys v0.0.0-20210415045647-66c3f260301c // indirect
	golang.org/x/tools v0.0.0-20191127201027-ecd32218bd7f // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/galenliu/gateway-addon v1.0.0 => ./gateway-addon
