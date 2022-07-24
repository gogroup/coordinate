package storage

import (
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

// Init TODO 检查 sink 参数的合法性
func Init(storageType string) (Storage, error) {
	return storages[storageType]()
}
