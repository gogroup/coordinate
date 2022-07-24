package main

import (
	"github.com/gogroup/coordinate/region"
	"github.com/gogroup/coordinate/storage"
	"gopkg.in/alecthomas/kingpin.v2"
)

// TODO 参数检查，不传参无法执行，输出 -h
var (
	sink = kingpin.Flag(
		"sink",
		"持久化数据的系统类型，目前支持的类型有 [mysql]",
	).Required().String()
)

func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	coordinates, err := region.Collect()
	if err != nil {
		panic(err)
	}
	s, err := storage.Init(*sink)
	if err != nil {
		panic(err)
	}
	s.Store(coordinates)
}
