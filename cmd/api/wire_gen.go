// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ydssx/kratos-kit/common"
	"github.com/ydssx/kratos-kit/common/conf"
	"github.com/ydssx/kratos-kit/internal/biz"
	"github.com/ydssx/kratos-kit/internal/data"
	"github.com/ydssx/kratos-kit/internal/server"
	"github.com/ydssx/kratos-kit/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(ctx context.Context, c *conf.Bootstrap, logger log.Logger) (*kratos.App, func(), error) {
	wsService := common.NewWsService(ctx, logger)
	reader, cleanup := common.NewGeoipDB(c)
	client, err := common.NewRedisCLient(c)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	redisLimiter := common.NewRateLimiter(client)
	googleCloudStorage, cleanup2 := common.NewGoogleCloudStorage(c)
	db, err := common.NewMysqlDB(c)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	dataData, cleanup3, err := data.NewData(logger, client, db)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	transaction := data.NewTransaction(dataData)
	userRepo := data.NewUserRepo(dataData, logger)
	cache := data.NewRedisCache(client)
	bizUserRepo := data.NewUserRepoCacheDecorator(userRepo, cache)
	commonUseCase := biz.NewCommonUseCase(transaction, googleCloudStorage, bizUserRepo)
	uploadUseCase := biz.NewUploadUseCase(googleCloudStorage, c, commonUseCase)
	commonService := service.NewCommonService(uploadUseCase, commonUseCase)
	redisLocker := common.NewRedisLocker(client)
	config := common.InitGoogleOAuth(c)
	email := common.NewEmail(c)
	userUseCase := biz.NewUserUseCase(bizUserRepo, logger, transaction, commonUseCase, redisLocker, config, cache, email)
	userService := service.NewUserService(userUseCase)
	engine := server.NewGinServer(commonService, userService, reader)
	httpServer := server.NewHTTPServer(c, wsService, reader, redisLimiter, engine, userService)
	usecaseSet := biz.NewUsecaseSet(userUseCase, uploadUseCase)
	jobServer := server.NewJobServer(c, usecaseSet)
	grpcServer := server.NewGRPCServer(c, reader)
	v := server.NewServer(httpServer, jobServer, grpcServer)
	app := newApp(ctx, c, v...)
	return app, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}