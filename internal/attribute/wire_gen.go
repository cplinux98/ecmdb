// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package attribute

import (
	"github.com/Duke1616/ecmdb/internal/attribute/internal/repository"
	"github.com/Duke1616/ecmdb/internal/attribute/internal/repository/dao"
	"github.com/Duke1616/ecmdb/internal/attribute/internal/service"
	"github.com/Duke1616/ecmdb/internal/attribute/internal/web"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

// Injectors from wire.go:

func InitModule(db *mongo.Client) (*Module, error) {
	attributeDAO := dao.NewAttributeDAO(db)
	attributeRepository := repository.NewAttributeRepository(attributeDAO)
	service := NewService(attributeRepository)
	handler := web.NewHandler(service)
	module := &Module{
		Svc: service,
		Hdl: handler,
	}
	return module, nil
}

// wire.go:

var ProviderSet = wire.NewSet(web.NewHandler, repository.NewAttributeRepository, dao.NewAttributeDAO)

func NewService(repo repository.AttributeRepository) Service {
	return service.NewService(repo)
}
