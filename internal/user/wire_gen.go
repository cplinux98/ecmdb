// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package user

import (
	"github.com/Duke1616/ecmdb/internal/user/internal/repostory"
	"github.com/Duke1616/ecmdb/internal/user/internal/repostory/dao"
	"github.com/Duke1616/ecmdb/internal/user/internal/service"
	"github.com/Duke1616/ecmdb/internal/user/internal/web"
	"github.com/Duke1616/ecmdb/internal/user/ldapx"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

// Injectors from wire.go:

func InitModule(db *mongo.Client, ldapConfig ldapx.Config) (*Module, error) {
	userDAO := dao.NewUserDao(db)
	userRepository := repostory.NewResourceRepository(userDAO)
	serviceService := service.NewService(userRepository)
	ldapService := service.NewLdapService(ldapConfig)
	handler := web.NewHandler(serviceService, ldapService)
	module := &Module{
		Hdl: handler,
	}
	return module, nil
}

// wire.go:

var ProviderSet = wire.NewSet(service.NewLdapService, service.NewService, repostory.NewResourceRepository, dao.NewUserDao, web.NewHandler)