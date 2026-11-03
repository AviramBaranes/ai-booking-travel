package broker

import (
	"encoding/xml"
	"fmt"
	"strings"
	"sync"
	"time"

	"encore.dev/rlog"
)

const (
	// hertzGRStationsTTL is how long the station feed is reused on searches. It changes rarely.
	hertzGRStationsTTL = 12 * time.Hour
	// hertzGRStationsRetry spaces out refresh attempts after a failure, so an unreachable feed does
	// not add a timeout to every search while it is down.
	hertzGRStationsRetry = 5 * time.Minute

	// hertzGRScheduleLayout is how the feed writes a season's first and last moment.
	hertzGRScheduleLayout = "02/01/2006 15:04:05"
)

// hertzGRDays maps the feed's day attributes onto the day names the frontend translates, in the
// order a week is listed.
var hertzGRDays = []struct{ attr, name string }{
	{"Sun", "Sunday"},
	{"Mon", "Monday"},
	{"Tue", "Tuesday"},
	{"Wed", "Wednesday"},
	{"Thu", "Thursday"},
	{"Fri", "Friday"},
	{"Sat", "Saturday"},
}

// The cache is package-level because a new HertzGR is built for every search.
var hertzGRStationsCache struct {
	mu       sync.Mutex
	stations map[string]hertzGRStation
	expires  time.Time
}

// stations returns every station in the feed by code, refreshing the cached copy when it has
// expired. Station details only decorate results, so a failed refresh is logged and the search
// carries on with the previous copy, or with none.
func (h *HertzGR) stations() map[string]hertzGRStation {
	c := &hertzGRStationsCache
	c.mu.Lock()
	defer c.mu.Unlock()

	if time.Now().Before(c.expires) {
		return c.stations
	}

	stations, err := h.fetchStations()
	if err != nil {
		rlog.Warn("failed to refresh hertz gr stations, using the previous copy", "error", err)
		c.expires = time.Now().Add(hertzGRStationsRetry)
		return c.stations
	}

	c.stations, c.expires = stations, time.Now().Add(hertzGRStationsTTL)
	return c.stations
}

func (h *HertzGR) fetchStations() (map[string]hertzGRStation, error) {
	r, err := h.openLocationsFeed()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var feed hertzGRLocationsFeed
	if err := xml.NewDecoder(r).Decode(&feed); err != nil {
		return nil, fmt.Errorf("hertz gr stations: decode feed: %w", err)
	}

	// A station shared by several brands is listed under each, with the same details.
	stations := make(map[string]hertzGRStation)
	for _, country := range feed.Countries {
		for _, brand := range country.Brands {
			for _, s := range brand.Stations {
				code := strings.TrimSpace(s.Code)
				if _, ok := stations[code]; !ok && code != "" {
					stations[code] = s
				}
			}
		}
	}

	return stations, nil
}

// stationInfo describes a station for a rental on the given date (yyyy-mm-dd), showing the opening
// hours of the season that covers it.
func (s hertzGRStation) stationInfo(date string) StationInfo {
	return StationInfo{
		Address:      joinNonEmpty(", ", s.AddressLine, s.City),
		PhoneNumber:  strings.TrimSpace(s.Telephone),
		OpeningHours: s.openingHours(date),
	}
}

// openingHours lists the week of the season covering date, one row per opening and "Closed" for a
// day the season does not open. It is empty when no season covers the date.
func (s hertzGRStation) openingHours(date string) []OpeningHoursItem {
	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil
	}

	for _, sch := range s.Schedules {
		start, errStart := time.Parse(hertzGRScheduleLayout, strings.TrimSpace(sch.Start))
		end, errEnd := time.Parse(hertzGRScheduleLayout, strings.TrimSpace(sch.End))
		if errStart != nil || errEnd != nil || day.Before(start.Truncate(24*time.Hour)) || day.After(end) {
			continue
		}

		ranges := make(map[string][]string)
		for _, t := range sch.Times {
			for _, a := range t.Days {
				ranges[a.Name.Local] = append(ranges[a.Name.Local], a.Value)
			}
		}

		hours := make([]OpeningHoursItem, 0, len(hertzGRDays))
		for _, d := range hertzGRDays {
			if len(ranges[d.attr]) == 0 {
				hours = append(hours, OpeningHoursItem{Day: d.name, OpenTime: "Closed"})
				continue
			}
			for _, rg := range ranges[d.attr] {
				opens, closes, _ := strings.Cut(rg, "-")
				hours = append(hours, OpeningHoursItem{
					Day:       d.name,
					OpenTime:  strings.TrimSpace(opens),
					CloseTime: strings.TrimSpace(closes),
				})
			}
		}
		return hours
	}

	return nil
}
