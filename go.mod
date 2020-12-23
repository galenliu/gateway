module gateway

go 1.15

require (
	gitee.com/liu_guilin/WebThings-schema v1.0.0
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/gobuffalo/packr/v2 v2.8.1
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/json-iterator/go v1.1.10
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.3 // indirect
	github.com/mattn/go-sqlite3 v1.14.3
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/xiam/to v0.0.0-20200126224905-d60d31e03561
	go.uber.org/zap v1.16.0
	golang.org/x/sys v0.0.0-20200929083018-4d22bbb62b3c // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
	gorm.io/driver/sqlite v1.1.3
	gorm.io/gorm v1.20.7
)

replace gitee.com/liu_guilin/WebThings-schema v1.0.0 => ../gateway-schema
