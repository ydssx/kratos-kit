package main

import (
	"context"
	"flag"
	"time"

	"github.com/ydssx/kratos-kit/common"
	"github.com/ydssx/kratos-kit/common/conf"
	"github.com/ydssx/kratos-kit/pkg/logger"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport"
	_ "go.uber.org/automaxprocs"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "f", "./configs/config.local.yaml", "config path, eg: -conf config.yaml")
}

// main是程序的入口点。它会解析命令行参数,加载配置,初始化应用程序,并启动应用程序。
// 如果在启动过程中发生错误,它会panic。
func main() {
	flag.Parse()

	var config conf.Bootstrap
	closeConfig := conf.MustLoad(&config, flagconf)
	defer closeConfig()

	common.SetupLogger(&config)
	common.SetEnv(&config)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		time.Sleep(time.Millisecond * 10)
	}()

	application, cleanup, err := wireApp(ctx, &config, logger.DefaultLogger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// Start the application and wait for stop signal
	if err := application.Run(); err != nil {
		panic(err)
	}
}

// newApp 创建一个新的 Kratos 应用程序实例。它接收配置和服务器作为参数,
// 并根据配置来注册服务发现、追踪和指标中间件。
// 返回构建好的 Kratos 应用程序实例。
func newApp(ctx context.Context, c *conf.Bootstrap, srv ...transport.Server) *kratos.App {
	options := []kratos.Option{
		kratos.Name(c.Name),
		kratos.Context(ctx),
		kratos.Metadata(map[string]string{}),
		kratos.Server(srv...),
		kratos.BeforeStart(func(ctx context.Context) error {
			logger.Infof(ctx, "service %s is starting...", c.Name)
			return nil
		}),
	}

	return kratos.New(options...)
}
