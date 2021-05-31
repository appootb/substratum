module github.com/appootb/substratum

go 1.14

require (
	github.com/appootb/grc v0.0.0-20210413135854-1868c3b57f72
	github.com/appootb/protobuf/go v0.0.0-20210318030911-d42bd07a0a71
	github.com/elastic/go-elasticsearch/v6 v6.8.10
	github.com/elastic/go-elasticsearch/v7 v7.9.0
	github.com/gbrlsnchs/jwt/v3 v3.0.0-rc.2
	github.com/go-redis/redis/v8 v8.2.2
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.5
	github.com/jinzhu/gorm v1.9.12
	github.com/prometheus/client_golang v1.7.1
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.23.1-0.20200526195155-81db48ad09cc
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
