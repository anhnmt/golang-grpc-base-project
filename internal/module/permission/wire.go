package permissionmodule

import (
	"github.com/google/wire"

	permissionbiz "github.com/xdorro/golang-grpc-base-project/internal/module/permission/biz"
	permissionservice "github.com/xdorro/golang-grpc-base-project/internal/module/permission/service"
)

// ProviderModuleSet is Module providers.
var ProviderModuleSet = wire.NewSet(
	permissionbiz.ProviderBizSet,
	permissionservice.ProviderServiceSet,
)
