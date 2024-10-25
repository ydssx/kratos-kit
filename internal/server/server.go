package server

import (
	"github.com/ydssx/kratos-kit/common"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(common.NewRateLimiter, NewHTTPServer, NewJobServer, NewGinServer, NewGRPCServer, NewServer)

func NewServer(httpServer *http.Server, jobServer *JobServer, grpcServer *grpc.Server) []transport.Server {
	return []transport.Server{
		jobServer,
		httpServer,
		grpcServer,
	}
}
