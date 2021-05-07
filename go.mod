module github.com/galenliu/gateway

go 1.16

require (
	github.com/arsmn/fiber-swagger/v2 v2.6.0
	github.com/asaskevich/EventBus v0.0.0-20200907212545-49d423059eef
	github.com/brutella/hc v1.2.4
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fasthttp/websocket v1.4.3 // indirect
	github.com/galenliu/gateway-addon v1.0.0
	github.com/gin-gonic/gin v1.7.1
	github.com/go-oauth2/oauth2/v4 v4.3.0
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/gofiber/fiber/v2 v2.9.0
	github.com/gofiber/websocket/v2 v2.0.3
	github.com/gorilla/websocket v1.4.2
	github.com/json-iterator/go v1.1.10
	github.com/klauspost/compress v1.12.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.5
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/pkg/errors v0.9.1 // indirect
	github.com/savsgio/gotils v0.0.0-20210316171653-c54912823645 // indirect
	github.com/tidwall/gjson v1.7.4
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20210501142056-aec3718b3fa0 // indirect
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace github.com/galenliu/gateway-addon v1.0.0 => ./gateway-addon
