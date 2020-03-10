module isp-gate-service

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fasthttp/websocket v1.4.2
	github.com/go-pg/pg/v9 v9.1.3
	github.com/go-redis/redis v6.15.7+incompatible
	github.com/golang/protobuf v1.3.4
	github.com/integration-system/go-cmp v0.0.0-20190131081942-ac5582987a2f
	github.com/integration-system/isp-journal v1.4.0
	github.com/integration-system/isp-lib/v2 v2.2.0
	github.com/integration-system/isp-log v1.1.0
	github.com/json-iterator/go v1.1.9
	github.com/pkg/errors v0.9.1
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563
	github.com/spf13/cast v1.3.1
	github.com/stretchr/testify v1.4.0
	github.com/valyala/fasthttp v1.9.0
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a
	google.golang.org/grpc v1.27.1
)