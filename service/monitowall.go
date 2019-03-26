package service

import (
	"fmt"

	"github.com/jsdidierlaurent/monitowall/handlers/info"

	"github.com/jsdidierlaurent/echo-middleware/cache"

	"github.com/jsdidierlaurent/monitowall/config"
	"github.com/jsdidierlaurent/monitowall/handlers"
	"github.com/jsdidierlaurent/monitowall/handlers/ping"
	"github.com/jsdidierlaurent/monitowall/middlewares"
	"github.com/jsdidierlaurent/monitowall/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start(config *config.Config, buildInfo *config.BuildInfo) {
	router := echo.New()

	//  ----- Middlewares -----
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middlewares.ConfigMiddleware(config, buildInfo))

	cm := middlewares.NewCacheMiddleware(config) // Used as Handler wrapper in routes
	router.Use(cm.DownstreamStoreMiddleware())

	//  ----- CORS -----
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	// ----- Errors Handler -----
	router.HTTPErrorHandler = handlers.HTTPErrorHandler

	// ----- Routes -----
	v1 := router.Group("/api/v1")

	// ------------- INFO ------------- //
	v1.GET("/info", cm.UpstreamCacheHandlerWithExpiration(cache.NEVER, info.GetInfo))
	// ------------- PING ------------- //
	pingHandler := ping.NewHandler(models.NewPingModel())
	v1.GET("/ping", cm.UpstreamCacheHandler(pingHandler.GetPing))

	// Start service
	router.Logger.Fatal(router.Start(fmt.Sprintf(":%d", config.Port)))
}