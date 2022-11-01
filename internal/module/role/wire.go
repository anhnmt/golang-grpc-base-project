package rolemodule

import (
	"github.com/google/wire"

	rolebiz "github.com/xdorro/golang-grpc-base-project/internal/module/role/biz"
	roleservice "github.com/xdorro/golang-grpc-base-project/internal/module/role/service"
)

// ProviderModuleSet is Module providers.
var ProviderModuleSet = wire.NewSet(
	rolebiz.ProviderBizSet,
	roleservice.ProviderServiceSet,
)
