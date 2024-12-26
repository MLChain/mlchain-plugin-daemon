package main

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/mlchain/mlchain-plugin-daemon/internal/server"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/app"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/log"
)

func main() {
	var config app.Config

	// load env
	godotenv.Load()

	err := envconfig.Process("", &config)
	if err != nil {
		log.Panic("Error processing environment variables")
	}

	config.SetDefault()

	if err := config.Validate(); err != nil {
		log.Panic("Invalid configuration: %s", err.Error())
	}

	(&server.App{}).Run(&config)
}
