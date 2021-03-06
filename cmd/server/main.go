package main

import (
	"fmt"
	"log"

	"github.com/eriktate/toggle/env"
	"github.com/eriktate/toggle/http"
	"github.com/eriktate/toggle/pg"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func buildLogger() *zap.Logger {
	logCfg := zap.NewDevelopmentEncoderConfig()
	logCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(logCfg),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
}

func run() error {
	log := buildLogger()
	defer log.Sync()

	// load configs
	host := env.GetString("TOGGLE_HOST", "localhost")
	port := env.GetUint("TOGGLE_PORT", 9001)

	// initialize postgres
	db, err := pg.NewClient(pg.ConfigFromEnv("TOGGLE"))
	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", err)
	}

	// initialize services to be used
	services := http.Services{
		// ToggleService: mock.NewToggleService(nil),
		ToggleService: db,
	}

	// build server
	cfg := http.Config{
		Host:     host,
		Port:     port,
		Logger:   log,
		Services: services,
	}

	log.Info("starting server", zap.String("host", host), zap.Uint("port", port))
	return http.Listen(cfg)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
