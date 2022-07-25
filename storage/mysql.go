package storage

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const storageTypeMysql = "mysql"

func init() {
	registerInitializer(storageTypeMysql, initMysql)
}

type myMysql struct {
	*gorm.DB
}

var (
	dsn = kingpin.Flag(
		"storage.mysql.dsn",
		"MySQL data source name, doc: https://github.com/go-sql-driver/mysql#dsn-data-source-name.",
	).String()
)

func initMysql() (storage, error) {
	db, err := gorm.Open(mysql.Open(*dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&Coordinate{})
	if err != nil {
		return nil, err
	}
	return &myMysql{DB: db}, nil
}

func (m *myMysql) store(coordinates []*Coordinate) {
	for _, coordinate := range coordinates {
		m.Create(coordinate)
	}
}
