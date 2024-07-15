//go:build wireinject

package order

import (
	"context"
	"github.com/Duke1616/ecmdb/internal/order/internal/event"
	"github.com/Duke1616/ecmdb/internal/order/internal/event/consumer"
	"github.com/Duke1616/ecmdb/internal/order/internal/repository"
	"github.com/Duke1616/ecmdb/internal/order/internal/repository/dao"
	"github.com/Duke1616/ecmdb/internal/order/internal/service"
	"github.com/Duke1616/ecmdb/internal/order/internal/web"
	"github.com/Duke1616/ecmdb/internal/workflow"
	"github.com/Duke1616/ecmdb/pkg/mongox"
	"github.com/ecodeclub/mq-api"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	web.NewHandler,
	service.NewService,
	repository.NewOrderRepository,
	dao.NewOrderDAO,
)

func InitModule(q mq.MQ, db *mongox.Mongo, workflowModule *workflow.Module) (*Module, error) {
	wire.Build(
		ProviderSet,
		event.NewCreateProcessEventProducer,
		initWechatConsumer,
		InitProcessConsumer,
		wire.FieldsOf(new(*workflow.Module), "Svc"),
		wire.Struct(new(Module), "*"),
	)
	return new(Module), nil
}

func initWechatConsumer(svc service.Service, q mq.MQ) *consumer.WechatOrderConsumer {
	c, err := consumer.NewWechatOrderConsumer(svc, q)
	if err != nil {
		panic(err)
	}

	c.Start(context.Background())
	return c
}

func InitProcessConsumer(q mq.MQ, workflowSvc workflow.Service, svc service.Service) *consumer.ProcessEventConsumer {
	c, err := consumer.NewProcessEventConsumer(q, workflowSvc, svc)
	if err != nil {
		return nil
	}

	c.Start(context.Background())
	return c
}
