package middlewares

import (
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/kataras/iris.v6"
)

func LoggerMiddleware(ctx *iris.Context) {
	startTime := time.Now()

	w := ctx.Recorder()

	ctx.Next()

	endTime := time.Now()

	logger := logrus.New()

	log := logger.WithFields(logrus.Fields{
		"date":       startTime,
		"method":     ctx.Method(),
		"path":       ctx.Path(),
		"status":     w.StatusCode(),
		"latency":    endTime.Sub(startTime),
		"request_ip": ctx.RemoteAddr(),
		"referer":    ctx.Request.Referer(),
	})

	switch w.StatusCode() {
	case 200, 201, 202:
		log.Info("SUCCESS")
	case 400, 401, 404:
		log.Error("ERROR")
	case 500, 503:
		log.Error("FATAL")
	default:
		log.Info("REQUEST")
	}
}
