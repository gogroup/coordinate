package storage

import (
	"fmt"
	"github.com/gogroup/coordinate/region"
	"gopkg.in/alecthomas/kingpin.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const storageTypeMysql = "mysql"

func init() {
	registerInitializer(storageTypeMysql, initMysql)
}

type Mysql struct {
	*gorm.DB
}

var (
	dsn = kingpin.Flag(
		"storage.mysql.dsn",
		"MySQL data source name, doc: https://github.com/go-sql-driver/mysql#dsn-data-source-name.",
	).String()
)

func initMysql() (Storage, error) {
	db, err := gorm.Open(mysql.Open(*dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&region.Coordinate{})
	if err != nil {
		return nil, err
	}
	return &Mysql{DB: db}, nil
}

// Store TODO 这里可以优化一下，改成让 Store 去做，所有的统一控制日志
func (m *Mysql) Store(coordinates map[string]*region.Coordinate) {
	fmt.Println("Start write to mysql.")
	for s, coordinate := range coordinates {
		fmt.Printf("- Writing %s...\n", s)
		w(m.DB, coordinate)
		fmt.Println("- Done!")
	}
}

func w(db *gorm.DB, coordinate *region.Coordinate) {
	db.Create(coordinate)
	for _, subCoordinate := range coordinate.SubCoordinates {
		w(db, subCoordinate)
	}
}
