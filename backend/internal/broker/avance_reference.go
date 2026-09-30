package broker

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"encore.dev/rlog"
)

const (
	// avanceReferenceTTL is how long station and option lists are reused. Both change rarely.
	avanceReferenceTTL = 12 * time.Hour
	// avanceReferenceRetry spaces out refresh attempts after a failure, so an unreachable stations
	// page does not add a timeout to every search while it is down.
	avanceReferenceRetry = 5 * time.Minute
)

// avanceReference is the slow-changing data that decorates quotes: where each station is and when
// it opens, and what each option code is called.
type avanceReference struct {
	stations map[string]avanceStation
	options  map[string]avanceOption
}

type avanceStation struct {
	stationType string
	info        StationInfo
}

type avanceOption struct {
	name  string
	quant bool // whether the option can be booked in quantity
}

// The cache is package-level because a new Avance is built for every search.
var avanceReferenceCache struct {
	mu      sync.Mutex
	ref     avanceReference
	expires time.Time
}

// reference returns the cached station and option lists, refreshing them when they have expired.
// Reference data only decorates results, so a failed refresh is logged and the search carries on
// with the previous copy, or with none.
func (a *Avance) reference() avanceReference {
	c := &avanceReferenceCache
	c.mu.Lock()
	defer c.mu.Unlock()

	if time.Now().Before(c.expires) {
		return c.ref
	}

	ref, err := a.fetchReference()
	if err != nil {
		rlog.Warn("failed to refresh avance reference data, using the previous copy", "error", err)
		c.expires = time.Now().Add(avanceReferenceRetry)
		return c.ref
	}

	c.ref, c.expires = ref, time.Now().Add(avanceReferenceTTL)
	return c.ref
}

func (a *Avance) fetchReference() (avanceReference, error) {
	stationsBody, err := a.get("stations", url.Values{})
	if err != nil {
		return avanceReference{}, err
	}
	stations, err := parseAvanceStations(stationsBody)
	if err != nil {
		return avanceReference{}, err
	}

	optionsBody, err := a.get("options", url.Values{})
	if err != nil {
		return avanceReference{}, err
	}
	options, err := parseAvanceOptions(optionsBody)
	if err != nil {
		return avanceReference{}, err
	}

	return avanceReference{stations: stations, options: options}, nil
}

func parseAvanceStations(body []byte) (map[string]avanceStation, error) {
	var resp avanceStationsXML
	if err := xml.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("avance stations unmarshal response: %w", err)
	}

	stations := make(map[string]avanceStation, len(resp.Stations))
	for _, s := range resp.Stations {
		hours := make([]OpeningHoursItem, 0, len(s.Info.Hours))
		for _, h := range s.Info.Hours {
			item := OpeningHoursItem{Day: h.Day}
			if !h.Closed {
				item.OpenTime, item.CloseTime = h.OpensAt, h.ClosesAt
			}
			hours = append(hours, item)
		}

		stations[s.Code] = avanceStation{
			stationType: strings.TrimSpace(s.Info.Type),
			info: StationInfo{
				LocationInfo: strings.TrimSpace(s.Info.PickupInfo),
				Address:      joinNonEmpty(", ", s.Info.Address, s.Info.City, s.Info.ZipCode),
				PhoneNumber:  strings.TrimSpace(s.Info.Phone),
				OpeningHours: hours,
			},
		}
	}

	return stations, nil
}

func parseAvanceOptions(body []byte) (map[string]avanceOption, error) {
	var resp avanceOptionsXML
	if err := xml.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("avance options unmarshal response: %w", err)
	}

	options := make(map[string]avanceOption, len(resp.Options))
	for _, o := range resp.Options {
		options[o.Code] = avanceOption{name: strings.TrimSpace(o.Name), quant: o.Quant}
	}

	return options, nil
}

// joinNonEmpty joins the trimmed, non-empty parts with sep.
func joinNonEmpty(sep string, parts ...string) string {
	kept := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			kept = append(kept, p)
		}
	}
	return strings.Join(kept, sep)
}
