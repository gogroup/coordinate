package storage

import (
	"errors"
	"github.com/gogroup/coordinate/region"
)

type Storage interface {
	Store(coordinates map[string]*region.Coordinate)
}

var (
	storages = make(map[string]func() (Storage, error))
)

func registerInitializer(storageType string, initializer func() (Storage, error)) {
	storages[storageType] = initializer
}

func Init(storageType string) (Storage, error) {
	for s, f := range storages {
		if s == storageType {
			return f()
		}
	}
	return nil, errors.New("no storage type: " + storageType)
}
