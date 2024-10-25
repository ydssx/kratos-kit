package admin

import (
	"github.com/ydssx/kratos-kit/common"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	common.NewRateLimiter,
	NewServer,
	NewHttpServer,
	NewJobServer,
)

func NewServer(httpServer *http.Server, jobServer *JobServer) []transport.Server {
	return []transport.Server{
		httpServer,
		jobServer,
	}
}
