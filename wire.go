//go:build wireinject
// +build wireinject

package main

import (
	"library/config"
	"library/controller"
	"library/router"
	"library/service"

	"github.com/google/wire"
)

func InitApp(configPath string) (*router.App, error) {
	wire.Build(
		config.ProviderSet,
		service.ProviderSet,
		controller.ProviderSet,
		router.ProviderSet,
		wire.Bind(new(service.SeatService), new(*service.SeatServiceImpl)),
		wire.Bind(new(service.CreditPoints), new(*service.CreditServiceImpl)),
		wire.Bind(new(service.Discussion), new(*service.DiscussionImpl)),
	)
	return &router.App{}, nil
}
