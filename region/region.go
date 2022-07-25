package region

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
)

const topCode = "0"

type Coordinate struct {
	SuperCoordinateCode string
	Code                string
	Name                string
	Longitude           string
	Latitude            string
	SubCoordinates      []*Coordinate `gorm:"-:all"`
}

type Region interface {
	Convert(superCoordinateCode string) *Coordinate
}

var (
	amapKey = kingpin.Flag(
		"amap.key",
		"AMAP key, doc: https://console.amap.com/dev/key/app",
	).Required().String()
)

var (
	collectors = make(map[string]func() (Region, error))
)

func registerCollector(regionName string, collector func() (Region, error)) {
	collectors[regionName] = collector
}

func Collect() (map[string]*Coordinate, error) {
	fmt.Println("Start collect coordinates.")
	coordinates, err := collect()
	if err != nil {
		return nil, err
	}
	return coordinates, nil
}

func collect() (map[string]*Coordinate, error) {
	coordinates := make(map[string]*Coordinate)
	for s, f := range collectors {
		fmt.Printf("- Collecting %s...\n", s)
		region, err := f()
		if err != nil {
			return nil, err
		}
		coordinates[s] = region.Convert(topCode)
		fmt.Println("- Done!")
	}
	return coordinates, nil
}
