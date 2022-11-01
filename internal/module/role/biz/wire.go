package rolebiz

import (
	"github.com/google/wire"
)

// ProviderBizSet is Biz providers.
var ProviderBizSet = wire.NewSet(
	NewBiz,
)
