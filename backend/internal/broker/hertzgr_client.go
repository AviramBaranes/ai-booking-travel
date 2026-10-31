package broker

import (
	"io"
	"net/http"
	"time"
)

// hertzGRTimeout bounds each Autohellas call. The station feed is under half a megabyte.
const hertzGRTimeout = 30 * time.Second

// HertzGR implements the Broker interface for Autohellas, the Hertz franchisee in Greece. Its OTA
// 2006A service also sells the Thrifty, Dollar and Firefly brands, and the same service will serve
// Autohellas' other countries.
type HertzGR struct {
	httpClient *http.Client
	r          io.Reader
}

// NewHertzGR creates a new instance of the Hertz GR broker.
func NewHertzGR() *HertzGR {
	return &HertzGR{httpClient: &http.Client{Timeout: hertzGRTimeout}}
}

// NewHertzGRWithReader creates a Hertz GR broker that reads its station feed from r instead of
// fetching it.
func NewHertzGRWithReader(r io.Reader) *HertzGR {
	return &HertzGR{r: r}
}

// Name returns the name of the broker
func (h *HertzGR) Name() Name {
	return BrokerHertzGR
}
