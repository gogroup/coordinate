package main

import (
	"github.com/gogroup/coordinate/region"
	"github.com/gogroup/coordinate/storage"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	regionCoordinates, err := region.Collect()
	if err != nil {
		panic(err)
	}
	err = storage.Store(regionCoordinates)
	if err != nil {
		panic(err)
	}
}
