package storage

import (
	"github.com/morikuni/failure"
	"gopkg.in/alecthomas/kingpin.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dsn = kingpin.Flag(
		"storage.mysql.dsn",
		"MySQL data source name, doc: https://github.com/go-sql-driver/mysql#dsn-data-source-name.",
	).String()
)

const storageTypeMysql = "mysql"

func init() {
	registerInitializer(storageTypeMysql, initMysql)
}

type myMysql struct {
	*gorm.DB
}

func initMysql() (storage, error) {
	db, err := gorm.Open(mysql.Open(*dsn), &gorm.Config{})
	if err != nil {
		return nil, failure.Wrap(err)
	}
	err = db.AutoMigrate(&Coordinate{})
	if err != nil {
		return nil, failure.Wrap(err)
	}
	return &myMysql{DB: db}, nil
}

func (m *myMysql) store(coordinates []*Coordinate) {
	for _, coordinate := range coordinates {
		m.Create(coordinate)
	}
}
