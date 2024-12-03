package server

import (
	"errors"
	"net/http/pprof"

	"github.com/ydssx/kratos-kit/common/conf"
	"github.com/ydssx/kratos-kit/docs"
	"github.com/ydssx/kratos-kit/internal/middleware"
	"github.com/ydssx/kratos-kit/internal/service"
	"github.com/ydssx/kratos-kit/pkg/logger"
	mgin "github.com/ydssx/kratos-kit/pkg/middleware/gin"
	"github.com/ydssx/kratos-kit/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewGinMux 创建一个新的 Gin 路由，用于处理不方便通过 proto 定义的接口，如上传接口。
func NewGinMux(
	c *conf.Bootstrap,
	geoip *geoip2.Reader,
	commonSvc *service.CommonService,
	userSvc *service.UserService,
	aiSvc *service.AIService,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	mux := gin.New()
	mux.ContextWithFallback = true
	mux.Use(
		mgin.Logger(),
		gin.CustomRecoveryWithWriter(logger.Writer, func(c *gin.Context, err any) {
			logger.Errorf(c.Request.Context(), "panic recovered: %+v", err)
			c.AbortWithError(util.ERROR, errors.New("internal server error"))
			return
		}),
	)

	// Add a GET route for the API documentation
	mux.GET("/docs", gin.BasicAuth(gin.Accounts{"admin": "admin"}), docsHandler)

	// Add a GET route for the Swagger UI
	// The Swagger UI is accessible at http://localhost:9000/swagger/index.html
	mux.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs")))
	if c.Server.EnablePprof {
		mux.GET("/debug/pprof/", gin.WrapF(pprof.Index))
		mux.GET("/debug/pprof/cmdline", gin.WrapF(pprof.Cmdline))
		mux.GET("/debug/pprof/profile", gin.WrapF(pprof.Profile))
		mux.POST("/debug/pprof/symbol", gin.WrapF(pprof.Symbol))
		mux.GET("/debug/pprof/trace", gin.WrapF(pprof.Trace))
		mux.GET("/debug/pprof/allocs", gin.WrapH(pprof.Handler("allocs")))
		mux.GET("/debug/pprof/block", gin.WrapH(pprof.Handler("block")))
		mux.GET("/debug/pprof/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		mux.GET("/debug/pprof/heap", gin.WrapH(pprof.Handler("heap")))
		mux.GET("/debug/pprof/mutex", gin.WrapH(pprof.Handler("mutex")))
		mux.GET("/debug/pprof/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}

	mux.POST("/api/upload", middleware.AuthGin(geoip), commonSvc.Upload)
	mux.GET("/api/users/google-callback", userSvc.GoogleCallback)
	mux.Any("/api/v1/ai/chat", aiSvc.Chat)

	return mux
}

func docsHandler(c *gin.Context) {
	c.Writer.Write(docs.ApiDocs)
}
