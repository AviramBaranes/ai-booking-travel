package broker

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"encore.dev/rlog"
)

const (
	// hertzGRFleetURL is Autohellas' public fleet feed: every car group of every country and brand,
	// with a picture. Availability responses can leave PictureURL empty, so it fills the gap.
	hertzGRFleetURL = "https://prdota01.hertz.gr/OTAServices/Fleet.xml"
	// hertzGRImageHost serves the fleet feed's relative image paths.
	hertzGRImageHost = "https://images.autohellas.gr/"
)

type hertzGRFleetFeed struct {
	Countries []struct {
		Code   string `xml:"Code,attr"`
		Brands []struct {
			Code string `xml:"Code,attr"`
			Cars []struct {
				Group string `xml:"Group,attr"`
				Image string `xml:"SmallImage,attr"`
			} `xml:"Agent>Car"`
		} `xml:"Brand"`
	} `xml:"Country"`
}

// The cache is package-level because a new HertzGR is built for every search.
var hertzGRFleetCache struct {
	mu      sync.Mutex
	images  map[string]string
	expires time.Time
}

// hertzGRFleetKey identifies a car group's picture: groups are per country and brand.
func hertzGRFleetKey(country, brand, group string) string {
	return strings.ToUpper(strings.TrimSpace(country)) + "|" + strings.TrimSpace(brand) + "|" + strings.TrimSpace(group)
}

// fleetImages returns each car group's picture URL, keyed by hertzGRFleetKey, refreshing the cached
// copy when it has expired. Pictures only decorate results, so a failed refresh is logged and the
// search carries on with the previous copy, or with none.
func (h *HertzGR) fleetImages() map[string]string {
	c := &hertzGRFleetCache
	c.mu.Lock()
	defer c.mu.Unlock()

	if time.Now().Before(c.expires) {
		return c.images
	}

	images, err := h.fetchFleetImages()
	if err != nil {
		rlog.Warn("failed to refresh hertz gr fleet, using the previous copy", "error", err)
		c.expires = time.Now().Add(hertzGRStationsRetry)
		return c.images
	}

	c.images, c.expires = images, time.Now().Add(hertzGRStationsTTL)
	return c.images
}

func (h *HertzGR) fetchFleetImages() (map[string]string, error) {
	res, err := h.httpClient.Get(hertzGRFleetURL)
	if err != nil {
		return nil, fmt.Errorf("hertz gr fleet: fetch feed: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hertz gr fleet: fetch feed: status %d", res.StatusCode)
	}

	var feed hertzGRFleetFeed
	if err := xml.NewDecoder(res.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("hertz gr fleet: decode feed: %w", err)
	}

	return hertzGRFleetImages(feed), nil
}

func hertzGRFleetImages(feed hertzGRFleetFeed) map[string]string {
	images := make(map[string]string)
	for _, country := range feed.Countries {
		for _, brand := range country.Brands {
			for _, car := range brand.Cars {
				key := hertzGRFleetKey(country.Code, brand.Code, car.Group)
				if img := strings.TrimLeft(strings.TrimSpace(car.Image), "/"); img != "" && images[key] == "" {
					images[key] = hertzGRImageHost + img
				}
			}
		}
	}
	return images
}
