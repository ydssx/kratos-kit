package admin

import (
	"errors"

	// admin "github.com/ydssx/kratos-kit/api/admin/v1"
	"github.com/ydssx/kratos-kit/common/conf"
	"github.com/ydssx/kratos-kit/internal/middleware"
	"github.com/ydssx/kratos-kit/internal/server"
	"github.com/ydssx/kratos-kit/internal/service"
	"github.com/ydssx/kratos-kit/pkg/limit"
	"github.com/ydssx/kratos-kit/pkg/logger"
	mgin "github.com/ydssx/kratos-kit/pkg/middleware/gin"
	"github.com/ydssx/kratos-kit/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewHttpServer(
	c *conf.Bootstrap,
	limiter *limit.RedisLimiter,
	adminSvc *service.AdminService,
) *http.Server {
	opts := []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			middleware.RateLimit(limiter),
			middleware.Validator(),
			middleware.TraceServer(),
			middleware.AuthAdmin(),
			middleware.LanguageMiddleware(),
		),
		http.ResponseEncoder(server.CustomizeResponseEncoder),
		http.ErrorEncoder(server.CustomizeErrorEncoder),
	}
	server := c.Server
	if server.Http.Addr != "" {
		opts = append(opts, http.Address(server.Http.Addr))
	}
	if server.Http.Timeout != nil {
		opts = append(opts, http.Timeout(server.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	// admin.RegisterAdminServiceHTTPServer(srv, adminSvc)

	gin.SetMode(gin.ReleaseMode)
	ginServer := gin.New()
	ginServer.Use(
		mgin.Logger(),
		gin.CustomRecoveryWithWriter(logger.Writer, func(c *gin.Context, err any) {
			logger.Errorf(c.Request.Context(), "panic recovered: %+v", err)
			c.AbortWithError(util.ERROR, errors.New("internal server error"))
			return
		}),
		middleware.AuthGinAdmin(),
	)

	ginServer.POST("/admin/upload", adminSvc.Upload)

	srv.HandlePrefix("/admin", ginServer)

	return srv
}
