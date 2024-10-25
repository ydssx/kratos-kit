package run

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	GRPC_TIMEOUT int
	ENV          string
	ADDR         string
)

func getAddr(env string) string {
	switch env {
	case "local":
		return "127.0.0.1:8000"
	default:
		return "127.0.0.1:8000"
	}
}

var token = "3f047022-1f1f-11ef-adda-1e8020115572"

type ClientSet struct {
	conn *grpc.ClientConn
}

func NewClientSet(ctx context.Context) *ClientSet {
	endpoint := getAddr(ENV)
	if ADDR != "" {
		endpoint = ADDR
	}
	conn, err := kgrpc.DialInsecure(ctx,
		kgrpc.WithEndpoint(endpoint),
		kgrpc.WithTimeout(time.Second*time.Duration(GRPC_TIMEOUT)),
		kgrpc.WithMiddleware(authClient(token)),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to dial server: %v", err))
	}

	return &ClientSet{conn}
}

func (c *ClientSet) Close() {
	err := c.conn.Close()
	if err != nil {
		slog.Error("failed to close conn to %s: %v", c.conn.Target(), err)
	}
}

func authClient(tokenStr string) middleware.Middleware {
	return func(h middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if clientContext, ok := transport.FromClientContext(ctx); ok {
				clientContext.RequestHeader().Set("Authorization", fmt.Sprintf("Bearer %s", tokenStr))
				return h(ctx, req)
			}
			return nil, errors.Unauthorized("auth", "请登录")
		}
	}
}

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "run a service with grpc",
}

func init() {
	RunCmd.PersistentFlags().StringVarP(&ENV, "env", "e", "local", "the environment of the application, should be one of 'local', 'test', 'prod'")
	RunCmd.PersistentFlags().IntVarP(&GRPC_TIMEOUT, "timeout", "t", 60, "the timeout of the gRPC call in seconds")
	RunCmd.PersistentFlags().StringVarP(&ADDR, "addr", "a", "", "the address of the gRPC server")
}
