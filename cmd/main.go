package main

import (
	"github.com/gogroup/coordinate/region"
	"github.com/gogroup/coordinate/storage"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	storageType = kingpin.Flag(
		"storage.type",
		"Type of data sink, support: mysql.",
	).Required().String()
)

func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	s, err := storage.Init(*storageType)
	if err != nil {
		panic(err)
	}
	coordinates, err := region.Collect()
	if err != nil {
		panic(err)
	}
	s.Store(coordinates)
}
