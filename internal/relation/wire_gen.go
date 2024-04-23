// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package relation

import (
	"github.com/Duke1616/ecmdb/internal/attribute"
	"github.com/Duke1616/ecmdb/internal/model"
	"github.com/Duke1616/ecmdb/internal/relation/internal/repository"
	"github.com/Duke1616/ecmdb/internal/relation/internal/repository/dao"
	"github.com/Duke1616/ecmdb/internal/relation/internal/service"
	"github.com/Duke1616/ecmdb/internal/relation/internal/web"
	"github.com/Duke1616/ecmdb/internal/resource"
	"github.com/Duke1616/ecmdb/pkg/mongox"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitModule(db *mongox.Mongo, attributeModel *attribute.Module, resourceModel *resource.Module, modelModule *model.Module) (*Module, error) {
	relationResourceDAO := dao.NewRelationResourceDAO(db)
	relationResourceRepository := repository.NewRelationResourceRepository(relationResourceDAO)
	relationResourceService := service.NewRelationResourceService(relationResourceRepository)
	relationModelDAO := dao.NewRelationModelDAO(db)
	relationModelRepository := repository.NewRelationModelRepository(relationModelDAO)
	relationModelService := service.NewRelationModelService(relationModelRepository)
	relationTypeDAO := dao.NewRelationTypeDAO(db)
	relationTypeRepository := repository.NewRelationTypeRepository(relationTypeDAO)
	relationTypeService := service.NewRelationTypeService(relationTypeRepository)
	serviceService := attributeModel.Svc
	service2 := resourceModel.Svc
	relationResourceHandler := web.NewRelationResourceHandler(relationResourceService, serviceService, service2)
	service3 := modelModule.Svc
	relationModelHandler := web.NewRelationModelHandler(relationModelService, service3)
	relationTypeHandler := web.NewRelationTypeHandler(relationTypeService)
	module := &Module{
		RRSvc: relationResourceService,
		RMSvc: relationModelService,
		RTSvc: relationTypeService,
		RRHdl: relationResourceHandler,
		RMHdl: relationModelHandler,
		RTHdl: relationTypeHandler,
	}
	return module, nil
}

// wire.go:

var ProviderSet = wire.NewSet(web.NewRelationResourceHandler, web.NewRelationModelHandler, web.NewRelationTypeHandler, service.NewRelationResourceService, service.NewRelationModelService, service.NewRelationTypeService, repository.NewRelationModelRepository, repository.NewRelationResourceRepository, repository.NewRelationTypeRepository, dao.NewRelationModelDAO, dao.NewRelationResourceDAO, dao.NewRelationTypeDAO)
