package storage

import (
	"errors"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
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
	storageType = kingpin.Flag(
		"storage.type",
		"Type of data sink, support: mysql.",
	).Required().String()
)

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
	return nil, errors.New("no storage type: " + *storageType)
}

func Store(m map[string][]*Coordinate) error {
	s, err := create()
	if err != nil {
		return err
	}
	fmt.Printf("Start write to %s.\n", *storageType)
	for regionName, coordinates := range m {
		fmt.Printf("- Writing %s...\n", regionName)
		s.store(coordinates)
		fmt.Println("- Done!")
	}
	return nil
}
