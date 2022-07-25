package region

import (
	"encoding/json"
	"github.com/gogroup/coordinate/storage"
	"io/ioutil"
	"net/http"
	"strings"
)

const regionNameChina = "china"

func init() {
	registerCollector(regionNameChina, defaultEnabled, collectChina)
}

type chinaRegion struct {
	Citycode     interface{}    `json:"citycode"`
	Adcode       string         `json:"adcode"`
	Name         string         `json:"name"`
	Center       string         `json:"center"`
	Level        string         `json:"level"`
	ChinaRegions []*chinaRegion `json:"districts"`
}

func (c *chinaRegion) convert() []*storage.Coordinate {
	coordinates := make([]*storage.Coordinate, 0)
	var deal func(superCoordinateCode string, c *chinaRegion)
	deal = func(superCoordinateCode string, c *chinaRegion) {
		if c.Adcode == "" {
			if len(c.ChinaRegions) > 0 {
				deal(superCoordinateCode, c.ChinaRegions[0])
			}
			return
		}
		longitudeAndLatitude := strings.Split(c.Center, ",")
		coordinates = append(coordinates, &storage.Coordinate{
			SuperCoordinateCode: superCoordinateCode,
			Code:                c.Adcode,
			Name:                c.Name,
			Longitude:           longitudeAndLatitude[0],
			Latitude:            longitudeAndLatitude[1],
		})
		for _, region := range c.ChinaRegions {
			deal(c.Adcode, region)
		}
	}
	deal(storage.WorldCode, c)
	return coordinates
}

func collectChina() ([]*storage.Coordinate, error) {
	china := &chinaRegion{}
	res, err := http.Get("https://restapi.amap.com/v3/config/district?key=" + *amapKey + "&subdistrict=3")
	if err != nil {
		return nil, err
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resBody, china)
	if err != nil {
		return nil, err
	}
	return china.convert(), nil
}
