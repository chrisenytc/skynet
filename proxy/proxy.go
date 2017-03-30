package proxy

import (
	"context"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/gorillamux"
	"gopkg.in/kataras/iris.v6/middleware/recover"

	"github.com/chrisenytc/skynet/config"
	"github.com/chrisenytc/skynet/middlewares"
	"github.com/chrisenytc/skynet/rules"
)

type Context map[string]interface{}

func Load() {
	log.Info("Loading server configs.")

	// Create proxy
	app := iris.New()

	// Enable devlogs
	if config.IsDevelopment() {
		app.Adapt(iris.DevLogger())
	}

	log.Info("Loading HTTP request multiplexer.")

	// Enable http handler
	app.Adapt(gorillamux.New())

	// Enable recovery
	app.Use(recover.New())

	log.Info("Loading request logger.")

	// Logger
	app.UseFunc(middlewares.LoggerMiddleware)

	log.Info("Loading error handlers.")

	// Load not found handler
	app.OnError(iris.StatusNotFound, func(ctx *iris.Context) {
		middlewares.LoggerMiddleware(ctx)

		ctx.JSON(iris.StatusNotFound, map[string]string{"error": "Page not found."})
	})

	// Load error handler
	app.OnError(iris.StatusInternalServerError, func(ctx *iris.Context) {
		middlewares.LoggerMiddleware(ctx)

		ctx.JSON(iris.StatusInternalServerError, map[string]string{"error": "An internal server error occurred."})
	})

	log.Info("Loading routes.")

	// Define routes
	app.Any("/{skynet_router_fullpath:.*}", func(ctx *iris.Context) {
		path := ctx.Param("skynet_router_fullpath")
		body := &Context{}

		ctx.ReadJSON(body)

		new_url, headers := rules.Load(ctx.Method(), path)

		if new_url == "not_found" {
			ctx.EmitError(iris.StatusNotFound)
			return
		}

		client := gorequest.New()

		request := client.CustomMethod(ctx.Method(), new_url)

		request.Query(ctx.URLParams())
		request.Send(body)

		for _, h := range headers {
			if ctx.Request.Header.Get(h) != "" {
				request.Set(h, ctx.Request.Header.Get(h))
			}
		}

		request.End(func(resp gorequest.Response, body string, errs []error) {
			if len(errs) > 0 {
				ctx.EmitError(iris.StatusInternalServerError)
				return
			}

			ctx.SetContentType(resp.Header.Get("Content-Type"))
			ctx.SetStatusCode(resp.StatusCode)
			ctx.Writef(body)
		})
	})

	// Enable graceful shutdown
	app.Adapt(iris.EventPolicy{
		Interrupted: func(*iris.Framework) {
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			app.Shutdown(ctx)
		},
	})

	log.Info("Loading server.")

	log.Infof("Running on environment: %s.", config.Get().Environment)

	log.Infof("Listening on port %s.", config.Get().Port)

	// Start proxy
	app.Listen(":" + config.Get().Port)
}
