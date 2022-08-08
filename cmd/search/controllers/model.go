package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/typesense/typesense-go/typesense"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Server struct {
	*echo.Echo
	Logger   *zap.Logger
	Searcher *typesense.Client
	DB       *gorm.DB
}

// HTTPError matches the echo.HTTPError format and is only used
// so swagger can display the error response type.
type HTTPError struct {
	Message string `json:"message"`
}
