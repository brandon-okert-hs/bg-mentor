package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bcokert/bg-mentor/internal/pkg/database"
	"github.com/bcokert/bg-mentor/internal/pkg/handler"
	"github.com/go-sql-driver/mysql"
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
		"domain", config.GetString("server.domain"),
	)

	// Setup DB
	database.SetDBLogger(logger)
	database, err := database.NewDatabase("mysql", mysql.Config{
		User:                 config.GetString("database.user"),
		Passwd:               config.GetString("database.password"),
		Addr:                 fmt.Sprintf("%s:%s", config.GetString("database.host"), config.GetString("database.port")),
		DBName:               config.GetString("database.name"),
		Net:                  config.GetString("database.protocol"),
		AllowNativePasswords: true,
	})
	if err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}

	// Setup authentication config
	authConfig := handler.Auth0Config{
		Domain:          config.GetString("auth0.domain"),
		ClientID:        config.GetString("auth0.clientId"),
		ClientSecret:    config.GetString("auth0.clientSecret"),
		RedirectURLRoot: "http://" + config.GetString("server.domain") + address,
		JWTSecret:       []byte(config.GetString("jwtSecret")),
		CookieName:      config.GetString("auth0.cookieName"),
	}

	// Load Handlers and start server
	handler.SetRequestLogger(logger)
	rootHandler := &handler.RootHandler{
		StaticFileRoot: staticFileRoot,
		StatusHandler:  &handler.StatusHandler{},
		StaticHandler: &handler.StaticHandler{
			StaticFileRoot: staticFileRoot,
			FileServer:     http.FileServer(http.Dir(staticFileRoot)),
		},
		AuthHandler: &handler.AuthHandler{
			AuthConfig: &authConfig,
			DB:         database,
		},
		APIHandler: handler.Authenticated(&handler.APIHandler{
			MemberHandler:     &handler.MemberHandler{DB: database},
			UnitHandler:       &handler.UnitHandler{DB: database},
			TournamentHandler: &handler.TournamentHandler{DB: database},
			DABEntryHandler:   &handler.DABEntryHandler{DB: database},
		}, &authConfig),
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
