package main

import (
	"context"

	"github.com/evan-buss/openbooks/cmd/search/controllers"
	_ "github.com/evan-buss/openbooks/cmd/search/docs"
	"github.com/evan-buss/openbooks/cmd/search/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/typesense/typesense-go/typesense"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/mattn/go-sqlite3"
)

// @title OpenBooks Search Service
// @version 1.0
// @description This service provides advanced search capabilities to OpenBooks as well as a management interface.

// @contact.name OpenBooks
// @contact.url https://github.com/evan-buss/openbooks

// @license.name MIT
// @license.url https://mit-license.org/

// @Accept       json
// @Produce      json
func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	sugared := logger.Sugar()

	db := InitDbClient(sugared)

	initContext := context.Background()

	err = model.EnsureDefaultUser(db, initContext)
	if err != nil {
		sugared.Fatalf("error created / checking default user %v", err)
	}

	e := echo.New()
	e.HideBanner = true

	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	searcher := typesense.NewClient(
		typesense.WithServer("http://localhost:8108"),
		typesense.WithAPIKey("Hu52dwsas2AdxdE"))

	server := &controllers.Server{Echo: e, Logger: logger, Searcher: searcher, DB: db}
	server.GET("/swagger/*", echoSwagger.WrapHandler)
	server.APIRoutes()
	server.AuthRoutes()
	server.ProxyRoutes()
	server.ScheduleRoutes()

	e.Logger.Fatal(e.Start(":5699"))
}

func InitDbClient(logger *zap.SugaredLogger) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&model.User{}, &model.Schedule{}, &model.History{})

	return db.Debug()
}
