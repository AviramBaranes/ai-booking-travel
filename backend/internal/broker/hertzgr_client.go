package broker

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"encore.dev/rlog"
)

// hertzGRTimeout bounds each Autohellas call. The station feed is under half a megabyte, and quotes
// run on the availability search path.
const hertzGRTimeout = 30 * time.Second

// hertzGRNamespace is Autohellas' namespace for the SOAP header and the operation wrappers. The OTA
// messages inside them use the OTA namespace.
const (
	hertzGRNamespace = "http://www.hertz.gr/xml"
	otaNamespace     = "http://www.opentravel.org/OTA/2003/05"
)

// hertzGRErrorBodyLimit caps how much of a failed response goes into an error message.
const hertzGRErrorBodyLimit = 2000

// HertzGR implements the Broker interface for Autohellas, the Hertz franchisee in Greece. Its OTA
// 2006A service also sells the Thrifty, Dollar and Firefly brands, and the same service will serve
// Autohellas' other countries.
type HertzGR struct {
	baseURL    string
	rates      []HertzGRRate
	logXML     bool
	httpClient *http.Client
	r          io.Reader
}

// HertzGRConfig configures the Hertz GR broker. The station import needs none of it.
type HertzGRConfig struct {
	// BaseURL is per-environment: production calls Autohellas directly, while every other
	// environment goes through the egress proxy, since Autohellas allow-lists our static address.
	BaseURL string
	// Rates are the rate codes Autohellas assigned us. Each is quoted as one plan on every car of
	// its brand, which the code's 4th letter selects.
	Rates []HertzGRRate
	// LogXML logs every request and response, with the password and validation code masked.
	// Autohellas certifies us from these exchanges before giving production access.
	LogXML bool
}

// HertzGRRate is one of our rate codes and the plan it becomes.
type HertzGRRate struct {
	Code         string
	PlanName     string
	ErpDayCharge float64
}

// NewHertzGR creates a new instance of the Hertz GR broker.
func NewHertzGR(cfg HertzGRConfig) *HertzGR {
	return &HertzGR{
		baseURL:    cfg.BaseURL,
		rates:      cfg.Rates,
		logXML:     cfg.LogXML,
		httpClient: &http.Client{Timeout: hertzGRTimeout},
	}
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

// hertzGRCredentials holds the four identifiers Autohellas needs on every call. They arrive packed
// into a single secret and are split apart at init.
type hertzGRCredentials struct {
	userName string
	password string
	ccn      string
	vc       string
}

var (
	hertzGRCreds    hertzGRCredentials
	hertzGRCredsErr error
)

// init unpacks the comma-separated Hertz GR credentials secret. A malformed value is recorded rather
// than panicked on, so a misconfigured Hertz GR does not take the other brokers down with it.
func init() {
	hertzGRCreds, hertzGRCredsErr = parseHertzGRCredentials(secrets.hertzGRCredentials)
}

func parseHertzGRCredentials(s string) (hertzGRCredentials, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 4 {
		return hertzGRCredentials{}, fmt.Errorf("hertz gr credentials: expected 4 comma-separated fields, got %d", len(parts))
	}

	c := hertzGRCredentials{
		userName: strings.TrimSpace(parts[0]),
		password: strings.TrimSpace(parts[1]),
		ccn:      strings.TrimSpace(parts[2]),
		vc:       strings.TrimSpace(parts[3]),
	}
	if c.userName == "" || c.password == "" || c.ccn == "" || c.vc == "" {
		return hertzGRCredentials{}, fmt.Errorf("hertz gr credentials: user name, password, CCN and VC must all be set")
	}

	return c, nil
}

// hertzGREnvelope is the SOAP request: the credentials in the header, and the OTA message wrapped
// in the operation's request element in the body.
type hertzGREnvelope struct {
	XMLName xml.Name `xml:"soap:Envelope"`
	Soap    string   `xml:"xmlns:soap,attr"`
	Header  struct {
		Auth struct {
			XMLName  xml.Name `xml:"AuthenticationHeader"`
			Xmlns    string   `xml:"xmlns,attr"`
			UserName string   `xml:"UserName"`
			Password string   `xml:"Password"`
		}
	} `xml:"soap:Header"`
	Body struct {
		Request struct {
			XMLName xml.Name
			Xmlns   string `xml:"xmlns,attr"`
			Message any
		}
	} `xml:"soap:Body"`
}

// hertzGRFault is the SOAP fault Autohellas returns when it cannot process a call at all, such as
// failed authentication.
type hertzGRFault struct {
	Code   string `xml:"faultcode"`
	String string `xml:"faultstring"`
}

// call invokes one Autohellas operation, for example call("VehicleAvailRate", rq, &rs). It wraps
// rq, an OTA request message, in the SOAP envelope and decodes the OTA response inside the
// operation's result element into rs. Only transport and SOAP failures are errors here: OTA errors
// come back inside rs, and the caller decides what they mean.
func (h *HertzGR) call(operation string, rq any, rs any) error {
	if hertzGRCredsErr != nil {
		return fmt.Errorf("hertz gr %s: %w", operation, hertzGRCredsErr)
	}

	var env hertzGREnvelope
	env.Soap = "http://schemas.xmlsoap.org/soap/envelope/"
	env.Header.Auth.Xmlns = hertzGRNamespace
	env.Header.Auth.UserName = hertzGRCreds.userName
	env.Header.Auth.Password = hertzGRCreds.password
	env.Body.Request.XMLName = xml.Name{Local: operation + "Request"}
	env.Body.Request.Xmlns = hertzGRNamespace
	env.Body.Request.Message = rq

	payload, err := xml.Marshal(env)
	if err != nil {
		return fmt.Errorf("hertz gr %s marshal request: %w", operation, err)
	}

	req, err := http.NewRequest(http.MethodPost, h.baseURL, bytes.NewReader(append([]byte(xml.Header), payload...)))
	if err != nil {
		return fmt.Errorf("hertz gr %s create request: %w", operation, err)
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", operation)

	res, err := h.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("hertz gr %s do request: %w", operation, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("hertz gr %s read response: %w", operation, err)
	}

	if h.logXML {
		rlog.Info("hertz gr exchange", "operation", operation, "status", res.StatusCode,
			"request", maskHertzGRSecrets(string(payload)), "response", string(body))
	}

	if err := decodeHertzGRResponse(body, operation+"Result", rs); err != nil {
		return fmt.Errorf("hertz gr %s (status %d): %w: %s", operation, res.StatusCode, err, truncateRunes(string(body), hertzGRErrorBodyLimit))
	}

	return nil
}

// maskHertzGRSecrets hides the password and validation code in a logged request. The user name and
// CCN stay, since Autohellas needs them to find our calls.
func maskHertzGRSecrets(s string) string {
	for _, secret := range []string{hertzGRCreds.password, hertzGRCreds.vc} {
		if secret != "" {
			s = strings.ReplaceAll(s, secret, "****")
		}
	}
	return s
}

// decodeHertzGRResponse finds the operation's result element in a SOAP response and decodes the OTA
// message inside it into rs. A SOAP fault, or a response without the result, is an error.
func decodeHertzGRResponse(body []byte, resultElement string, rs any) error {
	d := xml.NewDecoder(bytes.NewReader(body))
	inResult := false

	for {
		tok, err := d.Token()
		if err == io.EOF {
			return fmt.Errorf("no %s in response", resultElement)
		}
		if err != nil {
			return fmt.Errorf("read response: %w", err)
		}

		// Some services send the OTA message as escaped text inside the result rather than as an
		// element, which the guide's templates leave open.
		if text, ok := tok.(xml.CharData); ok && inResult {
			if s := strings.TrimSpace(string(text)); strings.HasPrefix(s, "<") {
				if err := xml.Unmarshal([]byte(s), rs); err != nil {
					return fmt.Errorf("decode %s text: %w", resultElement, err)
				}
				return nil
			}
			continue
		}

		start, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}

		switch {
		case start.Name.Local == "Fault":
			var f hertzGRFault
			if err := d.DecodeElement(&f, &start); err != nil {
				return fmt.Errorf("decode soap fault: %w", err)
			}
			return fmt.Errorf("soap fault %s: %s", strings.TrimSpace(f.Code), strings.TrimSpace(f.String))
		case start.Name.Local == resultElement:
			inResult = true
		case inResult:
			// The first element inside the result is the OTA response message.
			if err := d.DecodeElement(rs, &start); err != nil {
				return fmt.Errorf("decode %s: %w", start.Name.Local, err)
			}
			return nil
		}
	}
}
