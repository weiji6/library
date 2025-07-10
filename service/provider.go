package service

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewSeatServiceImpl,
	NewCreditServiceImpl,
	NewDiscussionServiceImpl,
)
