package server

import (
	"errors"

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

// NewGinServer 创建一个新的 Gin 服务器实例，用于处理不方便通过 proto 定义的接口，如上传接口。
func NewGinServer(commonSvc *service.CommonService, userSvc *service.UserService, geoip *geoip2.Reader) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	server := gin.New()
	server.ContextWithFallback = true
	server.Use(
		mgin.Logger(),
		gin.CustomRecoveryWithWriter(logger.Writer, func(c *gin.Context, err any) {
			logger.Errorf(c.Request.Context(), "panic recovered: %+v", err)
			c.AbortWithError(util.ERROR, errors.New("internal server error"))
			return
		}),
	)

	// Add a GET route for the API documentation
	server.GET("/docs", docsHandler)

	// Add a GET route for the Swagger UI
	// The Swagger UI is accessible at http://localhost:9000/swagger/index.html
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs")))

	server.POST("/api/upload", middleware.AuthGin(geoip), commonSvc.Upload)
	server.GET("/api/users/google-callback", userSvc.GoogleCallback)

	return server
}

func docsHandler(c *gin.Context) {
	c.Writer.Write(docs.ApiDocs)
}
