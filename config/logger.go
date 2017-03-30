package config

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

func LoadLogger() {
	logType := os.Getenv("LOG_TYPE")
	logLevel := os.Getenv("LOG_LEVEL")

	if logType == "pretty" {
		log.SetFormatter(&log.TextFormatter{})
	}

	if logType == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	switch logLevel {
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
