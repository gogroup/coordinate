package storage

import (
	"errors"
	"fmt"
	"github.com/morikuni/failure"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	storageType = kingpin.Flag(
		"storage.type",
		"Type of data sink, support: mysql.",
	).Required().String()
)

const WorldCode = "0"

type Coordinate struct {
	SuperCoordinateCode string
	Code                string
	Name                string
	Longitude           string
	Latitude            string
}

type storage interface {
	store(coordinates []*Coordinate)
}

var (
	storages = make(map[string]func() (storage, error))
)

func registerInitializer(storageType string, initializer func() (storage, error)) {
	storages[storageType] = initializer
}

func create() (storage, error) {
	for s, f := range storages {
		if s == *storageType {
			return f()
		}
	}
	return nil, failure.Wrap(errors.New("no storage type: " + *storageType))
}

func Store(logger *log.Logger, m map[string][]*Coordinate) error {
	s, err := create()
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Start write to %s.", *storageType))
	for regionName, coordinates := range m {
		logger.Info(fmt.Sprintf("- Writing %s...", regionName))
		s.store(coordinates)
		logger.Info("- Done!")
	}
	return nil
}
