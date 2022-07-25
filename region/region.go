package region

import (
	"fmt"
	"github.com/gogroup/coordinate/storage"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	ddr = kingpin.Flag(
		"region.disable-defaults",
		"Set all regions to disabled by default.",
	).Default("false").Bool()
	amapKey = kingpin.Flag(
		"amap.key",
		"AMAP key, doc: https://console.amap.com/dev/key/app",
	).Required().String() // amapKey 用于获取 china 地区的数据
)

const (
	defaultEnabled  = true
	defaultDisabled = false
)

var (
	regionCollectors = make(map[string]func() ([]*storage.Coordinate, error))
	regionState      = make(map[string]*bool)
	forcedRegions    = map[string]bool{} // forcedRegions which have been explicitly enabled or disabled
)

func registerCollector(regionName string, isDefaultEnabled bool, collector func() ([]*storage.Coordinate, error)) {
	var helpDefaultState string
	if isDefaultEnabled {
		helpDefaultState = "enabled"
	} else {
		helpDefaultState = "disabled"
	}

	flagName := fmt.Sprintf("region.%s", regionName)
	flagHelp := fmt.Sprintf("Enable %s region (default: %s).", regionName, helpDefaultState)
	defaultValue := fmt.Sprintf("%v", isDefaultEnabled)

	flag := kingpin.Flag(flagName, flagHelp).Default(defaultValue).Action(regionFlagAction(regionName)).Bool()
	regionState[regionName] = flag

	regionCollectors[regionName] = collector
}

// regionFlagAction generates a new action function for the given region
// to track whether it has been explicitly enabled or disabled from the command line.
// A new action function is needed for each region flag because the ParseContext
// does not contain information about which flag called the action.
// See: https://github.com/alecthomas/kingpin/issues/294
func regionFlagAction(regionName string) func(ctx *kingpin.ParseContext) error {
	return func(ctx *kingpin.ParseContext) error {
		forcedRegions[regionName] = true
		return nil
	}
}

// disableDefaultRegions sets the region state to false for all regions which
// have not been explicitly enabled on the command line.
func disableDefaultRegions() {
	for regionName := range regionState {
		if _, ok := forcedRegions[regionName]; !ok {
			*regionState[regionName] = false
		}
	}
}

// Collect enabled regions coordinate
func Collect() (map[string][]*storage.Coordinate, error) {
	if *ddr {
		disableDefaultRegions()
	}

	enableRegionList := make([]string, 0)
	disableRegionList := make([]string, 0)
	for regionName, state := range regionState {
		if *state {
			enableRegionList = append(enableRegionList, regionName)
		} else {
			disableRegionList = append(disableRegionList, regionName)
		}
	}
	fmt.Printf("Enabled region list:  %v\n", enableRegionList)
	fmt.Printf("Disabled region list: %v\n", disableRegionList)

	fmt.Println("Start collect region coordinates.")
	regionCoordinates := make(map[string][]*storage.Coordinate)
	for regionName, state := range regionState {
		if *state {
			fmt.Printf("- Collecting %s...\n", regionName)
			coordinates, err := regionCollectors[regionName]()
			if err != nil {
				return nil, err
			}
			regionCoordinates[regionName] = coordinates
			fmt.Println("- Done!")
		}
	}
	return regionCoordinates, nil
}
