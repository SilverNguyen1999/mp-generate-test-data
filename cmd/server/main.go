package main

import (
	"mp-generate-test-data/config"
)

func main() {
	cfg := config.Load()
	svc := registerService(cfg)

	svc.Run()
}
