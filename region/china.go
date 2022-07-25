package region

import (
	"encoding/json"
	"github.com/gogroup/coordinate/storage"
	"github.com/morikuni/failure"
	"io/ioutil"
	"net/http"
	"strings"
)

const regionNameChina = "china"

func init() {
	registerCollector(regionNameChina, defaultEnabled, collectChina)
}

type china struct {
	Citycode   interface{} `json:"citycode"`
	Adcode     string      `json:"adcode"`
	Name       string      `json:"name"`
	Center     string      `json:"center"`
	Level      string      `json:"level"`
	SubRegions []*china    `json:"districts"`
}

func (c *china) convert() []*storage.Coordinate {
	coordinates := make([]*storage.Coordinate, 0)
	var deal func(superCoordinateCode string, c *china)
	deal = func(superCoordinateCode string, c *china) {
		if c.Adcode == "" {
			if len(c.SubRegions) > 0 {
				deal(superCoordinateCode, c.SubRegions[0])
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
		for _, region := range c.SubRegions {
			deal(c.Adcode, region)
		}
	}
	deal(storage.WorldCode, c)
	return coordinates
}

func collectChina() ([]*storage.Coordinate, error) {
	c := &china{}
	res, err := http.Get("https://restapi.amap.com/v3/config/district?key=" + *amapKey + "&subdistrict=3")
	if err != nil {
		return nil, failure.Wrap(err)
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	err = json.Unmarshal(resBody, c)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	return c.convert(), nil
}
