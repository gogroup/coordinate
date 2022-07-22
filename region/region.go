package region

import (
	"fmt"
	"gorm.io/gorm"
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
	collectors = make(map[string]func() (Region, error))
	regions    = make(map[string]Region)
)

func registerCollector(region string, collector func() (Region, error)) {
	collectors[region] = collector
}

func Collect(db *gorm.DB) error {
	err := db.AutoMigrate(&Coordinate{})
	if err != nil {
		return err
	}
	fmt.Println("Start collect regions.")
	err = collect()
	if err != nil {
		return err
	}
	fmt.Println("Start write to db.")
	write(db)
	return nil
}

func collect() error {
	for s, f := range collectors {
		fmt.Printf("- Collecting %s...\n", s)
		region, err := f()
		if err != nil {
			return err
		}
		regions[s] = region
		fmt.Println("- Done!")
	}
	return nil
}

func write(db *gorm.DB) {
	for s, region := range regions {
		fmt.Printf("- Writing %s...\n", s)
		coordinate := region.Convert(topCode)
		w(db, coordinate)
		fmt.Println("- Done!")
	}
}

func w(db *gorm.DB, coordinate *Coordinate) {
	db.Create(coordinate)
	for _, subCoordinate := range coordinate.SubCoordinates {
		w(db, subCoordinate)
	}
}
