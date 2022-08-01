package main

import (
	"mp-generate-test-data/config"
	"mp-generate-test-data/internal/services"
	"mp-generate-test-data/internal/stores"
)

func registerService(cfg *config.Config) *services.Service {
	db := mustConnectPostgres(cfg)

	mainStore := stores.NewMainStore(db)
	ordersStore := stores.NewOrdersStore(db)
	orderAssetsStore := stores.NewOrderAssetsStore(db)
	matchedOrdersStore := stores.NewMatchedOrdersStore(db)

	return services.New(cfg, mainStore, ordersStore, orderAssetsStore, matchedOrdersStore)
}
