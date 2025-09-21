package server

import (
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/ydssx/kratos-kit/common"
	"github.com/ydssx/kratos-kit/pkg/limit"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	common.NewRateLimiter,
	wire.Bind(new(limit.Limiter), new(*limit.RedisLimiter)),
	NewHTTPServer,
	NewJobServer,
	NewGinMux,
	NewGRPCServer,
	NewServer,
)

func NewServer(
	httpServer *http.Server,
	jobServer *JobServer,
	grpcServer *grpc.Server,
) []transport.Server {
	return []transport.Server{
		jobServer,
		httpServer,
		grpcServer,
	}
}
