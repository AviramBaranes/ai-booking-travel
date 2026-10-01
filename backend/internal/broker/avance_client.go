package broker

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// avanceTimeout bounds each Wheelsys call. Quotes answer in well under a second, and they run on
// the availability search path, so a hung call must not hold the whole search for defaultTimeout.
const avanceTimeout = 30 * time.Second

// Avance is a struct that implements the Broker interface for the Avance car rental service,
// which is served by the Wheelsys "Link v3" API.
type Avance struct {
	baseURL string
	// erpDayCharges is our own ERP day charge, keyed by product (CDP): it differs between the
	// standard and the zero-excess product.
	erpDayCharges map[string]float64
	httpClient    *http.Client
	r             io.Reader
}

// AvanceConfig configures the Avance broker. Only availability reads the ERP day charges, and the
// voucher needs none of it.
type AvanceConfig struct {
	// BaseURL is per-environment: production reaches Wheelsys directly, while every other
	// environment goes through the egress proxy, since only production's address is allow-listed.
	BaseURL                string
	StandardErpDayCharge   float64
	ZeroExcessErpDayCharge float64
}

// NewAvance creates a new instance of the Avance broker.
func NewAvance(cfg AvanceConfig) *Avance {
	return &Avance{
		baseURL: cfg.BaseURL,
		erpDayCharges: map[string]float64{
			avanceCDPStandard:   cfg.StandardErpDayCharge,
			avanceCDPZeroExcess: cfg.ZeroExcessErpDayCharge,
		},
		httpClient: &http.Client{Timeout: avanceTimeout},
	}
}

// NewAvanceWithReader creates an Avance broker that reads its locations from the given workbook.
func NewAvanceWithReader(r io.Reader) *Avance {
	return &Avance{r: r}
}

// Name returns the name of the broker
func (a *Avance) Name() Name {
	return BrokerAvance
}

// avanceCredentials holds the three identifiers Wheelsys needs on every call. They arrive packed
// into a single secret and are split apart at init.
type avanceCredentials struct {
	accountNo string
	linkCode  string
	agentCode string
}

var (
	avanceCreds    avanceCredentials
	avanceCredsErr error
)

// init unpacks the comma-separated Avance credentials secret. Encore populates the secrets struct
// before package initialisation, so reading it here is safe. A malformed value is recorded rather
// than panicked on, so a misconfigured Avance does not take the other brokers down with it.
func init() {
	avanceCreds, avanceCredsErr = parseAvanceCredentials(secrets.avanceCredentials)
}

func parseAvanceCredentials(s string) (avanceCredentials, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		return avanceCredentials{}, fmt.Errorf("avance credentials: expected 3 comma-separated fields, got %d", len(parts))
	}

	c := avanceCredentials{
		accountNo: strings.TrimSpace(parts[0]),
		linkCode:  strings.TrimSpace(parts[1]),
		agentCode: strings.TrimSpace(parts[2]),
	}
	if c.accountNo == "" || c.linkCode == "" || c.agentCode == "" {
		return avanceCredentials{}, fmt.Errorf("avance credentials: account number, link code and agent code must all be set")
	}

	return c, nil
}

// get calls a single Wheelsys Link v3 page, for example get("price-quote", q), and returns the raw
// response body. The link code is part of the page name and the agent code is a query parameter,
// so both are applied here rather than by the caller.
func (a *Avance) get(page string, q url.Values) ([]byte, error) {
	if avanceCredsErr != nil {
		return nil, fmt.Errorf("avance %s: %w", page, avanceCredsErr)
	}

	q.Set("agent", avanceCreds.agentCode)
	endpoint := fmt.Sprintf("%s/%s/link/v3/%s_%s.html?%s",
		a.baseURL, avanceCreds.accountNo, page, avanceCreds.linkCode, q.Encode())

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("avance %s create request: %w", page, err)
	}

	req.Header.Set("Accept", "application/xml,text/xml;q=0.9,*/*;q=0.8")

	res, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("avance %s do request: %w", page, err)
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("avance %s read response: %w", page, err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("avance %s unexpected status %d: %s", page, res.StatusCode, string(b))
	}

	return b, nil
}
