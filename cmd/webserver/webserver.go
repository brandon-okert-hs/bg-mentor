package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bcokert/bg-mentor/internal/pkg/handler"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	var err error
	var config *viper.Viper
	var logger *zap.SugaredLogger

	// Load Config
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("Invalid arguments. Must have exactly 1 argument, the folder containing the config file")
	}
	configFolder := args[0]
	config, err = loadConfig(configFolder)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Load Logger
	env := config.GetString("env")
	staticFileRoot := config.GetString("server.staticFileRoot")
	address := config.GetString("server.address")

	logger, err = loadLogger(env)
	if err != nil {
		log.Fatalf("Error setting up logger: %v", err)
	}
	defer logger.Sync() // write out all pending logs when server exists

	logger.Infow("Sucessfully loaded logger and config",
		"configFolder", configFolder,
		"env", env,
		"address", address,
		"staticFileRoot", staticFileRoot,
	)

	// Load Handlers and start server
	handler.SetRequestLogger(logger)
	rootHandler := &handler.RootHandler{
		StaticFileRoot: staticFileRoot,
		StatusHandler:  &handler.StatusHandler{},
		StaticHandler: &handler.StaticHandler{
			StaticFileRoot: staticFileRoot,
			FileServer:     http.FileServer(http.Dir(staticFileRoot)),
		},
	}
	log.Fatal(http.ListenAndServe(address, rootHandler))
}

func loadConfig(configFolder string) (*viper.Viper, error) {
	var config = viper.New()
	config.SetConfigName("config")
	config.AddConfigPath(configFolder)
	err := config.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func loadLogger(env string) (*zap.SugaredLogger, error) {
	var loggerInit *zap.Logger
	var err error
	switch env {
	case "production":
		loggerInit, err = zap.NewProduction()
	case "dev":
		loggerInit, err = zap.NewDevelopment()
	case "local":
		loggerInit, err = zap.NewDevelopment()
	}
	if err != nil {
		return nil, err
	}
	return loggerInit.Sugar(), nil
}
