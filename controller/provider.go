package controller

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewSeatController,
	NewCreditController,
	NewDiscussionController,
	NewHistoryRecordController,
)
