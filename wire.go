//go:build wireinject

package main

import (
	"github.com/dadaxiaoxiao/go-pkg/customserver"
	"github.com/dadaxiaoxiao/search/internal/events"
	"github.com/dadaxiaoxiao/search/internal/grpc"
	"github.com/dadaxiaoxiao/search/internal/repository"
	"github.com/dadaxiaoxiao/search/internal/repository/dao"
	"github.com/dadaxiaoxiao/search/internal/service"
	"github.com/dadaxiaoxiao/search/ioc"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitEtcdClient,
	ioc.InitRedis,
	ioc.InitLogger,
	ioc.InitOTEL,
	ioc.InitESClient,
	ioc.InitKafka,
	ioc.NewConsumers)

var serviceProvider = wire.NewSet(
	dao.NewTagESDAO,
	dao.NewAnyElasticDAO,
	dao.NewUserElasticDAO,
	dao.NewArticleElasticDAO,
	repository.NewAnyRepository,
	repository.NewUserRepository,
	repository.NewArticleRepository,
	service.NewSyncService,
	service.NewSearchService,
)

var consumerProvider = wire.NewSet(
	events.NewSyncAnyDataEventConsumer,
	events.NewSyncUserConsumer,
	events.NewSyncArticleConsumer,
)

func InitApp() *customserver.App {
	wire.Build(
		thirdProvider,
		serviceProvider,
		grpc.NewSearchServiceServer,
		grpc.NewSyncServiceServer,
		ioc.InitGRPCServer,
		consumerProvider,
		// wire 帮忙构造
		wire.Struct(new(customserver.App), "GRPCServer", "Consumers"),
	)
	return new(customserver.App)
}
