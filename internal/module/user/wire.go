package usermodule

import (
	"github.com/google/wire"

	userbiz "github.com/xdorro/golang-grpc-base-project/internal/module/user/biz"
	userservice "github.com/xdorro/golang-grpc-base-project/internal/module/user/service"
)

// ProviderModuleSet is Module providers.
var ProviderModuleSet = wire.NewSet(
	userbiz.ProviderBizSet,
	userservice.ProviderServiceSet,
)
