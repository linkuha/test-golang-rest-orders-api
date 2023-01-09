package main

import (
	"flag"
	"fmt"
	"github.com/linkuha/test-golang-rest-orders-api/config"
	"github.com/linkuha/test-golang-rest-orders-api/internal/app"
)

func main() {
	flags := initFlags()
	cfg := config.InitConfig(flags)

	app.Run(cfg)
}

func initFlags() config.AppFlags {
	flags := config.AppFlags{}
	flag.Var(
		&flags.LogDir,
		"logdir",
		fmt.Sprintf("Path to log directory (default: %s, can be replaced by APP_LOG_DIR env)", config.DefaultLogDir),
	)
	flag.Var(
		&flags.LogLevel,
		"loglevel",
		fmt.Sprintf("Log level (default: %s)", config.DefaultLogLevel),
	)
	flag.Var(
		&flags.ConfigPath,
		"config",
		fmt.Sprintf("Path to config (default: %s, can be replaced by APP_CONFIG_PATH env)", config.DefaultConfigPath),
	)
	flag.Parse()
	return flags
}
