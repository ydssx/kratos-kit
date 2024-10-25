package service

import (
	"github.com/ydssx/kratos-kit/common"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	common.NewAsynqClient,
	common.NewAsynqInspector,
	NewUserService,
	NewJobService,
	NewCommonService,
	NewAdminService,
)
