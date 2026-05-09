package pdf

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"
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
	return &gotenberg{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: time.Second * 20},
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
