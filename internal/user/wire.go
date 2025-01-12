//go:build wireinject

package user

import (
	"github.com/Duke1616/ecmdb/internal/policy"
	"github.com/Duke1616/ecmdb/internal/user/internal/repository"
	"github.com/Duke1616/ecmdb/internal/user/internal/repository/dao"
	"github.com/Duke1616/ecmdb/internal/user/internal/service"
	"github.com/Duke1616/ecmdb/internal/user/internal/web"
	"github.com/Duke1616/ecmdb/internal/user/ldapx"
	"github.com/Duke1616/ecmdb/pkg/mongox"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	service.NewLdapService,
	service.NewService,
	repository.NewResourceRepository,
	dao.NewUserDao,
	web.NewHandler)

func InitModule(db *mongox.Mongo, ldapConfig ldapx.Config, policyModule *policy.Module) (*Module, error) {
	wire.Build(
		ProviderSet,
		wire.Struct(new(Module), "*"),
		wire.FieldsOf(new(*policy.Module), "Svc"),
	)
	return new(Module), nil
}
