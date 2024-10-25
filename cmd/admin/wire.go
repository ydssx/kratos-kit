//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"github.com/ydssx/kratos-kit/common/conf"
	"github.com/ydssx/kratos-kit/internal/biz"
	"github.com/ydssx/kratos-kit/internal/data"
	"github.com/ydssx/kratos-kit/internal/server/admin"
	"github.com/ydssx/kratos-kit/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(ctx context.Context, c *conf.Bootstrap, logger log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(admin.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
