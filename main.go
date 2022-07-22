package main

import (
	"github.com/gogroup/coordinate/region"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TODO
var (
	dsn = kingpin.Flag("dsn", "Mysql 链接, 参考文档: https://github.com/go-sql-driver/mysql#dsn-data-source-name").String()
)

func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	db, err := gorm.Open(mysql.Open(*dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = region.Collect(db)
	if err != nil {
		panic(err)
	}
}
