package region

import (
	"encoding/json"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"net/http"
	"strings"
)

const regionNameChina = "China"

func init() {
	registerCollector(regionNameChina, collectChina)
}

type ChinaRegion struct {
	Citycode     interface{}    `json:"citycode"`
	Adcode       string         `json:"adcode"`
	Name         string         `json:"name"`
	Center       string         `json:"center"`
	Level        string         `json:"level"`
	ChinaRegions []*ChinaRegion `json:"districts"`
}

func (c *ChinaRegion) Convert(superCoordinateCode string) *Coordinate {
	if c.Adcode == "" {
		if len(c.ChinaRegions) == 0 {
			return nil
		}
		return c.ChinaRegions[0].Convert(superCoordinateCode)
	}
	longitudeAndLatitude := strings.Split(c.Center, ",")
	subCoordinates := make([]*Coordinate, 0)
	for _, region := range c.ChinaRegions {
		subCoordinates = append(subCoordinates, region.Convert(c.Adcode))
	}
	return &Coordinate{
		SuperCoordinateCode: superCoordinateCode,
		Code:                c.Adcode,
		Name:                c.Name,
		Longitude:           longitudeAndLatitude[0],
		Latitude:            longitudeAndLatitude[1],
		SubCoordinates:      subCoordinates,
	}
}

var (
	key = kingpin.Flag("amap.key", "高德应用平台 key, 请在 https://console.amap.com/dev/key/app 申请").String()
)

func collectChina() (Region, error) {
	chinaRegion := &ChinaRegion{}
	res, err := http.Get("https://restapi.amap.com/v3/config/district?key=" + *key + "&subdistrict=3")
	if err != nil {
		return nil, err
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resBody, chinaRegion)
	if err != nil {
		return nil, err
	}
	return chinaRegion, nil
}
