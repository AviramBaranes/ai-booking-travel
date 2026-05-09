package pdf

// PDFConverter is an interface for converting HTML to PDF.
type PDFConverter interface {
	ConvertHTMLToPDF(html string) ([]byte, error)
}
