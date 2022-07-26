package region

import (
	"errors"
	"fmt"
	"github.com/gogroup/coordinate/storage"
	"github.com/morikuni/failure"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"time"
)

var (
	fromSnapshots = kingpin.Flag(
		"region.from-snapshots",
		"Get data from snapshots instead of online.",
	).Default("false").Bool()
	ddr = kingpin.Flag(
		"region.disable-defaults",
		"Set all regions to disabled by default.",
	).Default("false").Bool()
)

const (
	defaultEnabled  = true
	defaultDisabled = false
)

var (
	collectors    = make(map[string]func() ([]*storage.Coordinate, error))
	snapshots     = make(map[string]func() ([]*storage.Coordinate, time.Time, error))
	regionState   = make(map[string]*bool)
	forcedRegions = map[string]bool{} // forcedRegions which have been explicitly enabled or disabled
)

func registerRegion(regionName string, isDefaultEnabled bool, collector func() ([]*storage.Coordinate, error), snapshot func() ([]*storage.Coordinate, time.Time, error)) {
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

	collectors[regionName] = collector
	snapshots[regionName] = snapshot
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
func Collect(logger *log.Logger) (map[string][]*storage.Coordinate, error) {
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
	logger.Info(fmt.Sprintf("Enable  %1d regions: %v", len(enableRegionList), enableRegionList))
	logger.Info(fmt.Sprintf("Disable %1d regions: %v", len(disableRegionList), disableRegionList))
	if len(enableRegionList) == 0 {
		return nil, failure.Wrap(errors.New("no region enable"))
	}

	logger.Info("Start get region coordinates.")
	regionCoordinates := make(map[string][]*storage.Coordinate)
	for _, regionName := range enableRegionList {
		var (
			coordinates []*storage.Coordinate
			err         error
		)
		if *fromSnapshots {
			logger.Info(fmt.Sprintf("- Parsing %s snapshot...", regionName))
			var modTime time.Time
			coordinates, modTime, err = snapshots[regionName]()
			if err == nil {
				logger.Info("- Snapshot time: ", modTime.Format("2006-01-02 15:04:06"))
			}
		} else {
			logger.Info(fmt.Sprintf("- Collecting %s...", regionName))
			coordinates, err = collectors[regionName]()
		}
		if err != nil {
			return nil, err
		}
		regionCoordinates[regionName] = coordinates
		logger.Info("- Done!")
	}
	return regionCoordinates, nil
}
