package server

import (
	"github.com/ydssx/kratos-kit/common/conf"
	"github.com/ydssx/kratos-kit/internal/middleware"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/oschwald/geoip2-golang"
)

func NewGRPCServer(c *conf.Bootstrap, geoip *geoip2.Reader) *grpc.Server {
	server := grpc.NewServer(
		grpc.Address(c.Server.Grpc.Addr),
		grpc.Timeout(c.Server.Grpc.Timeout.AsDuration()),
		grpc.Middleware(recovery.Recovery(), middleware.AuthServer(geoip)),
	)

	return server
}
