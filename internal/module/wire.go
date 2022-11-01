package module

import (
	"github.com/google/wire"

	authmodule "github.com/xdorro/golang-grpc-base-project/internal/module/auth"
	usermodule "github.com/xdorro/golang-grpc-base-project/internal/module/user"
)

// ProviderModuleSet is Module providers.
var ProviderModuleSet = wire.NewSet(
	usermodule.ProviderModuleSet,
	authmodule.ProviderModuleSet,
)
