package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (api *Server) APIRoutes() {
	api.GET("/", api.Index())
	api.GET("/api/health", api.ServiceHealth())
}

func (api *Server) Index() echo.HandlerFunc {
	return func(c echo.Context) error {
		api.Logger.Info("hello world")
		return c.String(http.StatusOK, "Hello, World!")
	}
}

// ServiceHealth godoc
// @Summary      Get the health of the service
// @Description  Get the health of the service and the underlying search provider.
// @Tags         health
// @Success      200
// @Failure      500
// @Router       /api/health [get]
func (api *Server) ServiceHealth() echo.HandlerFunc {
	return func(c echo.Context) error {
		ok, err := api.Searcher.Health(time.Second * 5)
		fmt.Println(ok, err)
		if err != nil || !ok {
			api.Logger.Error("error", zap.Error(err))
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}

		return c.String(http.StatusOK, "OK")
	}
}
