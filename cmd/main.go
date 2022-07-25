package main

import (
	"fmt"
	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gogroup/coordinate/region"
	"github.com/gogroup/coordinate/storage"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := getLogger()

	regionCoordinates, err := region.Collect(logger)
	if err != nil {
		logger.Error(fmt.Sprintf("%+v", err))
		return
	}

	err = storage.Store(logger, regionCoordinates)
	if err != nil {
		logger.Error(fmt.Sprintf("%+v", err))
		return
	}
}

func getLogger() *log.Logger {
	logger := log.New()
	logger.SetFormatter(&formatter.Formatter{
		TimestampFormat: "2006-01-02 | 15:04:05.000",
	})
	return logger
}
