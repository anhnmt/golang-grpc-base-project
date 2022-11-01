package module

import (
	"github.com/google/wire"

	authmodule "github.com/xdorro/golang-grpc-base-project/internal/module/auth"
	rolemodule "github.com/xdorro/golang-grpc-base-project/internal/module/role"
	usermodule "github.com/xdorro/golang-grpc-base-project/internal/module/user"
)

// ProviderModuleSet is Module providers.
var ProviderModuleSet = wire.NewSet(
	usermodule.ProviderModuleSet,
	authmodule.ProviderModuleSet,
	rolemodule.ProviderModuleSet,
)
