package region

import (
	"encoding/json"
	"errors"
	"github.com/gogroup/coordinate/storage"
	"github.com/morikuni/failure"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	amapKey = kingpin.Flag(
		"region.china.amap.key",
		"AMAP key, doc: [1] https://console.amap.com/dev/key/app; [2] https://lbs.amap.com/api/webservice/guide/api/district.",
	).String()
)

const regionNameChina = "china"

func init() {
	registerRegion(regionNameChina, defaultEnabled, chinaCollector, chinaSnapshot)
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

func chinaCollector() ([]*storage.Coordinate, error) {
	if *amapKey == "" {
		return nil, failure.Wrap(errors.New("need flag --amap.key"))
	}
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

func chinaSnapshot() ([]*storage.Coordinate, error) {
	c := &china{}
	file, err := ioutil.ReadFile("region/china.json")
	if err != nil {
		return nil, failure.Wrap(err)
	}
	err = json.Unmarshal(file, c)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	return c.convert(), nil
}
