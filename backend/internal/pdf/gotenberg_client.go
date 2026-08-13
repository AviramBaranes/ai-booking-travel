package pdf

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"encore.dev"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

const (
	marginTop                = "0"
	marginBottom             = "0"
	marginLeft               = "0"
	marginRight              = "0"
	paperWidth               = "8.27"  // A4 width in inches
	paperHeight              = "11.69" // A4 height in inches
	gotenbergHtmlConvertPath = "/forms/chromium/convert/html"
)

// gotenberg is a PDFConverter implementation that uses the Gotenberg API for HTML to PDF conversion.
type gotenberg struct {
	baseURL    string
	httpClient *http.Client
}

// NewPdfConverter creates a new PDFConverter that uses Gotenberg for HTML to PDF conversion.
func NewPdfConverter(baseURL string) PDFConverter {
	return NewPdfConverterWithHTTPClient(baseURL, nil)
}

// The service account Gotenberg's IAM check accepts. Declared here rather than per service so that
// every caller signs its requests the same way without redeclaring the secret.
var secrets struct {
	GoogleServiceAccountJSON string
}

// NewDeployedPdfConverter creates a PDFConverter for the Gotenberg deployment at baseURL, signing
// requests with an identity token when one is needed. Deployed Gotenberg sits behind Google IAM, so
// it only answers authenticated calls; locally and under test it is a plain unauthenticated
// container and the identity token is skipped.
func NewDeployedPdfConverter(ctx context.Context, baseURL string) (PDFConverter, error) {
	meta := encore.Meta()
	if meta.Environment.Type == encore.EnvTest || meta.Environment.Cloud == encore.CloudLocal {
		return NewPdfConverter(baseURL), nil
	}

	httpClient, err := idtoken.NewClient(
		ctx,
		baseURL,
		option.WithAuthCredentialsJSON(option.ServiceAccount, []byte(secrets.GoogleServiceAccountJSON)),
	)
	if err != nil {
		return nil, fmt.Errorf("creating gotenberg identity token client: %w", err)
	}

	return NewPdfConverterWithHTTPClient(baseURL, httpClient), nil
}

func NewPdfConverterWithHTTPClient(baseURL string, httpClient *http.Client) PDFConverter {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: time.Second * 20}
	} else {
		httpClient.Timeout = time.Second * 20
	}

	return &gotenberg{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

// ConvertHTMLToPDF sends the HTML content to the Gotenberg API and returns the generated PDF as a byte slice.
func (g *gotenberg) ConvertHTMLToPDF(html string) ([]byte, error) {
	body, contentType, err := newMultipartWriter(html)
	if err != nil {
		return nil, fmt.Errorf("creating multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", g.baseURL+gotenbergHtmlConvertPath, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Gotenberg request failed (Is it running?): %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gotenberg returned non-OK status: %s, body: %s", resp.Status, body.String())
	}

	pdfData := &bytes.Buffer{}
	_, err = pdfData.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading Gotenberg response: %w", err)
	}

	return pdfData.Bytes(), nil
}

// newMultipartWriter creates a multipart form writer with the necessary fields and the HTML content as a file part.
func newMultipartWriter(html string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("marginTop", marginTop)
	_ = writer.WriteField("marginBottom", marginBottom)
	_ = writer.WriteField("marginLeft", marginLeft)
	_ = writer.WriteField("marginRight", marginRight)
	_ = writer.WriteField("paperWidth", paperWidth)
	_ = writer.WriteField("paperHeight", paperHeight)

	part, err := writer.CreateFormFile("files", "index.html")
	if err != nil {
		return nil, "", fmt.Errorf("creating form file: %w", err)
	}

	htmlBuffer := bytes.NewBufferString(html)
	part.Write(htmlBuffer.Bytes())

	err = writer.Close()
	if err != nil {
		return nil, "", fmt.Errorf("closing multipart writer: %w", err)
	}

	return body, writer.FormDataContentType(), nil
}
