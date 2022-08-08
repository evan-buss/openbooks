package controllers

import (
	"net/url"

	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func (api *Server) ProxyRoutes() {
	searchGroup := api.Group("/multi_search")

	searchUrl, err := url.Parse("http://localhost:8108")
	if err != nil {
		api.Logger.Info("error", zap.Error(err))
	}

	targets := []*middleware.ProxyTarget{
		{
			URL: searchUrl,
		},
	}

	searchGroup.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
}
